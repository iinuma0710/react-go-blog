# Redis と JWT を用いた認証・認可機能の実装 ―Redisの準備―
別のブランチでユーザの登録機能を実装したので、その情報を使ってログインしたり、アクセス管理を行ったりする機能を実装します。
それにあたって、まずは Redis を使う準備を行います。

## Redis とは
[Redis](https://redis.io/) とは **REmote DIctionary Server** の略で、オープンソースの Key-Value 型 NoSQL データベースです。
Redis の特徴としては、

- インメモリデータベース  
  データをメモリ上に保存するため、二次記憶にデータを保存する MySQL 等の RDBMS に比べて高速に動作します。
- 多様なデータ構造  
  基本的には Key-Value 型ですが、リストやセット、ハッシュなどのデータ構造にも対応しています。
- 永続化機能  
  原則的にインメモリで動作しますが、データを二次記憶上に永続化することも可能です。
- 分散性  
  複数台の仮想サーバやコンテナにスケールアウト可能で、クラスタ化やレプリケーションで高可用性と対象外性を実現しています。
- Pub/Sub 機能  
  通知機能やリアルタイム更新機能の実装に利用することもできます。

今回実装する認証・認可機能ではアクセストークンを発行しますが、有効期限が切れたトークンは無効化する必要があるので、永続化機能は利用しません。
また、クラウドネイティブな運用も考慮して、仮想サーバやコンテナがステートレスとなるよう、ミドルウェア上で保管するように実装を進めていきます。

## ```docker-compose.yml``` の準備
Docker コンテナとして Redis を立ち上げるため、```blog_redis``` サービスを追記します。

```docker-compose.yml```
```yml
services:
  # ...略...

  blog_redis:
    image: redis:7.4.1
    container_name: blog_redis_container
    ports:
      - "6379:6379"
    volumes:
      - blog_redis_data:/data

  # ...略...
```

Go の HTTP サーバを動かしているコンテナからアクセスできるように、```blog_backend``` に環境変数を追加しておきます。

```docker-compose.yml```
```yml
services:
  blog_backend:
    # ...略...
    
    environment:
      # ...略...
      - BLOG_REDIS_HOST=blog_redis
      - BLOG_REDIS_PORT=6379
   
    # ...略...
```

また、データを保持するボリュームも用意しておきます。

```docker-compose.yml```
```yml
# ...略...

volumes:
  blog_database_data:
  blog_redis_data:
```

## アプリケーションコードの準備
Go のコンテナから Redis のデータを操作するコードを用意します。
これに先立って、Redis クライアントの ```github.com/redis/go-redis/v9``` を ```go get``` コマンドで取得しておきます。  
```store``` パッケージ内に、Redis への接続を行う関数と、アクセストークンを保存するメソッドを定義します。

```backend/store/kvs.go```
```go
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
```

```config.go``` の ```Config``` 構造体に、Redis 関連の設定項目を追加しておきます。

```backend/config/config.go```
```go
type Config struct {
	...
	RedisHost  string `env:"BLOG_REDIS_HOST" envDefault:"blog_redis"`
	RedisPort  int    `env:"BLOG_REDIS_PORT" envDefault:"6379"`
}
```