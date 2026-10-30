package model

import "gorm.io/gorm"

type Article struct { //自带id不用在家articleID了
	gorm.Model
	Title     string `json:"title" gorm:"index"`
	Context   string `json:"context"`
	UserID    uint   `json:"userid" gorm:"index"`
	LikeCount uint   `json:"likecount" gorm:"default:0"`
}
