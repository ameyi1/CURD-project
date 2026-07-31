package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string    `json:"username" binding:"required,min=4,max=16" gorm:"primaryKey"`
	Password string    `json:"password" binding:"required"`
	Name     string    `json:"name" binding:"required"`
	Age      int       `json:"age" binding:"gte=0,lte=120"`
	Email    string    `json:"email" binding:"email"`
	Articles []Article `json:"articles"`
}

//type User struct {
//	ID       uint   `json:"id" gorm:"primaryKey"`
//	Username string `json:"username" binding:"required,min=4,max=16"`
//Password string `json:"-" binding:"required,min=6,max=20"`
//}
