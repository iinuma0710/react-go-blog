package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/iinuma0710/react-go-blog/backend/clock"
	"github.com/iinuma0710/react-go-blog/backend/config"
	"github.com/jmoiron/sqlx"
)

func New(ctx context.Context, cfg *config.Config, maxTrial int) (*sqlx.DB, func(), error) {
	// 接続先のデータベースのパス
	path := fmt.Sprintf(
		// parseTime=true は時刻情報の取得に必須
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)
	
	var db *sql.DB
	var err error
	for i := 0; i < maxTrial; i++ {
		fmt.Printf("mysql connection trial: %d", i + 1)

		// database/sql の Open メソッドで接続
		db, err = sql.Open("mysql", path)
		if err != nil {
			fmt.Printf("sql.Open method failed: %v", err)
			time.Sleep(time.Second * 2)
			continue
		}

		// *sql.DB.PingContest メソッドで疎通確認
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
		if err := db.PingContext(ctx); err != nil {
			fmt.Printf("*sql.DB.PingContext method failed: %v", err)
			time.Sleep(time.Second * 2)
			continue
		}

		// ここまでエラーがなければ接続確認できているので接続試行のループを抜ける
		break
	}

	// 何らかのエラーで接続できなかった場合の処理
	if err != nil {
		if db != nil {
			return nil, func() { _ = db.Close() }, fmt.Errorf("Cannot open sql connection: %v", err)
		} else {
			return nil, func() {}, fmt.Errorf("Cannot confirm sql connection: %v", err)
		}
	}
	
	// *sqlx.DB に変換して返す
	xdb := sqlx.NewDb(db, "mysql")
	return xdb, func() { _ = db.Close() }, nil
}

type Beginner interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type Preparer interface {
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
}

type Execer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}

type Queryer interface {
	Preparer
	QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error)
	QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row
	GetContext(ctx context.Context, dest interface{}, query string, args ...any) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...any) error
}

// インターフェースが期待通りに宣言されているかを確認
var (
	_ Beginner = (*sqlx.DB)(nil)
	_ Preparer = (*sqlx.DB)(nil)
	_ Queryer  = (*sqlx.DB)(nil)
	_ Execer   = (*sqlx.DB)(nil)
	_ Execer   = (*sqlx.Tx)(nil)
)

type Repository struct {
	Clocker clock.Clocker
}
