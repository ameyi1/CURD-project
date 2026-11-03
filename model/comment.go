package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	ArticleID uint   `json:"articleID" gorm:"index;not null"`
	UserID    uint   `json:"userID" gorm:"index;not null"`
	ParentID  uint   `json:"parentID" gorm:"index;default:0"`
	Content   string `json:"content"   gorm:"type:text;not null"`
	LikeCount uint   `json:"like_count" gorm:"default:0"`
}
