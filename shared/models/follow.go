package models

import (
	"gorm.io/gorm"
	"time"
)

type Follow struct {
	ID          string         `gorm:"primaryKey;type:varchar(32)" json:"id"`
	FollowerID  string         `gorm:"type:varchar(32);not null;index" json:"follower_id"`
	FollowingID string         `gorm:"type:varchar(32);not null;index" json:"following_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Follow) TableName() string {
	return "follows"
}
