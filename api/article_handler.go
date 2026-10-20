package api

import (
	"context"
	"errors"
	"gin-study/config"
	"gin-study/model"
	"gin-study/service"
	"gin-study/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 增
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

// 删
func DeleteArticle(c *gin.Context) {
	//先从URL里面获取文章id
	idStr := c.Param("id")
	articleID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.Fail(c, "文章id必须为数字", utils.CodeInvalidParam)
		return
	}
	//再拿userid
	userIDVal, exists := c.Get("userID")
	if !exists {
		utils.Fail(c, "未登录", utils.CodeUnauthorized)
		return
	}
	//再类型断言
	userID := userIDVal.(uint)
	//建ctx,设超时时间
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	//调用service_deleteArticle
	if err := service.DeleteArticle(ctx, uint(articleID), userID); err != nil {
		var appError *utils.AppError
		//判断下是什么错误
		if errors.As(err, &appError) {
			//如何这个as的返回值bool的结果是true，传给前端，是业务错误,不用手动填了
			utils.Fail(c, appError.Msg, appError.Code)
		} else {
			//bool是fail的话就是系统错误了,日志记录
			utils.Log.Error("删除文章失败", zap.Error(err))
			utils.Fail(c, "删除文章失败", utils.CodeServerError)

		}
		return
	}
	utils.Success(c, nil, "成功删除文章")
}

// 查
func SearchArticle(c *gin.Context) {
	idStr := c.Param("id")
	articleID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.Fail(c, "文章id必须为数字", utils.CodeInvalidParam)
		return
	}

	utils.Success(c, gin.H{"id": articleID}, "查询文章成功")
}

// 点赞
func LikeArticle(c *gin.Context) {
	//获取文章id
	idstr := c.Param("ArticleID")
	articleID, err := strconv.Atoi(idstr)
	if err != nil {
		utils.Fail(c, "文章ID必须是数字", utils.CodeInvalidParam)
		return
	} //获取用户id
	userID := c.GetUint("userID")
	//调用service
	ctx := c.Request.Context()
	if err := service.LikeArticle(ctx, uint(articleID), userID); err != nil {
		var appErr *utils.AppError
		if errors.As(err, &appErr) {
			utils.Fail(c, appErr.Msg, appErr.Code)
		} else {
			utils.Log.Error("点赞失败", zap.Error(err))
			utils.Fail(c, "服务器内部错误", utils.CodeServerError)

		}
		return
	}
	//调用service的点赞计数
	count, err := service.GetLikeCount(ctx, uint(articleID))
	if err != nil {
		var appErr *utils.AppError
		if errors.As(err, &appErr) {
			utils.Fail(c, appErr.Msg, appErr.Code)
		} else {
			utils.Log.Error("查询点赞个数失败", zap.Error(err))
			utils.Fail(c, "服务器内部错误", utils.CodeServerError)

		}
		return
	}

	utils.Success(c, gin.H{"like_count": count}, "点赞成功")

} //和上面的几乎一致
func UnlikeArticle(c *gin.Context) {
	idStr := c.Param("id")
	articleID, err := strconv.Atoi(idStr)
	if err != nil {
		utils.Fail(c, "文章ID必须是数字", utils.CodeInvalidParam)
		return
	}
	userID := c.GetUint("userID")
	ctx := c.Request.Context()
	if err := service.UnlikeArticle(ctx, uint(articleID), userID); err != nil {
		var appErr *utils.AppError
		if errors.As(err, &appErr) {
			utils.Fail(c, appErr.Msg, appErr.Code)
		} else {
			utils.Log.Error("取消点赞失败", zap.Error(err))
			utils.Fail(c, "服务器内部错误", utils.CodeServerError)
		}
		return
	}
	count, err := service.GetLikeCount(ctx, uint(articleID))
	if err != nil {
		utils.Log.Error("查询点赞个数失败", zap.Error(err))
		utils.Fail(c, "服务器内部错误", utils.CodeServerError)
		return
	}
	utils.Success(c, gin.H{"like_count": count}, "取消点赞成功")
}

// 增加文章浏览量
func GetViewCount(c *gin.Context) {
	idstr := c.Param("articleID")
	articleID, err := strconv.Atoi(idstr)
	if err != nil {
		utils.Fail(c, "必须是纯数字", utils.CodeServerError)

	}
	//先查文章浏览总数
	ctx := c.Request.Context()
	count, err := service.GetViewCount(ctx, uint(articleID))
	if err != nil {
		var appErr *utils.AppError
		if !errors.As(err, *appErr) {
			utils.Log.Error("获取浏览量失败", zap.Error(err))
			utils.Fail(c, appErr.Msg, appErr.Code)
			return

		} else {
			utils.Log.Error("获取浏览量错误", zap.Error(err))
			utils.Fail(c, "获取浏览量失败", utils.CodeServerError)
			return
		}
		//然后增加浏览量的话再开一个携程吧，要不然浏览量增加还要堵人有点那啥了
		//不行，别开，有点不值得一个获取浏览量的东西，而且，context，background高并发直接会导致携程数量爆炸

	}

	service.Addview(ctx, uint(articleID))
	// 4. 获取最新浏览量
	Count, err := service.GetViewCount(ctx, uint(articleID))
	if err != nil {
		var appErr *utils.AppError
		if !errors.As(err, *appErr) {
			utils.Log.Error("获取浏览量失败", zap.Error(err))
			utils.Fail(c, appErr.Msg, appErr.Code)
			return

		} else {
			utils.Log.Error("获取浏览量错误", zap.Error(err))
			utils.Fail(c, "获取浏览量失败", utils.CodeServerError)
			return
		}
	}
	// 5. 返回
	utils.Success(c, gin.H{
		"article":       articleID,
		"Newview_count": Count,
		"oldView_count": count,
	}, "查询成功")

}

//没招了，为什么要旧数据，纯脱裤子放屁
