package service

import (
	"context"
	"gin-study/dao"
	"gin-study/model"
	"gin-study/utils"

	"go.uber.org/zap"
)

func CreateUser(username, password, name, email string, age int) (*model.User, error) {
	encryptPwd, err := utils.EncryptPassword(password)
	if err != nil {
		utils.Log.Error("加密失败", zap.Error(err))
		return nil, err
	}

	user := &model.User{
		Username: username,
		Password: encryptPwd,
		Name:     name,
		Age:      age,
		Email:    email,
	}

	if err := dao.CreateUser(user); err != nil {
		utils.Log.Error("创建用户失败", zap.Error(err))
		return nil, err
	}
	return user, nil
}
func CreateNewUser(user *model.User) (*model.User, *model.Article, error) { //string类型

	encryptPwd, err := utils.EncryptPassword(user.Password)
	if err != nil {
		utils.Log.Error("密码加密失败", zap.Error(err))
		return nil, nil, err
	}
	user.Password = encryptPwd
	newuser, HelloArticle, err2 := dao.CreateNewUser(user)
	if err2 != nil {
		utils.Log.Error("创建用户失败", zap.Error(err2))
		return nil, nil, err2
	}
	return newuser, HelloArticle, nil
}
func GetUser(ctx context.Context, id uint) (*model.User, error) {
	return dao.GetUserByID(ctx, id)
}
func GetUsers(offset, size int) ([]model.User, error) {
	return dao.GetUsers(offset, size)
}
