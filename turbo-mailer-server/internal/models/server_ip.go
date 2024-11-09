package models

import (
	"time"
)

type ServerIP struct {
	BaseModel
	IP           string    `gorm:"column:ip;type:varchar(128);not null;unique" json:"ip"`
	LastActiveAt time.Time `gorm:"column:last_active_at;not null" json:"last_active_at"`
}

type ServerIPLog struct {
	ID        uint      `gorm:"column:id;type:bigint;primaryKey;autoIncrement" json:"id"`
	ServerIP  string    `gorm:"column:server_ip;type:varchar(128);not null" json:"server_ip"`
	Status    string    `gorm:"column:status;type:varchar(64);not null" json:"status"`
	Host      string    `gorm:"column:host;type:varchar(128);not null" json:"host"`
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"created_at"`
}
