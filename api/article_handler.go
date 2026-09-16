package api

import (
	"gin-study/config"
	"gin-study/model"
	"gin-study/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CreateArticle(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.Fail(c, "未登录或 token 无效", utils.CodeUnauthorized)
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		utils.Fail(c, "用户信息异常", utils.CodeUnauthorized)
		return
	}

	var article model.Article
	if err := c.ShouldBindJSON(&article); err != nil {
		utils.Log.Error("参数错误:", zap.Error(err))
		utils.Fail(c, "参数错误："+err.Error())
		return
	}

	article.UserID = uid
	if err := config.DB.Create(&article).Error; err != nil {
		utils.Log.Error("发布文章失败:", zap.Error(err))
		utils.Fail(c, "发布文章失败："+err.Error())
		return
	}

	utils.Success(c, article, "发布文章成功")
	utils.Log.Info("发布文章成功", zap.String("title", article.Title), zap.Uint("userID", article.UserID), zap.Uint("articleID", article.ID))
}
