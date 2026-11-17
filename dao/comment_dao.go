package dao

import (
	"context"

	"gin-study/config"
	"gin-study/model"

	"gorm.io/gorm"
)

func CreateComment(ctx context.Context, comment *model.Comment) error {
	return config.DB.WithContext(ctx).Create(comment).Error
}

func GetCommentByID(ctx context.Context, commentID uint) (*model.Comment, error) {
	var comment model.Comment

	err := config.DB.WithContext(ctx).First(&comment, commentID).Error

	return &comment, err
}
func IncrCommentCount(ctx context.Context, articleID uint, delta int) error {
	return config.DB.WithContext(ctx).
		Model(&model.Article{}).
		Where("id = ?", articleID).
		UpdateColumn("comment_count", gorm.Expr("comment_count + ?", delta)).
		Error
}
func GetRootComment(ctx context.Context, articleID uint) ([]model.Comment, error) {
	var comment []model.Comment
	err := config.DB.WithContext(ctx).Where("article_id=? AND parent_id=0", articleID).Order("created_at DESC").Find(&comment).Error
	return comment, err

}
func GetRepliesComment(ctx context.Context, parentID []uint) ([]model.Comment, error) {
	var replies []model.Comment
	if len(parentID) == 0 {
		return replies, nil
	}
	err := config.DB.WithContext(ctx).Where("parent_id IN ?", parentID).Order("created_at ASC").Find(&replies).Error
	return replies, err
}

func DeleteComment(ctx context.Context, parentID uint) error {
	return config.DB.WithContext(ctx).Delete(&model.Comment{}, parentID).Error
}
