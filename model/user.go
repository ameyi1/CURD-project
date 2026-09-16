package model

import "gorm.io/gorm"

type User struct {
	//User 模型主键冲突（数据模型从根上坏了）p0错误啊，疯了
	//Username string `json:"username" ... gorm:"primaryKey"`
	//gorm.Model 里已经有 ID uint gorm:"primaryKey"，
	// 你又给 Username 加了 gorm:"primaryKey"。GORM 会把它俩合成联合主键 (ID, Username),然后userid不再唯一
	//还没有唯一索引，甚至能注册个一模一样的用户
	//===================================================================================
	//	Password string    `json:"password" binding:"required"`
	//密码哈希被直接返回给前端，又一个p0，
	gorm.Model
	Username string    `json:"username" binding:"required,min=4,max=16" gorm:"type:varchar(32);uniqueIndex;not null"`
	Password string    `json:"-" binding:"required"`
	Name     string    `json:"name" binding:"required"`
	Age      int       `json:"age" binding:"gte=0,lte=150"`
	Email    string    `json:"email" binding:"email"`
	Articles []Article `json:"articles"`
}

//type User struct {
//	ID       uint   `json:"id" gorm:"primaryKey"`
//	Username string `json:"username" binding:"required,min=4,max=16"`
//Password string `json:"-" binding:"required,min=6,max=20"`
//}
