package models

import (
	"gorm.io/gorm"
)

type Pool struct {
	gorm.Model
	Name        string       `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Description string       `gorm:"column:description;type:varchar(512)" json:"description,omitempty"`
	SenderCount int          `gorm:"column:sender_count;type:int;default:0" json:"sender_count"`
	Senders     []PoolSender `gorm:"foreignKey:PoolID;references:ID" json:"senders"`
}

type PoolSender struct {
	gorm.Model
	PoolID    uint   `gorm:"column:pool_id;type:bigint;not null" json:"pool_id"`
	FromName  string `gorm:"column:from_name;type:varchar(255);not null" json:"from_name"`
	FromEmail string `gorm:"column:from_email;type:varchar(255);not null" json:"from_email"`
	ReplyTo   string `gorm:"column:reply_to;type:varchar(255)" json:"reply_to,omitempty"`
	Domain    string `gorm:"column:domain;type:varchar(255);not null" json:"domain"`
}
