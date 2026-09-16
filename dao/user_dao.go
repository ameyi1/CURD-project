package dao

import (
	"context"
	"gin-study/config"
	"gin-study/model"
	"gin-study/utils"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// /CreateUser 插入用户记录，返回错误
func CreateUser(user *model.User) error {
	return config.DB.Create(user).Error
}
func CreateNewUser(user *model.User) (*model.User, *model.Article, error) {
	var HelloArticle model.Article
	txErr := config.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(user)
		if err.Error != nil {
			utils.Log.Warn("创建用户失败", zap.Error(err.Error))
			return err.Error
		}
		HelloArticle = model.Article{
			Title:   "Hello",
			Context: "这是一篇欢迎文章",
			UserID:  user.ID,
		}
		err2 := tx.Create(&HelloArticle)
		if err2.Error != nil {
			utils.Log.Warn("创建欢迎文章失败", zap.Error(err2.Error))
			return err2.Error
		}
		return nil
	})
	if txErr != nil {
		utils.Log.Warn("启动事务失败", zap.Error(txErr))
		return nil, nil, txErr
	}
	return user, &HelloArticle, nil
}

// GetUserByID 根据 ID 查询用户
func GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := config.DB.WithContext(ctx).First(&user, id).Error
	return &user, err
}
func GetUsers(offset int, size int) ([]model.User, error) {
	var list []model.User
	err := config.DB.Offset(offset).Limit(size).Find(&list)
	if err.Error != nil {
		utils.Log.Warn("查询用户列表失败", zap.Int("offset", offset), zap.Int("size", size), zap.Error(err.Error))
		return nil, err.Error
	}
	return list, nil
}
