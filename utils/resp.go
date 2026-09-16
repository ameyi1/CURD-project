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
		Code: CodeSuccess,
		Msg:  msg,
		Data: data,
	})
}

// 失败返回
// 兼容历史调用方式：Fail(c, "msg") 和 Fail(c, "msg", 401)
func Fail(c *gin.Context, msg string, code ...int) {
	errCode := CodeServerError
	if len(code) > 0 {
		errCode = code[0]
	}
	c.JSON(200, Response{
		Code: errCode,
		Msg:  msg,
		Data: nil,
	})
}
