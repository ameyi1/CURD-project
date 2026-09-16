package api

import (
	"encoding/json"
	"gin-study/config"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetuser_NotFound(t *testing.T) {
	//1准备阶段=========================
	config.InitTestDB()
	//让gin不输出调试日志
	gin.SetMode(gin.TestMode)
	r := gin.New()              //创建一个干净的引擎
	r.GET("/user/:id", GetUser) //注册测试路由
	//2构造请求==============================
	req, err := http.NewRequest("GET", "/user/99999", nil)
	if err != nil {
		t.Fatalf("构建请求失败：%v", err)
	}
	//3执行请求=========================
	//创建一个假响应的记录器，用来记录handler的响应
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	//4断言===========================
	//先看状态码
	if w.Code != http.StatusOK {
		t.Errorf("期望HTTP状态码为200，实际为:%d", w.Code)
	}
	//解析响应体
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Errorf("解析响应体失败：%v", err)
	}
	// 4.3 业务码应该是 404
	if resp["code"].(float64) != 404 {
		t.Errorf("期望业务码 404，实际 %v", resp["code"])
	}
	// 4.4 消息应该包含“不存在”
	msg := resp["msg"].(string)
	if !strings.Contains(msg, "不存在") {
		t.Errorf("期望消息包含'不存在'，实际 %s", msg)
	}

}
