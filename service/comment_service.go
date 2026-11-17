package service

import (
	"context"
	"errors"
	"fmt"
	"gin-study/dao"
	"gin-study/model"
	"gin-study/utils"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func CreateComment(ctx context.Context, articleID, userID uint, content string, parentID uint) (*model.Comment, error) {
	//1先确定文章存不存在
	_, err := dao.GetArticleByID(ctx, articleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.NewAppError(utils.CodeNotFound, "文章不存在")

		}
		return nil, fmt.Errorf("查看评论时查询的文章不存在(articleID=%d,userID=%d):%w", articleID, userID, err)

	}
	//2确认父评论是否合法
	if parentID != 0 {
		//父评论存在吗
		parent, err := dao.GetCommentByID(ctx, parentID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, utils.NewAppError(utils.CodeNotFound, "父评论不存在")
			}
			return nil, err
		}
		//父评论属于同一篇文章吗
		if parent.ArticleID != articleID {
			return nil, utils.NewAppError(utils.CodeInvalidParam, "父评论不属于这篇文章")
		}
		//父评论是是主评论吗
		if parent.ParentID != 0 {
			return nil, utils.NewAppError(utils.CodeInvalidParam, "不能回复回复，只能回复主评论")

		}
	}

	comment := &model.Comment{
		ArticleID: articleID,
		UserID:    userID,
		ParentID:  parentID,
		Content:   content,
	}
	if err := dao.CreateComment(ctx, comment); err != nil {
		return nil, err
	}
	if err := dao.IncrCommentCount(ctx, articleID, 1); err != nil {
		//  计数失败不影响评论本身
		utils.Log.Warn("评论数+1失败", zap.Uint("article_id", articleID), zap.Error(err))
	}
	return comment, nil
}

type CommentTree struct {
	model.Comment
	Replies []model.Comment `json:"replies"` //二级评论
}

// 列表
func ListComment(ctx context.Context, articleID uint) ([]CommentTree, error) {
	//获取一级评论
	parents, err := dao.GetRootComment(ctx, articleID)
	if err != nil {
		return nil, fmt.Errorf("获取一级评论失败(articleID=%v):%w", articleID, err)
	}
	if len(parents) == 0 {
		return []CommentTree{}, nil
	}
	//遍历获取一级评论的ID
	var parentIDs []uint
	for _, r := range parents {
		parentIDs = append(parentIDs, r.ID)
	}
	//根据一级评论ID查二级评论
	replies, err := dao.GetRepliesComment(ctx, parentIDs)
	if err != nil {
		return nil, fmt.Errorf("获取二级评论失败(rootID=%v):%w", parentIDs, err)
	}
	//%d 只能格式化整数，[]uint 会输出 %!d([]uint=[1 2 3]),用%v

	//将二级评论按parentID分组组装
	replyMap := make(map[uint][]model.Comment)
	for _, reply := range replies {
		replyMap[reply.ParentID] = append(replyMap[reply.ParentID], reply)
	}

	result := make([]CommentTree, 0, len(parents)) // 用 make，不用 var,不用内存优化分配内存
	for _, parent := range parents {
		rs := replyMap[parent.ID]
		if rs == nil {
			rs = []model.Comment{} //  nil → 空切片，防前端通过json。mashal导致切片的nil变为null
		}
		result = append(result, CommentTree{
			Comment: parent,
			Replies: rs,
		})
	}

	return result, nil
}
func DeleteComment(ctx context.Context, parentID, userID uint) error {
	//先确定评论存不存在
	comment, err := dao.GetCommentByID(ctx, parentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.NewAppError(utils.CodeNotFound, "评论不存在")
		}
		return fmt.Errorf("查询评论失败(commentID=%d):%w", parentID, err)
	}
	//在判断归属问题
	if comment.UserID != userID {
		return utils.NewAppError(utils.CodeForbidden, "只能删除自己的评论")
	}
	if err := dao.DeleteComment(ctx, parentID); err != nil {
		return fmt.Errorf("删除评论失败(commentID=%d):%w", parentID, err)
	}
	if err := dao.IncrCommentCount(ctx, comment.ArticleID, -1); err != nil {
		utils.Log.Warn("评论数-1失败", zap.Uint("article_id", comment.ArticleID), zap.Error(err))
	}
	return nil
}
