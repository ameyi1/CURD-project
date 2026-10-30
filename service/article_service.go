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

// 增
func CreateArticle(ctx context.Context, article model.Article, userID uint) error {
	if err := dao.CreateArticle(&article); err != nil {
		return fmt.Errorf("创建文章失败(userID=%d,articleID=%d,title=%s):%w", userID, article.ID, article.Title, err)
	}
	utils.Log.Info("增加文章成功", zap.Uint("userID", userID), zap.Uint("articleID", article.ID), zap.String("title", article.Title))
	return nil
}

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

// 软删
func DeleteArticle(ctx context.Context, articleID, id uint) error {
	//要校验下是不是本人
	//先调用下getarticlebyid
	article, err := dao.GetArticleByID(ctx, articleID)
	if err != nil {
		//先判断下文章存不存在
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return utils.NewAppError(utils.CodeNotFound, "文章不存在")

		}

		return utils.NewAppError(utils.CodeServerError, "调用dao_getarticlebyid失败")

	}
	//判断下
	if id != article.UserID {
		return utils.NewAppError(utils.CodeForbidden, "用户id不匹配")

	}
	//再调用删除
	if err := dao.DeleteArticle(ctx, articleID); err != nil {
		return utils.NewAppError(utils.CodeServerError, "调用dao_deletearticle失败")

	}
	return nil
} //查单个
func SearchArticle(ctx context.Context, articleID uint) (*model.Article, error) {
	article, err := dao.SearchArticle(ctx, articleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil,
				utils.NewAppError(utils.CodeNotFound, "文章不存在")

		}
		return nil,
			utils.NewAppError(utils.CodeServerError, "查询文章失败")

	}
	return article, nil
}

// 查分页
func ListArticles(ctx context.Context, keyword string, page, size int) ([]model.Article, int64, error) {
	offset := (page - 1) * size
	article, total, err := dao.ListArticles(ctx, keyword, offset, size)
	if err != nil {
		return nil, 0, fmt.Errorf("分页查文章失败(keyword=%s,page=%d,size=%d):%w", keyword, page, size, err)
	}
	return article, total, nil
}

// 改，先校验归属，再改，最后删缓存
func UpdateArticle(ctx context.Context, articleID, userID uint, title, context string) error {
	article, err := dao.GetArticleByID(ctx, articleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.NewAppError(utils.CodeNotFound, "文章不存在")
		}
		return fmt.Errorf("查询文章失败articleID=%d:%w", articleID, err)
	}
	if article.UserID != userID {
		return utils.NewAppError(utils.CodeForbidden, "无权操作别人的文章")
	}
	if err := dao.UpdateArticle(ctx, title, context, articleID); err != nil {
		return fmt.Errorf("更新文章记录失败 articleID=%d: %w", articleID, err)
	}
	//删缓存
	if err := config.RedisClient.Del(ctx, keyArticleDetail(articleID)).Err(); err != nil {
		utils.Log.Warn("删除文章缓存失败", zap.Uint("article_id", articleID), zap.Error(err))
	}
	return nil
}

func keyArticleDetail(articleID uint) string {
	panic("unimplemented")
}

// 点赞,有原子性问题没解决2026-10-14
func LikeArticle(ctx context.Context, ArticleID uint, UserID uint) error {
	//先设置key,设置两个key，一个记录点赞用户id，一个记录点赞次数
	userSetkey := fmt.Sprintf("like:user_id:%d", ArticleID)
	countkey := fmt.Sprintf("like:count_id:%d", ArticleID)
	//然后set通过redis，也就是存内存里
	added, err := config.RedisClient.SAdd(ctx, userSetkey, UserID).Result()
	//先判断added的情况，如果为零就说明重复点赞，这是redis里面set（无序，不重复）的功能去重，重复返回0
	if added == 0 {
		return utils.NewAppError(utils.CodeForbidden, "不能重复点赞")
	}
	if err != nil {
		return utils.NewAppError(utils.CodeServerError, "出错了请重试")

	}
	if err := config.RedisClient.Incr(ctx, countkey).Err(); err != nil {
		return err
	}
	return nil
} //取消点赞，和点赞几乎一样，无非是srem，decr
func UnlikeArticle(ctx context.Context, ArticleID uint, UserID uint) error {
	//同样先设置key
	usersetkey := fmt.Sprintf("like:user_id:%d", ArticleID)
	countkey := fmt.Sprintf("like:count:%d", ArticleID)
	removed, err := config.RedisClient.SRem(ctx, usersetkey, UserID).Result()
	if removed == 0 {
		return utils.NewAppError(utils.CodeInvalidParam, "不能重复取消点赞")
	}
	if err != nil {
		return utils.NewAppError(utils.CodeServerError, "出错了请重新再试")

	}
	if err := config.RedisClient.Decr(ctx, countkey).Err(); err != nil {
		return err
	}
	return nil
}

// 查询用户是否点赞
func IsLike(ctx context.Context, UserID uint, ArticleID uint) error {
	//同样的先设置key
	usersetKey := fmt.Sprintf("like:user_id:%d", ArticleID)
	err := config.RedisClient.SIsMember(ctx, usersetKey, UserID).Err()
	if err != nil {
		return utils.NewAppError(utils.CodeServerError, "出错了请重试")
	}
	return nil
}

// 查询点赞个数
func GetLikeCount(ctx context.Context, ArticleID uint) (int64, error) {
	countkey := fmt.Sprintf("like:count:%d", ArticleID)
	count, err := config.RedisClient.Get(ctx, countkey).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, utils.NewAppError(utils.CodeServerError, "查询失败")
	}
	return count, nil
}

// 写个文章浏览量吧,相当于简化的点赞系统，不用去重和获取用户id
func Addview(ctx context.Context, articleID uint) {
	key := fmt.Sprintf("view:count:%d", articleID)
	err := config.RedisClient.Incr(ctx, key).Err()
	if err != nil {
		utils.Log.Error("浏览量增加失败", zap.Error(err))
	}

} //查询浏览量
func GetViewCount(ctx context.Context, articleID uint) (int64, error) {
	countkey := fmt.Sprintf("view:count:%d", articleID)
	Count, err := config.RedisClient.Get(ctx, countkey).Int64()
	if err == redis.Nil { //key不存在，说明没人看过，正常现象，不是错误
		return 0, nil
	}
	if err != nil {
		return 0, utils.NewAppError(utils.CodeServerError, "查询浏览量失败")
	}
	return Count, nil
}
