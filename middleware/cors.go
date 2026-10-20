package middleware

import "github.com/gin-gonic/gin"

// 跨域中间件
//纯抄的一点也不会2026-4-12
//主要解决浏览器同源策略，前端后端loaclhost位置不同，浏览器禁止这个，设置cors响应头，告诉浏览器
// 允许这个跨域，然后options请求直接放行，不过有个*alloworigin允许任何域名需要注意下2026-10-12

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,token")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	}
}
