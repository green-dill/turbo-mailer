package models

// Pool
type Pool struct {
	BaseModel
	Name        string        `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Description string        `gorm:"column:description;type:varchar(512)" json:"description,omitempty"`
	SenderCount int           `gorm:"column:sender_count;type:int;default:0" json:"sender_count"`
	Senders     []*PoolSender `gorm:"foreignKey:PoolID;references:ID" json:"senders"`
}

// PoolSender
type PoolSender struct {
	BaseModel
	PoolID    uint   `gorm:"column:pool_id;type:bigint;not null" json:"pool_id"`
	FromName  string `gorm:"column:from_name;type:varchar(255);not null" json:"from_name" validate:"required"`
	FromEmail string `gorm:"column:from_email;type:varchar(255);not null" json:"from_email" validate:"required"`
	ReplyTo   string `gorm:"column:reply_to;type:varchar(255)" json:"reply_to,omitempty"`
	Domain    string `gorm:"column:domain;type:varchar(255);not null" json:"domain" validate:"required"`
}
