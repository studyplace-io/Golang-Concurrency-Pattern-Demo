package task_job_mode

import (
	"container/list"
	"golanglearning/new_project/Golang-Concurrency-Pattern-Demo/task-job-mode/model"
	"sync"
)

// taskNode 从server发来的task以queue方式存储在内存中。
// taskNode为链表中节点的定义
type taskNode struct {
	task   *model.Task // 任务信息
	id     int64       // 任务id
	jobIds []int64     // 任务下job id列表
}

// taskQueue 任务队列，存储server发来的task信息
type taskQueue struct {
	list list.List
	lock sync.Mutex // 控制写入并发
}

// Push 新增任务至队列中
func (q *taskQueue) push(node *taskNode) {
	// 加入queue中
	q.lock.Lock()
	defer q.lock.Unlock()
	q.list.PushBack(node)
}

// RequestTask 请求一个新任务。
// 返回结果中bool表示是否请求到了新任务，为true时任务信息必不为nil。
func (q *taskQueue) requestTask() (*taskNode, bool) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if ele := q.list.Front(); ele != nil {
		q.list.Remove(ele)
		return ele.Value.(*taskNode), true
	}
	return nil, false
}

// AbortTask 取消任务
func (q *taskQueue) abortTask(taskId string) bool {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.list.Back()
	for ele := q.list.Front(); ele != nil; ele = ele.Next() {
		if node := ele.Value.(*taskNode); node.task.TaskId == taskId {
			q.list.Remove(ele)
			return true
		}
	}
	return false
}

func (q *taskQueue) length() int {
	return q.list.Len()
}
