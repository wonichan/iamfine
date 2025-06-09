package models

import (
	"gorm.io/gorm"
	"time"
)

type Like struct {
	ID         string         `gorm:"primaryKey;type:varchar(32)" json:"id"`
	UserID     string         `gorm:"type:varchar(32);not null;index" json:"user_id"`
	TargetID   string         `gorm:"type:varchar(32);not null;index" json:"target_id"`
	TargetType int            `gorm:"not null;index" json:"target_type"` // 1: post, 2: comment
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Like) TableName() string {
	return "likes"
}
