package pipeline_mode

import (
	"testing"
)

func TestTaskPipeline(t *testing.T) {
	task1 := NewTask("task1")
	task2 := NewTask("task2")
	task3 := NewTask("task3")

	// 流水线
	prepareTaskC := PrepareTask(task3, task2, task1)
	resultTaskC := ExecuteTask(prepareTaskC)
	AnalyzeTask(resultTaskC)

}
