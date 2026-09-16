package main

import (
	"context"
	"errors"
	"fmt"
	"gin-study/api"
	"gin-study/config"
	"gin-study/middleware"
	"gin-study/utils"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	//config.InitConfig(),初始化两次的来了
	config.InitRedis()
	// 2. 创建gin引擎
	r := gin.Default()

	// 3. 注册全局跨域中间件
	r.Use(middleware.GlobalRecovery()) // 注册全局异常处理
	r.Use(middleware.Cors())
	//获取日志requestID
	r.Use(middleware.RequsetID())
	r.Use(middleware.Logger())

	// 4. 路由分组 + 注册用户接口
	apiGroup := r.Group("/api")
	// 给整个/api分组添加中间件：所有接口必须登录

	//!!!!!!!!!!!!!!!!!!!!!!!p0级别错误，所有接口都需要token包括登录，造成无限循环|
	//apiGroup.Use(middleware.AuthMiddleware())                            +  !  +
	//这个是限流
	apiGroup.POST("/login", middleware.RateLimit(5, time.Minute), api.GetUser)
	apiGroup.POST("/user", middleware.RateLimit(5, time.Minute), api.CreateNewUser)
	apiGroup.POST("/article/create", middleware.Idempotency(), api.CreateArticle)
	apiGroup.POST("/refreash", utils.Refrashtoken)
	// 其他接口
	authgroup := r.Group("")
	authgroup.Use(middleware.AuthMiddleware())
	api.RegisterUserRoutes(authgroup)
	//api.RegisterArticleRoutes(apiGroup)
	api.RegisterUserRoutes(authgroup)
	//api.RegisterArticleRoutes(apiGroup)

	//==================================================================================
	//优雅关机
	//http.Server 没有任何超时（会被 Slowloris 拖死）p0级别错误        这个看不懂，回头看
	server := &http.Server{
		Addr:    fmt.Sprintf("%d", config.GlobalConfig.Server.Port),
		Handler: r,
		//=======================看不懂下面的========================
		ReadHeaderTimeout: 5 * time.Second,  // 读完请求头的最长时间，防 Slowloris 的核心
		ReadTimeout:       10 * time.Second, // 读完整请求（含 body）的最长时间
		WriteTimeout:      10 * time.Second, // 写完响应的最长时间
		IdleTimeout:       60 * time.Second, // keep-alive 空闲连接存活时间
		MaxHeaderBytes:    1 << 20,          // 请求头最大 1MB
		//=======================================================================
	} //监控
	go func() {
		utils.Log.Info("后端服务启动", zap.Int("listen_port", config.GlobalConfig.Server.Port))        //不懂
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) { //这个表示排除了主动关闭之外的特殊错误，
			utils.Log.Panic("服务启动失败", zap.Error(err))

		}

	}()
	//接收系统信号
	quit := make(chan os.Signal, 1)
	//这些把SIGINT（CTRL+c)SIGTERM(kill/容器终止)信号传给channel里
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	//没信息就阻塞等待
	<-quit
	utils.Log.Info("开始尝试关闭...")
	//设置超时时长
	ctx, cancal := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancal()
	//这个是关键，go帮你写好了，不在接收新请求然后解决正在进行的请求之后关闭
	if err := server.Shutdown(ctx); err != nil {
		utils.Log.Warn("关闭失败，准备强制关闭", zap.Error(err))
		return
	}

	// 5. 启动服务
	//Port := config.GlobalConfig.Server.Port
	//utils.Log.Info("后端服务完成", zap.Int("listen_port", Port))
	//err1 := r.Run(fmt.Sprintf(":%d", Port))
	//if err1 != nil {
	//	utils.Log.Panic("启动错误", zap.Error(err1))
	//}
}
