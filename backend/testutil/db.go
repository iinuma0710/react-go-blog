package testutil

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func OpenDBForTest(t *testing.T) *sqlx.DB {
	t.Helper()

	// ローカル環境でテストする場合のポート番号
	port := 3306
	if _, defined := os.LookupEnv("CI"); defined {
		// GitHub Actions で用いるポート番号 (未実装)
		port = 3306
	}

	db, err := sql.Open(
		"mysql",
		fmt.Sprintf("blog:blog@tcp(blog_database:%d)/blog?parseTime=true", port),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(
		func() { _ = db.Close() },
	)

	return sqlx.NewDb(db, "mysql")
}
