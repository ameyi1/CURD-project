package model

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	Title   string `json:"title"`
	Context string `json:"context"`
	UserID  uint   `json:"userid"`
}
