package dao

import (
	"context"

	"gin-study/config"
	"gin-study/model"
)

func CreateArticle(article *model.Article) error {
	return config.DB.Create(article).Error
}

func GetArticleByID(ctx context.Context, id uint) (*model.Article, error) {
	var article model.Article
	if err := config.DB.WithContext(ctx).First(&article, id).Error; err != nil {
		return nil, err
	}
	return &article, nil

}

// DeleteArticle 删除文章（软删）。根据文章 id 删除；硬删可使用 Unscoped。
func DeleteArticle(ctx context.Context, id uint) error {
	if err := config.DB.WithContext(ctx).Delete(&model.Article{}, id).Error; err != nil {
		return err
	}
	return nil
} //查单个,具体搜索
func SearchArticle(ctx context.Context, id uint) (*model.Article, error) {
	var article model.Article
	if err := config.DB.WithContext(ctx).Find(&article, id).Error; err != nil {
		return nil, err
	}
	return &article, nil
}
