package middleware_test

import (
	"bytes"
	"context"
	"encoding/json"
	"gin-study/api"
	"gin-study/config"
	"gin-study/middleware"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gin-study/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func TestIdempotency(t *testing.T) {
	config.InitTestDB()
	utils.Initzap()
	gin.SetMode(gin.TestMode)

	r := gin.New()
	// 为测试添加一个简单的登录中间件，设置 userID，跳过真实鉴权
	r.POST("/articles", func(c *gin.Context) {
		c.Set("userID", uint(1))
		c.Next()
	}, middleware.Idempotency(), api.CreateArticle)

	// 1. 先在 Redis 里造一个 Token
	token := "test-token-123"
	config.RedisSet(context.Background(), "idempotency:"+token, "1", time.Minute)

	// 2. 第一次请求，应该成功
	body := bytes.NewBufferString(`{"title":"test","content":"hello"}`)
	req1, _ := http.NewRequest("POST", "/articles", body)
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("X-Idempotency-Token", token)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	var resp1 map[string]interface{}
	json.Unmarshal(w1.Body.Bytes(), &resp1)
	assert.Equal(t, float64(200), resp1["code"])

	// 3. 第二次用同一个 Token，应该被拒绝
	body2 := bytes.NewBufferString(`{"title":"test","content":"hello"}`)
	req2, _ := http.NewRequest("POST", "/articles", body2)
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("X-Idempotency-Token", token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	var resp2 map[string]interface{}
	json.Unmarshal(w2.Body.Bytes(), &resp2)
	assert.Equal(t, float64(403), resp2["code"]) // 应该被拒绝
}
