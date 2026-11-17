package model

import "gorm.io/gorm"

type Article struct { //自带id不用在家articleID了
	gorm.Model
	Title        string `json:"title" gorm:"index"`
	Context      string `json:"context"`
	UserID       uint   `json:"user_id" gorm:"index"`
	LikeCount    uint   `json:"like_count" gorm:"default:0"`
	CommentCount uint   `json:"comment_count" gorm:"default:0"`
}
