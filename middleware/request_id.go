package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

type requestIDkey struct{}

// 永远不要用 nil 作为 Context 参数。如果你不确定用什么，用 context.TODO() 代替。
// 因为 nil context 会导致 panic，而且无法判断类型。
func RequsetID() gin.HandlerFunc {
	return func(c *gin.Context) {
		//从头请求获得requestID
		requestID := c.GetHeader("X-Requset-ID")
		//没有就创建
		if requestID == "" {
			requestID = GenerateRequestID()
		}
		//把reuquestId放进context
		ctx := context.WithValue(c.Request.Context(), requestIDkey{}, requestID)
		//替换发起真正请求的context，就是那个可能没有名称的
		//这句是基于旧的 Request，生成一个新的 Request，但这个新 Request 的关联 Context 被替换成你指定的 Context
		c.Request = c.Request.WithContext(ctx)
		//把 Request ID 写入响应头c
		c.Header("X-Request-ID", requestID)
		c.Next()

	}
}

// 生成请求名称的
func GenerateRequestID() string {
	//生成一个byte 长16的切片
	b := make([]byte, 16)
	//然后随机填
	rand.Read(b)
	//把字节转为string
	return hex.EncodeToString(b)

}

// 为了方便写这个函数，后续只要middle.GetRequsetID(ctx)就好了!!!
func GetRequsetID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	//!!!ctx.Value(requestIDKey{})：从 context 里取之前放入的值 (string)：类型断言。因为我们知道放进去的是 string，所以要断言成字符串类型。
	//ok 判断：如果 context 里没有这个 key 或值类型不对，ok 会是 false，我们返回空字符串
	if RequsetID, ok := ctx.Value(requestIDkey{}).(string); ok {
		return RequsetID

	}
	return ""
}
