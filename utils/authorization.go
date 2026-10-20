package utils

//越权访问拦截
import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CheckOwnership(c *gin.Context, resourceOwnerID uint) bool {
	//先获取请求用户的id
	currentUserID := c.GetUint("userID")
	//c然后比较这个id和传递的id一样不一样
	if currentUserID != resourceOwnerID {
		Log.Warn("越权访问拦截",
			zap.Uint("currentUserID", currentUserID),
			zap.Uint("target_owner", resourceOwnerID),
			zap.String("path", c.Request.URL.Path))
		Fail(c, "无权操作该资源", CodeForbidden)
		return false
	}
	return true
}
