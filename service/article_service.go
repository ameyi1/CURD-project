package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"gin-study/config"
	"gin-study/dao"
	"gin-study/model"
	"gin-study/utils"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// GetArticleDetailWithCache 按缓存旁路模式查询文章
func GetArticleDetailWithCache(ctx context.Context, articleID uint) (*model.Article, error) {
	key := fmt.Sprintf("article:detail:%d", articleID)

	cached, err := config.RedisClient.Get(ctx, key).Result()
	if err == nil {
		if cached == "NULL" {
			utils.Log.Warn("文章缓存命中空值", zap.Uint("article_id", articleID))
			return nil, utils.NewAppError(utils.CodeNotFound, "文章不存在")
		}

		var article model.Article
		if json.Unmarshal([]byte(cached), &article) == nil {
			return &article, nil
		}
		utils.Log.Warn("文章缓存反序列化失败，回落数据库", zap.Uint("article_id", articleID), zap.Error(err))
	} else if !errors.Is(err, redis.Nil) {
		utils.Log.Warn("查询文章缓存失败，降级走数据库", zap.Uint("article_id", articleID), zap.Error(err))
	}

	article, err := dao.GetArticleByID(ctx, articleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if setErr := config.RedisClient.Set(ctx, key, "NULL", 30*time.Second).Err(); setErr != nil {
				utils.Log.Warn("写入空值缓存失败", zap.Uint("article_id", articleID), zap.Error(setErr))
			}
			return nil, utils.NewAppError(utils.CodeNotFound, "文章不存在")
		}
		return nil, err
	}

	data, _ := json.Marshal(article)
	expire := 10*time.Minute + time.Duration(rand.Intn(120))*time.Second
	if err := config.RedisClient.Set(ctx, key, data, expire).Err(); err != nil {
		utils.Log.Warn("写入文章缓存失败", zap.Uint("article_id", articleID), zap.Error(err))
	}
	return article, nil
}

