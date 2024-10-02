package service

import (
	"context"
	"fmt"

	"github.com/iinuma0710/react-go-blog/backend/entity"
	"github.com/iinuma0710/react-go-blog/backend/store"
)

type AddArticle struct {
	DB   store.Execer
	Repo ArticleAdder
}

func (aa *AddArticle) AddArticle(ctx context.Context, title string) (*entity.Article, error) {
	a := &entity.Article{
		Title:  title,
		Status: entity.ArticleDraft,
	}

	err := aa.Repo.AddArticle(ctx, aa.DB, a)
	if err != nil {
		return nil, fmt.Errorf("failed to resister: %w", err)
	}

	return a, nil
}
