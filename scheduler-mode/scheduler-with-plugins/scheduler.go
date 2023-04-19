package scheduler_with_plugins

import (
	"fmt"
	_interface "golanglearning/new_project/Golang-Concurrency-Pattern-Demo/scheduler-mode/scheduler-with-plugins/interface"
	"golanglearning/new_project/Golang-Concurrency-Pattern-Demo/scheduler-mode/scheduler-with-plugins/nodes"
	"sync"
)

// Scheduler 调度器
type Scheduler struct {
	pods      chan _interface.Pod // pod队列
	nodeInfos *nodes.NodeInfos    // 存储所有node的信息
	workers   int                 // 控制并发数
	plugins   []_interface.Plugin // 插件
	wg        sync.WaitGroup
}

// AddPlugin 加入插件
func (s *Scheduler) AddPlugin(plugin _interface.Plugin) {
	s.plugins = append(s.plugins, plugin)
}

// AddPod 放入pod
func (s *Scheduler) AddPod(pod _interface.Pod) {
	s.pods <- pod
}

// run 启动调度器
func (s *Scheduler) run() {

	for i := 0; i < s.workers; i++ {
		s.wg.Add(1)
		// 并发
		go func() {
			defer s.wg.Done()

			for {
				// 取出队列中的pod
				t := <-s.pods
				// 过滤pod
				t = s.runFilter(t)
				// 打分node
				t = s.runScorer(t)
				// 选出最好的node
				nodeInfo := s.selectHost(t)
				// 异步绑定
				go bind(t, nodeInfo)

			}
		}()
	}
}

// runFilter 执行过滤插件
func (s *Scheduler) runFilter(pod _interface.Pod) _interface.Pod {
	fmt.Println("runFilter。。。")
	if s.podFiltered(pod) {
		return pod
	}
	fmt.Println("没有runFilter。。。")
	return nil
}

// runScorer 执行打分插件
func (s *Scheduler) runScorer(pod _interface.Pod) _interface.Pod {
	fmt.Println("runScorer。。。")
	var totalScore float64
	if pod == nil {
		fmt.Println("没有runScorer操作")
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
		fmt.Println("没有selectHost操作")
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
		fmt.Println("没有bind操作")
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
