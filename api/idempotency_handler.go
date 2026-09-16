package api

//幂等token的生成
import (
	"gin-study/config"
	"gin-study/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func GetindempotencyToken(c *gin.Context) {
	//先生成个token
	//uuid.New().String()：生成一个全局唯一的字符串，比如 550e8400-e29b-41d4-a716-446655440000
	token := uuid.New().String()
	//生成完token后纯净redis里设置过期时间5min
	ctx := c.Request.Context()
	err := config.RedisClient.Set(ctx, "indempotency"+token, "1", 5*time.Minute).Err()
	if err != nil {
		utils.Fail(c, "存储indempotencytoken到redis失败", utils.CodeServerError)
		utils.Log.Error("存储indempotencytoken到redis失败", zap.Error(err))
		return
	}
	//然后返回给前端
	utils.Success(c, gin.H{"token": token}, "前端获取indempotencytoken成功")

}
