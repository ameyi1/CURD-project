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

// 查 分页
func ListArticles(ctx context.Context, keyword string, offset, size int) ([]model.Article, int64, error) {
	var list []model.Article
	var total int64
	base := config.DB.WithContext(ctx).Model(&model.Article{})
	//模糊查询
	if keyword != "" {
		listkeyword := "%" + keyword + "%"
		base = base.Where("title LIKE ?", listkeyword)
	}
	//算总页数
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err

	}
	//分页查
	if err := base.Limit(size).Offset(offset).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// 改
func UpdateArticle(ctx context.Context, title, context string, articleID uint) error {
	return config.DB.WithContext(ctx).
		Model(&model.Article{}).
		Where("articleID=?", articleID).
		Updates(map[string]interface{}{
			"title":   title,
			"context": context,
		}).Error
}
