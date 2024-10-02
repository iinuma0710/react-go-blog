package service

import (
	"context"
	"fmt"

	"github.com/iinuma0710/react-go-blog/backend/entity"
	"github.com/iinuma0710/react-go-blog/backend/store"
)

type ListArticle struct {
	DB   store.Queryer
	Repo ArticleLister
}

func (l *ListArticle) ListArticles(ctx context.Context) (entity.Articles, error) {
	as, err := l.Repo.ListArticles(ctx, l.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to list: %w", err)
	}
	return as, nil
}
