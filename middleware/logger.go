package middleware

import (
	"gin-study/utils"

	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取时间
		starttime := time.Now()
		c.Next()
		//获取事件后再计算耗时
		latency := time.Since(starttime)
		//获取客户端的IP
		clientIP := c.ClientIP()
		//获取状态码!
		statusCode := c.Writer.Status()
		//记录日志,http_request,方法，路径，客户端ip，状态码，耗时
		utils.Log.Info("http_Request",
			zap.String("Method", c.Request.Method),
			zap.String("Path", c.Request.URL.Path),
			zap.String("clientIP", clientIP),
			zap.Int("statusCode", statusCode),
			zap.Int64("latency_ms", latency.Microseconds()),
		)
	}
}
