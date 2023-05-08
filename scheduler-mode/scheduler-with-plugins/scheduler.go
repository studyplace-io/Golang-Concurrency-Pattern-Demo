package scheduler_with_plugins

import (
	_interface "golanglearning/new_project/Golang-Concurrency-Pattern-Demo/scheduler-mode/scheduler-with-plugins/interface"
	"golanglearning/new_project/Golang-Concurrency-Pattern-Demo/scheduler-mode/scheduler-with-plugins/nodes"
	"k8s.io/klog/v2"
	"sync"
)

// Scheduler 调度器
type Scheduler struct {
	options *schedulerOptions // 调度器配置
	name    string

	queue     *Queue
	pods      chan _interface.Pod // pod队列
	nodeInfos *nodes.NodeInfos    // 存储所有node的信息
	workers   int                 // 控制并发数
	plugins   []_interface.Plugin // 插件

	wg     sync.WaitGroup
	stopC  chan struct{} // 通知
	logger klog.Logger
}

type schedulerOptions struct {
	test          string
	test1         int
	numWorker     int
	queueCapacity int
}

// SchedulerOption 选项模式
type SchedulerOption func(options *schedulerOptions)

// defaultOptions 默认配置
var defaultOptions = schedulerOptions{
	numWorker:     3,
	queueCapacity: 10,
	test:          "test",
	test1:         10,
}

func WithNumWorker(numWorker int) SchedulerOption {
	return func(options *schedulerOptions) {
		options.numWorker = numWorker
	}
}

func WithQueueCapacity(queueCapacity int) SchedulerOption {
	return func(options *schedulerOptions) {
		options.queueCapacity = queueCapacity
	}
}

func WithTest1(test1 int) SchedulerOption {
	return func(options *schedulerOptions) {
		options.test1 = test1
	}
}

func WithTest(test string) SchedulerOption {
	return func(options *schedulerOptions) {
		options.test = test
	}
}

// AddPlugin 加入插件
func (s *Scheduler) AddPlugin(plugin _interface.Plugin) {
	s.plugins = append(s.plugins, plugin)
}

// AddPod 放入pod
func (s *Scheduler) AddPod(pod _interface.Pod) {
	s.queue.activeQ <- pod
}

// run 启动调度器
func (s *Scheduler) run() {
	// 启动队列
	go s.queue.Run(s.stopC)

	for i := 0; i < s.workers; i++ {
		s.wg.Add(1)
		// 并发
		go func() {
			defer s.wg.Done()

			for {
				select {
				case <-s.stopC: // 退出通知
					s.logger.Info("return scheduler...")
					return
				case t := <-s.queue.Get(): // 取出队列中的pod

					t = s.runFilter(t)
					// 打分node
					t = s.runScorer(t)
					// 选出最好的node
					nodeInfo := s.selectHost(t)
					// 异步绑定
					go bind(t, nodeInfo)
				}

			}
		}()
	}
}

// Stop 停止
func (s *Scheduler) Stop() {
	if s.queue.Len() > 0 {
		s.logger.Info("scheduler queue still have element...")
	}
	// 通知退出
	close(s.stopC) // 通知
}

// runFilter 执行过滤插件
func (s *Scheduler) runFilter(pod _interface.Pod) _interface.Pod {
	s.logger.Info("runFilter...")
	if s.podFiltered(pod) {
		return pod
	}
	s.logger.Info("have no pod to run...")
	// 把无法调度的放入backoffQ中
	s.queue.Backoff(pod)
	return nil
}

// runScorer 执行打分插件
func (s *Scheduler) runScorer(pod _interface.Pod) _interface.Pod {
	s.logger.Info("runScorer...")
	var totalScore float64
	if pod == nil {
		s.logger.Info("have no pod to score...")
		return nil
	}
	for _, nodeInfo := range s.nodeInfos.NodeInfos {
		for _, plugin := range s.plugins {
			totalScore += plugin.Score(pod, nodeInfo)
		}
		if totalScore != 0 {
			pod.SetPodRecordNode(nodeInfo.NodeName, totalScore)
			totalScore = 0
		}
	}

	return pod
}

// selectHost 选出最好的node
func (s *Scheduler) selectHost(pod _interface.Pod) *nodes.NodeInfo {
	var resNodeInfo *nodes.NodeInfo
	if pod == nil {
		s.logger.Info("have no pod to select host...")
		return resNodeInfo
	}
	nodeList := pod.GetPodRecordNodeList()
	var maxScore float64
	var maxScoreNodeName string
	for _, node := range nodeList {
		if node.Score >= maxScore {
			maxScore = node.Score
			maxScoreNodeName = node.NodeName
		}
	}

	for _, schedulerNode := range s.nodeInfos.NodeInfos {
		if schedulerNode.NodeName == maxScoreNodeName {
			resNodeInfo = schedulerNode
		}
	}
	return resNodeInfo
}

// bind 绑定
func bind(pod _interface.Pod, nodeInfo *nodes.NodeInfo) {
	if pod == nil || nodeInfo == nil {
		klog.Info("have no pod or node to bind...")
		return
	}
	pod.SetNode(nodeInfo.NodeName)
}

// podFiltered 过滤插件
func (s *Scheduler) podFiltered(pod _interface.Pod) bool {
	for _, plugin := range s.plugins {
		if plugin.Filter(pod) {
			return true
		}
	}
	return false
}
