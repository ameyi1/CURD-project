package main

import (
	"fmt"
	"gin-study/api"
	"gin-study/config"
	"gin-study/middleware"
	"gin-study/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	utils.Initzap() // 初始化日志
	err := config.InitConfig()
	if err != nil {
		panic("读取配置失败：" + err.Error())
	}
	utils.InitValidator() // 初始化验证器

	// 1. 初始化数据库
	config.InitDB()

	// 2. 创建gin引擎
	r := gin.Default()

	// 3. 注册全局跨域中间件
	r.Use(middleware.GlobalRecovery()) // 注册全局异常处理
	r.Use(middleware.Cors())

	// 4. 路由分组 + 注册用户接口
	apiGroup := r.Group("/api")
	// 给整个/api分组添加中间件：所有接口必须登录
	apiGroup.Use(middleware.AuthMiddleware())

	api.RegisterUserRoutes(apiGroup)
	api.RegisterArticleRoutes(apiGroup)

	// 5. 启动服务
	Port := config.GlobalConfig.Server.Port
	utils.Log.Info("后端服务完成", zap.Int("listen_port", Port))
	err1 := r.Run(fmt.Sprintf(":%d", Port))
	if err1 != nil {
		utils.Log.Panic("启动错误", zap.Error(err1))
	}
}
