package middleware

import (
	"gin-study/config"
	"gin-study/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func Idempotency() gin.HandlerFunc {
	return func(c *gin.Context) {
		//从请求头获取token,赋值
		token := c.GetHeader("X-Indempotency-Token")
		//如果为空就踢出
		if token == "" {
			utils.Fail(c, "token为空，请重试", utils.CodeServerError)
			c.Abort()
			return

		}
		// 然后执行原始性取出并删除getdel
		ctx := c.Request.Context()
		_, err := config.RedisClient.GetDel(ctx, "indempotency"+token).Result()
		if err == redis.Nil { //这种情况可能是重复提交或者token过期了
			utils.Fail(c, "token可能过期或者重复提交", utils.CodeServerError)
			c.Abort()
			return
		} else if err != nil {
			utils.Log.Error("redis操作失败-幂等操作失败", zap.Error(err))
			utils.Fail(c, "服务器发生错误", utils.CodeServerError)
			c.Abort()
			return

		} //没问题就放行
		c.Next()
	}
}
