package tree_template

import (
	"github.com/study-io/Golang-Concurrency-Pattern-Demo/tree-template-engine/model"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type TreeTemplateEngine struct {
}

func NewTreeTemplateEngine() *TreeTemplateEngine {
	return &TreeTemplateEngine{}
}

func (tte *TreeTemplateEngine) RunTask(task *model.Task) {
	status := &model.Status{
		Tasks: make(map[string]*model.Task),
	}
	runTask(task, status)
}

func (tte *TreeTemplateEngine) ParseYaml(path string) *model.Task {
	// 读取YAML文件
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading YAML file: %v\n", err)
	}

	// 解析YAML文件

	var taskTree model.Task
	err = yaml.Unmarshal(data, &taskTree)
	if err != nil {
		log.Fatalf("Error parsing YAML: %v\n", err)
	}
	return &taskTree
}
