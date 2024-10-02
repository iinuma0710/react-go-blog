package service

import (
	"context"

	"github.com/iinuma0710/react-go-blog/backend/entity"
	"github.com/iinuma0710/react-go-blog/backend/store"
)

//go:generate go run github.com/matryer/moq -out moq_test.go . ArticleAdder ArticleLister
type ArticleAdder interface {
	AddArticle(ctx context.Context, db store.Execer, a *entity.Article) error
}

type ArticleLister interface {
	ListArticles(ctx context.Context, db store.Queryer) (entity.Articles, error)
}
