package dao

import (
	"gin-study/config"
	"gin-study/model"
)

// /CreateUser 插入用户记录，返回错误
func CreateUser(user *model.User) error {
	return config.DB.Create(user).Error
}

// GetUserByID 根据 ID 查询用户
func GetUserByID(id uint) (*model.User, error) {
	var user model.User
	err := config.DB.First(&user, id).Error
	return &user, err
}
