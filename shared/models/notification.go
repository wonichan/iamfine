package models

import (
	"gorm.io/gorm"
	"time"
)

type Notification struct {
	ID        string         `gorm:"primaryKey;type:varchar(32)" json:"id"`
	UserID    string         `gorm:"type:varchar(32);not null;index" json:"user_id"`
	Title     string         `gorm:"type:varchar(255);not null" json:"title"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	Type      int            `gorm:"not null;index" json:"type"` // 1: like, 2: comment, 3: follow
	TargetID  string         `gorm:"type:varchar(32);not null" json:"target_id"`
	IsRead    bool           `gorm:"default:false" json:"is_read"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Notification) TableName() string {
	return "notifications"
}
