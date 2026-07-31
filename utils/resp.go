package utils

import "github.com/gin-gonic/gin"

// 统一返回结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// 成功返回
func Success(c *gin.Context, data interface{}, msg string) {
	c.JSON(200, Response{
		Code: 200,
		Msg:  msg,
		Data: data,
	})
}

// 失败返回
func Fail(c *gin.Context, msg string) {
	c.JSON(200, Response{
		Code: 400,
		Msg:  msg,
		Data: nil,
	})
}
