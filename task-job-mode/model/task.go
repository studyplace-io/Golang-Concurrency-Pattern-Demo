package model

type Task struct {
	TaskId    string `json:"task_id"`
	TaskGroup string `json:"task_group"`
	ClientId  string `json:"client_id"`
	Jobs      []Job  `json:"jobs"`
}

type Job struct {
	JobId   string            `json:"job_id"`
	Service string            `json:"service"`
	Action  string            `json:"action"`
	Args    map[string]string `json:"args"`
}

type Result struct {
	JobId      string      `json:"job_id"`
	TaskId     string      `json:"task_id"`
	ClientName string      `json:"client_name"`
	Status     string      `json:"status"`
	TaskStatus string      `json:"task_status"`
	Data       interface{} `json:"data"`
}
