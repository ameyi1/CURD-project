package middleware

import (
	"fmt"
	"gin-study/config"
	"gin-study/utils"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 限流 防止给服务器卡爆
func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		//先获取请求ip
		clientIP := c.ClientIP()
		//再设置key
		key := fmt.Sprintf("rate_limit:%s", clientIP)
		ctx := c.Request.Context()
		//redis内置的计数器
		count, err := config.RedisClient.Incr(ctx, key).Result()
		//判断情况，如果失效就驳出
		if err != nil {
			utils.Log.Warn("限流失败，放行", zap.Error(err))
			//是c.next不是c.abort

			c.Next()
			return
		} //如果是首次请求那设置过期期限
		if count == 1 {
			config.RedisClient.Expire(ctx, key, window)

		}
		//然后设置请求次数限制,限流的关键
		if count > int64(limit) {
			utils.Fail(c, "请求次数过多请稍后再试", utils.CodeServerError)
			c.Abort()
			return
		}
		c.Next()

	}
}
