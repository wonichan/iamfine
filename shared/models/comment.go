package models

import (
	"gorm.io/gorm"
	"time"
)

type Comment struct {
	ID             string         `gorm:"primaryKey;type:varchar(32)" json:"id"`
	PostID         string         `gorm:"type:varchar(32);not null;index" json:"post_id"`
	UserID         string         `gorm:"type:varchar(32);not null;index" json:"user_id"`
	ParentID       *string        `gorm:"type:varchar(32);index" json:"parent_id"`
	Content        string         `gorm:"type:text;not null" json:"content"`
	LikeCount      int32          `gorm:"default:0" json:"like_count"`
	ReplyCount     int32          `gorm:"default:0" json:"reply_count"`
	IsAnonymous    bool           `gorm:"default:false" json:"is_anonymous"`
	AnonymousName  *string        `gorm:"type:varchar(50)" json:"anonymous_name"`
	AnonymousColor *string        `gorm:"type:varchar(20)" json:"anonymous_color"`
	Images         StringArray    `gorm:"type:json" json:"images"`
	IsDeleted      bool           `gorm:"default:false" json:"is_deleted"`
	Location       *string        `gorm:"type:varchar(255)" json:"location"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Comment) TableName() string {
	return "comments"
}
