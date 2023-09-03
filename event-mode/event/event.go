package event

import "time"

const (
	Normal = "normal"
	Warning = "warning"
)

// EventType 事件种类
type EventType string

// Event 事件对象
type Event struct {
	// Reason 理由
	Reason    string
	// Message 消息内容
	Message   string
	// Source 事件源
	Source    string
	// Type 事件类型
	Type      EventType
	// Timestamp 写入时间
	Timestamp time.Time
}
