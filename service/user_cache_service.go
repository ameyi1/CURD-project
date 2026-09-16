package service

/*
//现在还不会但看着高级就先放着吧,这个叫 Cache Aside 模式
import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"gin-study/config"
	"gin-study/dao"
	"gin-study/model"
)

func GetUserWithCache(ctx context.Context, id uint) (*model.User, error) {
	// 1. 先从 Redis 找
	key := "user_cache:" + strconv.FormatUint(uint64(id), 10)
	cached, err := config.RedisClient.Get(ctx, key).Result()

	if err == nil {
		// 找到了！反序列化
		var user model.User
		if json.Unmarshal([]byte(cached), &user) == nil {
			return &user, nil
		}
	}

	// 2. 没找到，查数据库
	user, err := dao.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 3. 存进 Redis（设置过期时间 10 分钟）
	data, _ := json.Marshal(user)
	config.RedisClient.Set(ctx, key, data, 10*time.Minute)

	return user, nil
}
*/
