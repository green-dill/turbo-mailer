package models

import (
	"time"

	"gorm.io/gorm"
)

var (
	ALL = []any{
		&Pool{},
		&PoolSender{},
		&ServerIP{},
		&ServerIPLog{},
		&Task{},
		&TaskLog{},
		&TaskLogArchive{},
		&TaskPool{},
		&User{},
	}
)

type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
