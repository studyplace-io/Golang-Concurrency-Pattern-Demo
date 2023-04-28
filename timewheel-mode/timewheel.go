package timewheel

import (
	"container/list"
	"errors"
	"fmt"
	"sync"
	"time"
)

// TimeWheel 时间轮struct
type TimeWheel struct {
	// 时间轮盘时间间隔
	interval time.Duration
	// 时间轮盘每个位置存储的Task列表，使用链表
	slots  []*list.List
	ticker *time.Ticker
	// 时间轮盘当前的位置
	currentPos int
	// 时间轮盘的齿轮数 interval*slotNums就是时间轮盘转一圈走过的时间
	slotNums    int
	addTaskC    chan *Task
	removeTaskC chan *Task
	stopC       chan bool
	// Map结构来存储Task对象，key是Task.key，value是Task在双向链表中的存储对象，结构是list.Element
	taskRecords *sync.Map
	// 需要执行的任务，如果时间轮盘上的Task执行同一个Job，可以直接实例化到TimeWheel结构体中。
	// 此处的优先级低于Task中的Job参数
	job       Job
	isRunning bool
}

type TaskName string

// Job 需要执行的Job的函数结构体
type Job func(TaskName)

// Task 时间轮盘上需要执行的任务
type Task struct {
	// 用来标识task对象，是唯一的
	TaskName TaskName
	// 任务周期
	interval time.Duration
	// 任务的创建时间
	createdTime time.Time
	// 任务在轮盘的位置
	pos int
	// 任务需要在轮盘走多少圈才能执行
	circle int
	// 任务需要执行的Job，优先级高于TimeWheel中的Job
	job Job
	// 任务需要执行的次数，如果需要一直执行，设置成-1
	times int
}

var tw *TimeWheel
var once sync.Once

// ErrDuplicateTaskKey is an error for duplicate task key
var ErrDuplicateTaskKey = errors.New("duplicate task key")

// ErrTaskKeyNotFount is an error when task key is not found
var ErrTaskKeyNotFount = errors.New("task key doesn't existed in task list, please check your input")

// CreateTimeWheel 用来实现TimeWheel的单例模式
func CreateTimeWheel(interval time.Duration, slotNums int, job Job) *TimeWheel {
	once.Do(func() {
		tw = New(interval, slotNums, job)
	})
	return tw
}

// GetTimeWheel 返回tw对象
func GetTimeWheel() *TimeWheel {
	return tw
}

// New 初始化一个TimeWheel对象
func New(interval time.Duration, slotNums int, job Job) *TimeWheel {
	if interval <= 0 || slotNums <= 0 {
		return nil
	}
	tw := &TimeWheel{
		interval:    interval,
		slots:       make([]*list.List, slotNums),
		currentPos:  0,
		slotNums:    slotNums,
		addTaskC:    make(chan *Task),
		removeTaskC: make(chan *Task),
		stopC:       make(chan bool),
		taskRecords: &sync.Map{},
		job:         job,
		isRunning:   false,
	}

	tw.initSlots()
	return tw
}

// Start 启动时间轮盘
func (tw *TimeWheel) Start() {
	tw.ticker = time.NewTicker(tw.interval)
	go tw.start()
	tw.isRunning = true
}

// Stop 关闭时间轮盘
func (tw *TimeWheel) Stop() {
	tw.stopC <- true
	tw.isRunning = false
}

// IsRunning 检查一下时间轮盘的是否在正常运行
func (tw *TimeWheel) IsRunning() bool {
	return tw.isRunning
}

// AddTask 向时间轮盘添加任务的开放函数
// @param interval    任务的周期
// @param taskName    任务名，必须是唯一的，否则添加任务的时候会失败
// @param createTime  任务的创建时间
func (tw *TimeWheel) AddTask(interval time.Duration, taskName TaskName, createdTime time.Time, times int, job Job) error {
	if interval <= 0 || taskName == "" {
		return errors.New("invalid task params")
	}

	// 检查Task.Key是否已经存在
	_, ok := tw.taskRecords.Load(taskName)
	if ok {
		return ErrDuplicateTaskKey
	}

	tw.addTaskC <- &Task{
		TaskName:    taskName,
		interval:    interval,
		createdTime: createdTime,
		job:         job,
		times:       times,
	}

	return nil
}

// RemoveTask 从时间轮盘删除任务的公共函数
func (tw *TimeWheel) RemoveTask(taskName TaskName) error {
	if taskName == "" {
		return nil
	}

	// 检查该Task是否存在
	val, ok := tw.taskRecords.Load(taskName)
	if !ok {
		return ErrTaskKeyNotFount
	}

	// 需要删除的task，统一放入removeTaskC中
	task := val.(*list.Element).Value.(*Task)
	tw.removeTaskC <- task
	return nil
}

