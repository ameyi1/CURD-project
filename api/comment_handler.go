// api/comment_handler.go
package api

import (
	"strconv"

	"gin-study/service"
	"gin-study/utils"

	"github.com/gin-gonic/gin"
)

func RegisterCommentRoutes(r *gin.RouterGroup) {
	r.POST("/articles/:id/comments", CreateComment)
	r.GET("/articles/:id/comments", ListComments)
	r.DELETE("/comments/:id", DeleteComment)
}

// CreateComment 发评论
func CreateComment(c *gin.Context) {
	articleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Fail(c, "文章ID必须是数字", utils.CodeInvalidParam)
		return
	}

	userID := c.GetUint("userID")
	if userID == 0 {
		utils.Fail(c, "请先登录", utils.CodeUnauthorized)
		return
	}

	var req struct {
		Content  string `json:"content"   binding:"required,max=500"`
		ParentID uint   `json:"parent_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, "参数错误："+err.Error(), utils.CodeInvalidParam)
		return
	}

	comment, err := service.CreateComment(c.Request.Context(), uint(articleID), userID, req.Content, req.ParentID)
	if err != nil {
		utils.HandleError(c, "评论失败", err)
		return
	}
	utils.Success(c, comment, "评论成功")
}

// ListComments 查评论列表
func ListComments(c *gin.Context) {
	articleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Fail(c, "文章ID必须是数字", utils.CodeInvalidParam)
		return
	}

	list, err := service.ListComment(c.Request.Context(), uint(articleID))
	if err != nil {
		utils.HandleError(c, "查询评论失败", err)
		return
	}
	utils.Success(c, list, "查询成功")
}

// DeleteComment 删评论
func DeleteComment(c *gin.Context) {
	commentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Fail(c, "评论ID必须是数字", utils.CodeInvalidParam)
		return
	}

	userID := c.GetUint("userID")
	if userID == 0 {
		utils.Fail(c, "请先登录", utils.CodeUnauthorized)
		return
	}

	if err := service.DeleteComment(c.Request.Context(), uint(commentID), userID); err != nil {
		utils.HandleError(c, "删除评论失败", err)
		return
	}
	utils.Success(c, nil, "删除成功")
}
