package model

import "sync"

type Task struct {
	Name      string   `yaml:"name"`
	Type      string   `yaml:"type"`
	Commands  []string `yaml:"commands"`
	Children  []Task   `yaml:"children"`
	OnFailure string   `yaml:"on_failure"` // 新增字段
	Failed    bool     // 新增字段
	Completed bool
}

type Status struct {
	sync.Mutex
	Tasks map[string]*Task
}
