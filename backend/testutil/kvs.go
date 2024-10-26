package testutil

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/redis/go-redis/v9"
)

func OpenRedisForTeest(t *testing.T) *redis.Client {
	t.Helper()

	host := "blog_redis"
	port := 6379
	if _, defined := os.LookupEnv("CI"); defined {
		// GitHub Actions で用いるポート番号 (未実装)
		port = 6379
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: "",
		DB:       0, // デフォルトのデータベース番号
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		t.Fatalf("failed to connect redis: %s", err)
	}

	return client
}
