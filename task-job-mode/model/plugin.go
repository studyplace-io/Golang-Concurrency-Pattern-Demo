package model

type Plugin interface {
	Execute(*Job, func(interface{})) (ExecStatus, interface{}) // 执行任务
	ActionList() []*Action                                     // Plugin支持的操作列表
}

type ExecStatus string // 插件执行结果，可选值参考 internal/common/constants/task.go

type Action struct {
	Name    string
	Version string
}
