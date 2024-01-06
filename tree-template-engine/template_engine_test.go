package tree_template

import (
	"fmt"
	"testing"
)

func TestTreeTemplateTask(t *testing.T) {

	te := NewTreeTemplateEngine()
	te.RunTask(te.ParseYaml("./task.yaml"))

	fmt.Println("All tasks completed successfully.")
}
