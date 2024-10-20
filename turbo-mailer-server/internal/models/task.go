package models

import (
	"time"

	"gorm.io/gorm"
)

var (
	TaskStatePending    = "pending"
	TaskStateDispatched = "dispatched"
	TaskStateFinished   = "finished"
)

type Task struct {
	gorm.Model
	Subject            string      `gorm:"column:subject;type:varchar(255);not null" json:"subject"`
	ContextType        string      `gorm:"column:context_type;type:varchar(255);not null" json:"context_type"` // MIME type
	Content            string      `gorm:"column:content;type:text;not null" json:"content"`
	Receivers          []string    `gorm:"column:receivers;type:text;not null" json:"receivers"`
	State              string      `gorm:"column:state;type:varchar(255);not null" json:"state"`
	Pools              []*TaskPool `gorm:"many2many:task_pool_ref;" json:"pools"`
	MaxDispatchPreHour int         `gorm:"column:max_dispatch_pre_hour;type:int;not null" json:"max_dispatch_pre_hour"`
	LastDispatchAt     *time.Time  `gorm:"column:last_dispatch_at;not null" json:"last_dispatch_at"`
}

type TaskPool struct {
	gorm.Model
	TaskID uint `gorm:"column:task_id;type:bigint;not null;uniqueIndex:idx_task_pool" json:"task_id"`
	PoolID uint `gorm:"column:pool_id;type:bigint;not null;uniqueIndex:idx_task_pool" json:"pool_id"`
	Pool   Pool `gorm:"foreignKey:id" json:"pool"`
	Weight int  `gorm:"column:weight;type:int;not null;default:1" json:"weight"`
}

var (
	TaskLogStatePending = "pending"
	TaskLogStateSending = "sending"
	TaskLogStateSent    = "sent"
	TaskLogStateFailed  = "failed"
)

type TaskLog struct {
	ID        uint      `gorm:"column:id;type:bigint;primaryKey;autoIncrement" json:"id"`
	TaskID    uint      `gorm:"column:task_id;type:bigint;not null" json:"task_id"`
	PoolID    uint      `gorm:"column:pool_id;type:bigint;not null" json:"pool_id"`
	Sender    string    `gorm:"column:sender;type:varchar(255);not null" json:"sender"`
	Receiver  string    `gorm:"column:receiver;type:varchar(255);not null" json:"receiver"`
	State     string    `gorm:"column:state;type:varchar(64);not null" json:"state"`
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"created_at"`
}

type TaskLogArchive struct {
	TaskLog
	ArchivedAt time.Time `gorm:"column:archived_at;not null" json:"archived_at"`
}
