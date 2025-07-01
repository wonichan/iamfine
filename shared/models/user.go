package models

import (
	"time"

	"gorm.io/gorm"

	kitexuser "hupu/kitex_gen/user"
)

type UserStatus int32

const (
	UserStatusActive    UserStatus = 0
	UserStatusInactive  UserStatus = 1
	UserStatusBanned    UserStatus = 2
	UserStatusSuspended UserStatus = 3
)

type AgeGroup int32

const (
	AgeGroupUnder18 AgeGroup = 0
	AgeGroup18To25  AgeGroup = 1
	AgeGroup26To35  AgeGroup = 2
	AgeGroup36To45  AgeGroup = 3
	AgeGroup46To55  AgeGroup = 4
	AgeGroupOver55  AgeGroup = 5
)

type AnonymousAvatar struct {
	ID       string `gorm:"primaryKey;type:varchar(32)" json:"id"`
	UserID   string `gorm:"type:varchar(32);not null;index" json:"user_id"`
	Name     string `gorm:"type:varchar(50);not null" json:"name"`
	Color    string `gorm:"type:varchar(20);not null" json:"color"`
	IsActive bool   `gorm:"default:true" json:"is_active"`
}

func (AnonymousAvatar) TableName() string {
	return "anonymous_avatars"
}

type User struct {
	ID                 string                        `gorm:"primaryKey;type:varchar(32)" json:"id"`
	Role               string                        `gorm:"type:varchar(20);default:'user'" json:"role"`
	Username           string                        `gorm:"uniqueIndex;type:varchar(50);not null" json:"username"`
	Password           string                        `gorm:"type:varchar(255);not null" json:"password"`
	Nickname           string                        `gorm:"type:varchar(50)" json:"nickname"`
	Avatar             string                        `gorm:"type:varchar(255)" json:"avatar"`
	Phone              string                        `gorm:"ype:varchar(20)" json:"phone"`
	Email              string                        `gorm:"type:varchar(100)" json:"email"`
	Status             UserStatus                    `gorm:"type:int;default:0" json:"status"`
	Bio                *string                       `gorm:"type:text" json:"bio"`
	RelationshipStatus *kitexuser.RelationshipStatus `gorm:"type:int" json:"relationship_status"`
	AgeGroup           *AgeGroup                     `gorm:"type:int" json:"age_group"`
	Location           *string                       `gorm:"type:varchar(255)" json:"location"`
	PostCount          int32                         `gorm:"default:0" json:"post_count"`
	CommentCount       int32                         `gorm:"default:0" json:"comment_count"`
	LikeCount          int32                         `gorm:"default:0" json:"like_count"`
	FavoriteCount      int32                         `gorm:"default:0" json:"favorite_count"`
	AverageScore       float64                       `gorm:"type:decimal(3,2);default:0" json:"average_score"`
	FollowerCount      int32                         `gorm:"default:0" json:"follower_count"`
	FollowingCount     int32                         `gorm:"default:0" json:"following_count"`
	IsVerified         bool                          `gorm:"default:false" json:"is_verified"`
	Tags               StringArray                   `gorm:"type:json" json:"tags"`
	CreatedAt          time.Time                     `json:"created_at"`
	UpdatedAt          time.Time                     `json:"updated_at"`
	DeletedAt          gorm.DeletedAt                `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}