// 初始化时间轮盘
// 帮每个轮盘上的卡槽初始化链表
func (tw *TimeWheel) initSlots() {
	for i := 0; i < tw.slotNums; i++ {
		tw.slots[i] = list.New()
	}
}

// 启动时间轮盘的内部函数
func (tw *TimeWheel) start() {
	for {
		select {
		case <-tw.ticker.C:
			tw.checkAndRunTask()
		case task := <-tw.addTaskC:
			// 此处利用Task.createTime来定位任务在时间轮盘的位置和执行圈数
			// 如果直接用任务的周期来定位位置，那么在服务重启的时候，任务周器相同的点会被定位到相同的卡槽，
			// 会造成任务过度集中
			tw.addTask(task, false)
		case task := <-tw.removeTaskC:
			tw.removeTask(task)
		case <-tw.stopC:
			tw.ticker.Stop()
			return
		}
	}
}

// checkAndRunTask 检查该轮盘点位上的Task，看哪个需要执行
func (tw *TimeWheel) checkAndRunTask() {

	// 获取该轮盘位置的双向链表
	currentList := tw.slots[tw.currentPos]

	if currentList != nil {
		for item := currentList.Front(); item != nil; {
			task := item.Value.(*Task)
			// 如果圈数>0，表示还没到执行时间，更新圈数
			if task.circle > 0 {
				task.circle--
				item = item.Next()
				continue
			}

			// 执行任务时，Task.job是第一优先级，然后是TimeWheel.job
			if task.job != nil {
				go task.job(task.TaskName)
			} else if tw.job != nil {
				go tw.job(task.TaskName)
			} else {
				fmt.Println(fmt.Sprintf("The task %v don't have job to run", task.TaskName))
			}

			// 执行完成以后，将该任务从时间轮盘删除
			next := item.Next()
			tw.taskRecords.Delete(task.TaskName)
			currentList.Remove(item)

			item = next

			// 重新添加任务到时间轮盘，用Task.interval来获取下一次执行的轮盘位置
			// 如果times==0,说明已经完成执行周期,不需要再添加任务回时间轮盘
			if task.times != 0 {
				if task.times < 0 {
					tw.addTask(task, true)
				} else {
					task.times--
					tw.addTask(task, true)
				}

			}
		}
	}

	// 轮盘前进一步
	if tw.currentPos == tw.slotNums-1 {
		tw.currentPos = 0
	} else {
		tw.currentPos++
	}
}

// 添加任务的内部函数
// @param task       Task  Task对象
// @param byInterval bool  生成Task在时间轮盘位置和圈数的方式，true表示利用Task.interval来生成，false表示利用Task.createTime生成
func (tw *TimeWheel) addTask(task *Task, byInterval bool) {
	var pos, circle int
	if byInterval {
		pos, circle = tw.getPosAndCircleByInterval(task.interval)
	} else {
		pos, circle = tw.getPosAndCircleByCreatedTime(task.createdTime, task.interval)
	}

	task.circle = circle
	task.pos = pos

	element := tw.slots[pos].PushBack(task)
	tw.taskRecords.Store(task.TaskName, element)
}

// 删除任务的内部函数
func (tw *TimeWheel) removeTask(task *Task) {
	// 从map结构中删除
	val, _ := tw.taskRecords.Load(task.TaskName)
	tw.taskRecords.Delete(task.TaskName)

	// 通过TimeWheel.slots获取任务的
	currentList := tw.slots[task.pos]
	currentList.Remove(val.(*list.Element))
}

// 该函数通过任务的周期来计算下次执行的位置和圈数
func (tw *TimeWheel) getPosAndCircleByInterval(d time.Duration) (int, int) {
	delaySeconds := int(d.Seconds())
	intervalSeconds := int(tw.interval.Seconds())
	circle := delaySeconds / intervalSeconds / tw.slotNums
	pos := (tw.currentPos + delaySeconds/intervalSeconds) % tw.slotNums

	// 特殊case，当计算的位置和当前位置重叠时，因为当前位置已经走过了，所以circle需要减一
	if pos == tw.currentPos && circle != 0 {
		circle--
	}
	return pos, circle
}

// 该函数用任务的创建时间来计算下次执行的位置和圈数
func (tw *TimeWheel) getPosAndCircleByCreatedTime(createdTime time.Time, d time.Duration) (int, int) {

	passedTime := time.Since(createdTime)
	passedSeconds := int(passedTime.Seconds())
	delaySeconds := int(d.Seconds())
	intervalSeconds := int(tw.interval.Seconds())

	circle := delaySeconds / intervalSeconds / tw.slotNums
	pos := (tw.currentPos + (delaySeconds-(passedSeconds%delaySeconds))/intervalSeconds) % tw.slotNums

	// 特殊case，当计算的位置和当前位置重叠时，因为当前位置已经走过了，所以circle需要减一
	if pos == tw.currentPos && circle != 0 {
		circle--
	}

	return pos, circle
}
