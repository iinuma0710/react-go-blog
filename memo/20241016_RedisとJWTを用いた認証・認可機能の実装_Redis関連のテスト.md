# Redis と JWT を用いた認証・認可機能の実装 ―Redis 関連のテスト―
3連休に出かけていたのと、妻にキングダムを借りて読み始めたら面白くて、しばらくさぼっていました。
今回は、前回定義した Redis 関連の ```KVS``` 型のテストを実装していきいます。

## ```testutil``` への実装
```testutil``` パッケージの下に ```kvs.go``` を実装します。
```testutil/db.go``` の ```OpenDBForTest``` と同様に、テストの実行環境に合わせて接続情報の差分を吸収する ```OpenRedisForTest``` 関数を実装します。


```backend/testuril/kvs.go```
```go
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

	host := "127.0.0.1"
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
```

## ```KVS``` 型のメソッドのテストコード
まずは、```Save``` メソッドに対するテストを実装します。
ここでは、キーとして関数名と同じ ```TestKVS_Save``` を利用しています。
Redis の全てのテストケースでは、キーが衝突しないように設定する必要がありますが、テストケースと同じ命名規則でキーを作成しておけば、その心配はありません。

```backend/store/kvs_test.go```
```go
package store

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/iinuma0710/react-go-blog/backend/entity"
	"github.com/iinuma0710/react-go-blog/backend/testutil"
)

func TestKVS_Save(t *testing.T) {
	t.Parallel()

	cli := testutil.OpenRedisForTeest(t)
	sut := &KVS{Cli: cli}
	key := "TestKVS_Save"
	uid := entity.UserID(1234)
	ctx := context.Background()
	t.Cleanup(func() {
		cli.Del(ctx, key)
	})
	if err := sut.Save(ctx, key, uid); err != nil {
		t.Errorf("want no error, but got %v", err)
	}
}
```

次に、```Load``` メソッドのテストを実装します。
登録済みのキーで値が取得できることを確認するサブテストと、登録されていないキーでは値を取得できないことを確認するサブテストを、それぞれ ```t.Run``` メソッドで実行します。
ただし、Redis クライアントをサブテストごとに実装する必要はないので、各サブテストで共有しています。

```backend/store/kvs_test.go```
```go
func TestKVS_Load(t *testing.T) {
	t.Parallel()

	cli := testutil.OpenRedisForTeest(t)
	sut := &KVS{Cli: cli}

	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		key := "TestKVS_Load_ok"
		uid := entity.UserID(1234)
		ctx := context.Background()
		cli.Set(ctx, key, int64(uid), 30*time.Minute)
		t.Cleanup(func() {
			cli.Del(ctx, key)
		})
		got, err := sut.Load(ctx, key)
		if err != nil {
			t.Fatalf("want no error, but got %v", err)
		}
		if got != uid {
			t.Errorf("want %d, but got %d", uid, got)
		}
	})

	t.Run("notFound", func(t *testing.T) {
		t.Parallel()

		key := "TestKVS_Load_notFound"
		ctx := context.Background()
		got, err := sut.Load(ctx, key)
		if err == nil || !errors.Is(err, ErrNotFound) {
			t.Errorf("want %v, but got %v (value = %v)", ErrNotFound, err, got)
		}
	})
}
```

最後に ```go test -v ./...``` コマンドでテストが通ることを確認しておきます。