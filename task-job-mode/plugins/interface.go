package plugins

import (
	"github.com/practice/Golang-Concurrency-Pattern-Demo/task-job-mode/model"
	"github.com/practice/Golang-Concurrency-Pattern-Demo/task-job-mode/plugins/demo"
)

func init() {
	RegisterPlugin("demo", &demo.Plugin{})
}

// plugins 插件列表。不关心并发，由代码保证串行注册
var plugins = make(map[string]model.Plugin)

// actionFilter 用于判断目标action是否存在
var actionFilter = make(map[string]map[string]struct{})

// RegisterPlugin 注册插件
func RegisterPlugin(service string, plugin model.Plugin) {
	plugins[service] = plugin
	actions := make(map[string]struct{})
	for _, action := range plugin.ActionList() {
		actions[action.Name] = struct{}{}
	}
	actionFilter[service] = actions
}

// SearchPlugin 搜索插件
func SearchPlugin(service string, action string) (model.Plugin, bool) {
	if plugin, ok := plugins[service]; ok {
		if _, ok := actionFilter[service][action]; ok {
			return plugin, true
		}
	}
	return nil, false
}
