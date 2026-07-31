package api

import (
	"strconv"

	"gin-study/service"
	"gin-study/utils"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes 注册用户相关路由
func registerUserRoutes(r *gin.RouterGroup) {
	r.POST("/users", CreateUser) // 规范命名：POST /api/users
	r.GET("/users/:id", GetUser) // 规范命名：GET /api/users/:id
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
		utils.Fail(c, "参数错误："+err.Error())
		return
	}

	user, err := service.CreateUser(req.Username, req.Password, req.Name, req.Email, req.Age)
	if err != nil {
		utils.Fail(c, "创建失败")
		return
	}
	utils.Success(c, user, "创建成功")
}

// GetUser 根据 ID 查询用户
func GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.Fail(c, "ID必须是纯数字")
		return
	}

	user, err := service.GetUser(uint(id))
	if err != nil {
		utils.Fail(c, "查询失败："+err.Error())
		return
	}

	utils.Success(c, user, "查询成功")
}
