package models

import (
	"time"

	"gorm.io/gorm"
)

type UserStats struct {
	ID             string         `gorm:"primaryKey;type:varchar(32)" json:"id"`
	UserID         string         `gorm:"type:varchar(32);not null;uniqueIndex" json:"user_id"`
	PostCount      int32          `gorm:"default:0" json:"post_count"`
	CommentCount   int32          `gorm:"default:0" json:"comment_count"`
	LikeCount      int32          `gorm:"default:0" json:"like_count"`
	FavoriteCount  int32          `gorm:"default:0" json:"favorite_count"`
	AverageScore   float64        `gorm:"type:decimal(3,2);default:0" json:"average_score"`
	FollowerCount  int32          `gorm:"default:0" json:"follower_count"`
	FollowingCount int32          `gorm:"default:0" json:"following_count"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (UserStats) TableName() string {
	return "user_stats"
}
