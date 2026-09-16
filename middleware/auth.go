package middleware

import (
	"strings"

	"gin-study/utils"

	"github.com/gin-gonic/gin"
)

// 这个功能是解决密码验证后续操作，比如阻止你进入什么的
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//1.从请求头获取token
		token := c.GetHeader("token")
		//2.先判空
		if strings.TrimSpace(token) == "" {
			utils.Fail(c, "密码不能为空", utils.CodeUnauthorized)
			c.Abort()
			return
		}
		//3根据写的那个Parsetoken函数来判断密码是否过期或伪造
		UserID, err := utils.ParseToken(token)
		// 解析失败：Token过期/无效
		if err != nil {
			utils.Fail(c, "登录已过期或无效", utils.CodeUnauthorized)
			c.Abort()
			return
		}

		// 4. 验证成功：把用户ID存起来，放行

		c.Set("userID", UserID)
		c.Next() // 继续访问接口

	}
}
