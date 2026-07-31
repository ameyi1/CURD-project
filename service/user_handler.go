package service

import (
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

func GetUser(id uint) (*model.User, error) {
	return dao.GetUserByID(id)
}
