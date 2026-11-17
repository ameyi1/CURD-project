package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	ArticleID uint   `json:"article_id" gorm:"index;not null"`
	UserID    uint   `json:"user_id" gorm:"index;not null"`
	ParentID  uint   `json:"parent_id" gorm:"index;default:0"`
	Content   string `json:"content"   gorm:"type:text;not null"`
	LikeCount uint   `json:"like_count" gorm:"default:0"`
}
