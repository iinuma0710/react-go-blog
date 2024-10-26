package store

import (
	"context"
	"fmt"
	"time"

	"github.com/iinuma0710/react-go-blog/backend/config"
	"github.com/iinuma0710/react-go-blog/backend/entity"
	"github.com/redis/go-redis/v9"
)

type KVS struct {
	Cli *redis.Client
}

func NewKVS(ctx context.Context, cfg *config.Config, maxTrial int) (*KVS, error) {
	cli := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
	})

	// 2秒ごとに maxTrial 回だけ Ping を送って接続を試みる
	var err error
	for i := 0; i < maxTrial; i++ {
		fmt.Printf("redis connection trial: %d", i+1)
		if err = cli.Ping(ctx).Err(); err != nil {
			fmt.Printf("*redis.Client.Ping method failed: %v", err)
			time.Sleep(time.Second * 2)
			continue
		} else {
			// Ping が返ってきたら接続成功としてループを抜ける
			break
		}
	}

	if err != nil {
		return nil, err
	} else {
		return &KVS{Cli: cli}, nil
	}
}

func (k *KVS) Save(ctx context.Context, key string, userID entity.UserID) error {
	id := int64(userID)
	return k.Cli.Set(ctx, key, id, 30*time.Minute).Err()
}

func (k *KVS) Load(ctx context.Context, key string) (entity.UserID, error) {
	id, err := k.Cli.Get(ctx, key).Int64()
	if err != nil {
		return 0, fmt.Errorf("failed to get by %q: %w", key, ErrNotFound)
	}
	return entity.UserID(id), nil
}
