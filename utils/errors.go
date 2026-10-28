package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	CodeSuccess = 200

	CodeInvalidParam = 400
	CodeUnauthorized = 401
	CodeForbidden    = 403
	CodeNotFound     = 404
	CodeServerError  = 500
)

// | 错误码 | 含义 | 什么时候用 |
// | :--- | :--- | :--- |
// | `200` | 成功 | 所有正常返回 |
// | `400` | 参数错误 | 请求参数格式不对、缺少必填字段、类型不对 |
// | `401` | 未认证 | 请求需要登录，但没带 token / token 无效 / token 过期 |
// | `403` | 没有权限 | 已登录，但没权限操作别人的资源（比如改别人的文章） |
// | `404` | 资源不存在 | 查的用户、文章、内容不存在 |
// | `500` | 服务器内部错误 | 数据库挂了、系统出错，不能暴露细节给前端 |
// 定义一个错误码类型的结构体
type AppError struct {
	Code int
	Msg  string
}

// 让APPerr变为一个error错误类型
func (e *AppError) Error() string {
	return e.Msg
}

// 制造一个业务错误模板service用
func NewAppError(code int, msg string) error {
	return &AppError{Code: code, Msg: msg}
}
func HandleError(c *gin.Context, defaultMsg string, err error) {
	//防御性判断
	if err == nil {
		return
	}
	var appErr *AppError
	//业务错误
	if errors.As(err, &appErr) {
		Fail(c, appErr.Msg, appErr.Code)
		return

	}
	//系统错误
	Log.Error("系统异常", zap.String("path", c.Request.URL.Path), zap.String("method", c.Request.Method), zap.Error(err))
	Fail(c, defaultMsg, CodeServerError)

}
