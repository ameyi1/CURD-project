package middleware

import (
	"gin-study/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 这个主要兜住panic，然后stack记录堆栈定位是那个panic
func GlobalRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			err := recover()
			if err != nil {
				utils.Log.Error("接口运行发生panic崩溃",
					zap.Any("err", err),
					zap.String("request_way", c.FullPath()),
					zap.String("request_method", c.Request.Method),
					zap.String("clientIP", c.ClientIP()),
					zap.Stack("stack"),
				)
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": 500,
					"msg":  "服务内部异常,请稍后再试",
					"data": nil,
				})
				c.Abort()
			}

		}()
		c.Next()

	}
}
