package service

import (
	"context"
	"errors"
	"gin-study/dao"
	"gin-study/model"
	"gin-study/utils"

	"go.uber.org/zap"
	"gorm.io/gorm"
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

	user, err := dao.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { //is,不是==,因为==就只是表面对着，is拆链逐层找
			return nil, utils.NewAppError(utils.CodeNotFound, "用户不存在")
		}
		return nil, err
	}
	return user, nil
}
func GetUsers(offset, size int) ([]model.User, error) {
	return dao.GetUsers(offset, size)
}

//我需要写个登录功能
//用户提供账号与密码然后验证tokeen后让他进入或进行错误通知
//1输入账号，密码，上下文ctx，->2处理->3输出错误，accesstoken，refreashtoken，个人信息
//处理：1先根据账号查数据库看看有没有这个人
//2然后如果没有就错误处理
//3对比密码，然后错误处理，之前写的jwt里面有这个东西
//4输出token，返回

func LoginUser(ctx context.Context, username string, password string) (string, string, *model.User, error) {
	//1
	user, err := dao.LoginUser(ctx, username)

	//2错误处理

	if err != nil {
		// 是不是为空
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", nil, utils.NewAppError(utils.CodeUnauthorized, "账号或密码错误")
		}
		return "", "", nil, err
	} //返回同一句话防止用户名枚举漏洞
	// 3+密码是不是错了
	if !utils.ComparePassword(password, user.Username) {
		return "", "", nil, utils.NewAppError(utils.CodeUnauthorized, "账号或密码错误")

	}
	//4生成token
	access, refresh, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "", "", nil, err
	}
	return access, refresh, user, nil

}
