package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string         `gorm:"primaryKey;type:varchar(32)" json:"id"`
	Username  string         `gorm:"uniqueIndex;type:varchar(50);not null" json:"username"`
	Password  string         `gorm:"type:varchar(255);not null" json:"password"`
	Nickname  string         `gorm:"type:varchar(50)" json:"nickname"`
	Avatar    string         `gorm:"type:varchar(255)" json:"avatar"`
	Phone     string         `gorm:"uniqueIndex;type:varchar(20)" json:"phone"`
	Email     string         `gorm:"type:varchar(100)" json:"email"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}
