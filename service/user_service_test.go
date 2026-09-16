package service

import (
	"gin-study/config"
	"gin-study/model"
	"testing"
)

func TestCreateUser(t *testing.T) {
	config.InitTestDB()
	user, err := CreateUser("testuser", "123456", "测试用户", "test@example.com", 18)
	if err != nil {
		t.Fatalf("创建用户失败：%v", err)
	}
	if user.Username != "testuser" {
		t.Errorf("期望用户名为testuser,实际为:%s", user.Username)
	}
	if user.Password == "123456" {
		t.Error("密码不应该为明文")
	}
	var count int64
	config.DB.Model(&model.User{}).Where("username=?", "testuser").Count(&count)
	if count != 1 {
		t.Errorf("期望数据库有一条用户记录，实际为:%d", count)
	}
}
