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
