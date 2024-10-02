package store

import (
	"context"

	"github.com/iinuma0710/react-go-blog/backend/entity"
)

func (r *Repository) ListArticles(ctx context.Context, db Queryer) (entity.Articles, error) {
	articles := entity.Articles{}
	sql := `SELECT id, title, status, created_at FROM article;`

	if err := db.SelectContext(ctx, &articles, sql); err != nil {
		return nil, err
	}

	return articles, nil
}

func (r *Repository) AddArticle(ctx context.Context, db Execer, a *entity.Article) error {
	a.CreatedAt = r.Clocker.Now()
	sql := `INSERT INTO article
		(title, status, created_at)
		VALUES (?, ?, ?)`

	result, err := db.ExecContext(ctx, sql, a.Title, a.Status, a.CreatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	a.ID = entity.ArticleID(id)
	return nil
}
