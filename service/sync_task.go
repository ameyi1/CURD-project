package service

import (
	"context"
	"fmt"
	"time"

	"gin-study/config"
	"gin-study/model"
)

// StartSyncTask 启动定时同步任务
func StartSyncTask() {
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for range ticker.C {
			syncLikeCount()
		}
	}()
}

// syncLikeCount 同步点赞数到 MySQ
func syncLikeCount() {
	ctx := context.Background()

	// 1. 找出所有 like:count:* 的 key
	keys, err := config.RedisClient.Keys(ctx, "like:count:*").Result()
	if err != nil {
		return
	}
	//这个肯定不用，n+1问题来不就炸了~·
	for _, key := range keys {
		// 2. 从 key 中解析文章ID
		var articleID uint
		fmt.Sscanf(key, "like:count:%d", &articleID)

		// 3. 获取当前计数
		count, err := config.RedisClient.Get(ctx, key).Int64()
		if err != nil {
			continue
		}

		// 4. 更新 MySQL
		config.DB.Model(&model.Article{}).
			Where("id = ?", articleID).
			Update("like_count", count)
	}
} //回头再看
/*time.NewTicker(10 * time.Second)：创建一个定时器，每 10 秒触发一次。
go func() { for range ticker.C { ... } }()：在 goroutine 里循环，每次触发执行同步。
Keys(ctx, "like:count:*")：找出所有匹配的 key。
fmt.Sscanf(key, "like:count:%d", &articleID)：从 key 里解析文章 ID。
Get(ctx, key).Int64()：获取计数。
Update("like_count", count)：更新 MySQL。*/
