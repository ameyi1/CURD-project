package api

import (
	"context"
	"errors"
	"strconv"
	"time"

	"gin-study/model"
	"gin-study/service"
	"gin-study/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RegisterUserRoutes 注册用户相关路由
func RegisterUserRoutes(r *gin.RouterGroup) {
	r.POST("/users", CreateUser) // 规范命名：POST /api/users
	r.GET("/users/:id", GetUser) // 规范命名：GET /api/users/:id
	//r.GET("/search/:id",SearchUser)//GET/api/search/:id
	r.GET("/list", GetUsers)                //GET/api/list
	r.POST("/createNewUser", CreateNewUser) //POST/api/createNewUser
}

// CreateUser 接收请求 → 调用 service → 返回响应
func CreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=4,max=16"`
		Password string `json:"password" binding:"required"`
		Name     string `json:"name" binding:"required"`
		Age      int    `json:"age" binding:"gte=0,lte=120"`
		Email    string `json:"email" binding:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, "参数错误：", utils.CodeInvalidParam)
		return
	}

	user, err := service.CreateUser(req.Username, req.Password, req.Name, req.Email, req.Age)
	if err != nil {
		var AppError *utils.AppError
		if errors.As(err, &AppError) {
			//p0错误，为什么写apperror不加&啊
			utils.Fail(c, AppError.Msg, AppError.Code)
		} else {
			utils.Log.Error("创建用户失败", zap.Error(err))
			utils.Fail(c, "创建用户失败", utils.CodeServerError)
		}
		return
	}
	utils.Success(c, user, "创建用户成功")
}
func CreateNewUser(c *gin.Context) {
	var user model.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		//统一错误处理
		var AppError *utils.AppError
		if errors.As(err, &AppError) {
			utils.Fail(c, AppError.Msg, AppError.Code)
		} else {
			utils.Log.Error("绑定错误", zap.Error(err))
			utils.Fail(c, "绑定错误", utils.CodeServerError)
		}

		return
	}
	newUser, HelloArticle, err2 := service.CreateNewUser(&user)
	if err2 != nil {
		var AppError *utils.AppError
		if errors.As(err2, &AppError) {
			utils.Fail(c, AppError.Msg, AppError.Code)
		} else {
			utils.Log.Error("创建用户失败", zap.Error(err2))
			utils.Fail(c, "创建用户失败：", utils.CodeServerError)
		}
		return
	}
	// 成功就把用户和欢迎文章一起返回给前端
	utils.Success(c, gin.H{
		"user":    newUser,
		"article": HelloArticle,
	}, "创建成功")
}

// GetUser 根据 ID 查询用户
func GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err1 := strconv.Atoi(idStr)
	if err1 != nil {
		utils.Fail(c, "ID必须是纯数字", utils.CodeForbidden)
		return
	}
	//纯浪费内存如果写ctx，defer下面
	//先解决授权问题
	if !utils.CheckOwnership(c, uint(id)) {
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	user, err2 := service.GetUser(ctx, uint(id))
	if err2 != nil {
		var AppError *utils.AppError
		if errors.As(err2, &AppError) {
			utils.Fail(c, AppError.Msg, AppError.Code)

		} else {
			utils.Log.Error("查询用户信息失败", zap.Error(err2))
			utils.Fail(c, "查询用户信息失败：", utils.CodeServerError)
		}
		return
	}

	utils.Success(c, user, "查询用户信息成功")
} //查询单个用户
/*写重复了
	func SearchUser(c *gin.Context){
		idStr:=c.Param("id")
		id,err1:=strconv.Atoi(idStr)
		if err1!=nil{
			utils.Fail(c,"查询用户失败"+err1.Error())
			return
		}
		user,err2:=service.GetUser(uint(id))
		if err2!=nil{
			utils.Fail(c,"查询失败："+err2.Error())
			return
		}
		utils.Success(c,user,"查询成功")

}
*/
//查多个用户并分页
func GetUsers(c *gin.Context) {
	var list []model.User
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 10
	}
	offset := (page - 1) * size
	list, err := service.GetUsers(offset, size)
	if err != nil {
		var AppError *utils.AppError
		if errors.As(err, &AppError) {
			utils.Fail(c, AppError.Msg, AppError.Code)
			return
		} else {
			utils.Log.Error("查询列表失败", zap.Error(err))
			utils.Fail(c, "查询列表失败", utils.CodeServerError)
			return
		}
		return
	}
	utils.Success(c, list, "查询列表成功")

}

// 登录LoginUser同样，需要获取ctx,username password,然后调用service的登录LoginUser
func LoginUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, "用户名和密码不能为空", utils.CodeInvalidParam)
		return
	}
	if len(req.Password) < 6 {
		utils.Fail(c, "密码长度不能少于6位", utils.CodeInvalidParam)
		return
	}
	accessToken, refreshToken, user, err := service.LoginUser(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		utils.Fail(c, "登录失败", utils.CodeServerError)
		return
	}
	utils.Success(c, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          user,
	}, "登录成功")
}
