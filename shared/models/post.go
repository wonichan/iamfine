package models

import (
	"time"

	"gorm.io/gorm"

	kitex_gen_post "hupu/kitex_gen/post"
)

type PostCategory uint8

const (
	PostCategoryGeneral PostCategory = 0
	PostCategoryTech    PostCategory = 1
	PostCategorySports  PostCategory = 2
	PostCategoryLife    PostCategory = 3
	PostCategoryGaming  PostCategory = 4
	PostCategoryOther   PostCategory = 5
)

type Topic struct {
	ID          string    `gorm:"primaryKey;type:varchar(32)" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	Description *string   `gorm:"type:text" json:"description"`
	PostCount   int32     `gorm:"default:0" json:"post_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Topic) TableName() string {
	return "topics"
}

type Post struct {
	ID            string         `gorm:"primaryKey;type:varchar(32)" json:"id"`
	UserID        string         `gorm:"type:varchar(32);not null;index:idx_user_created" json:"user_id"`
	TopicID       *string        `gorm:"type:varchar(32);index" json:"topic_id"`
	Title         string         `gorm:"type:varchar(255);not null" json:"title"`
	Content       string         `gorm:"type:text;not null" json:"content"`
	Images        StringArray    `gorm:"type:json" json:"images"`
	Category      PostCategory   `gorm:"type:int;default:0;index:idx_category_created,idx_category_hot" json:"category"`
	IsAnonymous   bool           `gorm:"default:false;index:idx_anonymous_created" json:"is_anonymous"`
	AnonymousName *string        `gorm:"type:varchar(50)" json:"anonymous_name"`
	LikeCount     int32          `gorm:"default:0;index:idx_category_hot" json:"like_count"`
	CommentCount  int32          `gorm:"default:0;index:idx_category_hot" json:"comment_count"`
	ViewCount     int32          `gorm:"default:0" json:"view_count"`
	ShareCount    int32          `gorm:"default:0" json:"share_count"`
	FavoriteCount int32          `gorm:"default:0" json:"favorite_count"`
	IsHot         bool           `gorm:"default:false;index" json:"is_hot"`
	IsTop         bool           `gorm:"default:false;index" json:"is_top"`
	Location      *string        `gorm:"type:varchar(255)" json:"location"`
	Tags          StringArray    `gorm:"type:json" json:"tags"`
	CreatedAt     time.Time      `gorm:"index:idx_created_category,idx_user_created,idx_category_created,idx_anonymous_created,idx_category_hot" json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Post) TableName() string {
	return "posts"
}

func PostToKitexPost(post *Post) *kitex_gen_post.Post {
	return &kitex_gen_post.Post{
		Id:            post.ID,
		UserId:        post.UserID,
		TopicId:       post.TopicID,
		Title:         post.Title,
		Content:       post.Content,
		Images:        post.Images,
		Category:      kitex_gen_post.PostCategory(post.Category),
		IsAnonymous:   post.IsAnonymous,
		AnonymousName: post.AnonymousName,
		LikeCount:     post.LikeCount,
		CommentCount:  post.CommentCount,
		ViewCount:     post.ViewCount,
		ShareCount:    post.ShareCount,
		CollectCount:  post.FavoriteCount,
		IsHot:         post.IsHot,
		IsTop:         post.IsTop,
		Location:      post.Location,
		Tags:          post.Tags,
		CreatedAt:     post.CreatedAt.Unix(),
		UpdatedAt:     post.UpdatedAt.Unix(),
	}
}

type PostFavorite struct {
	ID        string         `gorm:"primaryKey;type:varchar(32)" json:"id"`
	UserID    string         `gorm:"type:varchar(32);not null;index:idx_user_created;uniqueIndex:idx_user_post_favorite" json:"user_id"`
	PostID    string         `gorm:"type:varchar(32);not null;index;uniqueIndex:idx_user_post_favorite" json:"post_id"`
	CreatedAt time.Time      `gorm:"index:idx_user_created" json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (PostFavorite) TableName() string {
	return "post_favorites"
}

type PostRating struct {
	ID        string         `gorm:"primaryKey;type:varchar(32)" json:"id"`
	UserID    string         `gorm:"type:varchar(32);not null;index;uniqueIndex:idx_user_post_rating" json:"user_id"`
	PostID    string         `gorm:"type:varchar(32);not null;index:idx_post_score,idx_post_created;uniqueIndex:idx_user_post_rating" json:"post_id"`
	Score     int32          `gorm:"not null;index:idx_post_score" json:"score"`
	Comment   *string        `gorm:"type:text" json:"comment"`
	CreatedAt time.Time      `gorm:"index:idx_post_created" json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (PostRating) TableName() string {
	return "post_ratings"
}
