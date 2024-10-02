package store

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/iinuma0710/react-go-blog/backend/clock"
	"github.com/iinuma0710/react-go-blog/backend/entity"
	"github.com/iinuma0710/react-go-blog/backend/testutil"
	"github.com/jmoiron/sqlx"
)

func prepareArticles(ctx context.Context, t *testing.T, con Execer) entity.Articles {
	t.Helper()

	// 一旦データベースの中身をきれいにしておく
	if _, err := con.ExecContext(ctx, "DELETE FROM article;"); err != nil {
		t.Logf("failed to initialize article: %v", err)
	}

	c := clock.FixedClocker{}
	wants := entity.Articles{
		{
			Title:     "wants article 1",
			Status:    "published",
			CreatedAt: c.Now(),
		},
		{
			Title:     "wants article 1",
			Status:    "draft",
			CreatedAt: c.Now(),
		},
		{
			Title:     "wants article 3",
			Status:    "withdrawn",
			CreatedAt: c.Now(),
		},
	}

	result, err := con.ExecContext(ctx, `
		INSERT INTO article (title, status, created_at)
		VALUES
			(?, ?, ?),
			(?, ?, ?),
			(?, ?, ?);`,
		wants[0].Title, wants[0].Status, wants[0].CreatedAt,
		wants[1].Title, wants[1].Status, wants[1].CreatedAt,
		wants[2].Title, wants[2].Status, wants[2].CreatedAt,
	)
	if err != nil {
		t.Fatal(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	wants[0].ID = entity.ArticleID(id)
	wants[1].ID = entity.ArticleID(id + 1)
	wants[2].ID = entity.ArticleID(id + 2)

	return wants
}

func TestRepository_ListTasks(t *testing.T) {
	ctx := context.Background()

	// トランザクションをはることでこのテストケースの中だけのテーブル状態にする。
	tx, err := testutil.OpenDBForTest(t).BeginTxx(ctx, nil)
	// テストが完了したら元に戻す
	t.Cleanup(func() { _ = tx.Rollback() })
	if err != nil {
		t.Fatal(err)
	}
	wants := prepareArticles(ctx, t, tx)

	sut := &Repository{}
	gots, err := sut.ListArticles(ctx, tx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d := cmp.Diff(gots, wants); len(d) != 0 {
		t.Errorf("differs: (-got +want)\n%s", d)
	}
}

func TestRepository_AddArticle(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c := clock.FixedClocker{}
	var wantID int64 = 20
	okTask := &entity.Article{
		Title:     "ok article",
		Status:    "published",
		CreatedAt: c.Now(),
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })

	mock.ExpectExec(
		// エスケープが必要
		`INSERT INTO article \(title, status, created_at\) VALUES \(\?, \?, \?\)`,
	).WithArgs(okTask.Title, okTask.Status, c.Now()).
		WillReturnResult((sqlmock.NewResult(wantID, 1)))

	xdb := sqlx.NewDb(db, "mysql")
	r := &Repository{Clocker: c}
	if err := r.AddArticle(ctx, xdb, okTask); err != nil {
		t.Errorf("want no error, but got %v", err)
	}
}
