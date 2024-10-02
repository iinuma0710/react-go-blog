package handler

import (
	"context"

	"github.com/iinuma0710/react-go-blog/backend/entity"
)

//go:generate go run github.com/matryer/moq -out moq_test.go . ListArticlesService AddArticleService
type ListArticlesService interface {
	ListArticles(ctx context.Context) (entity.Articles, error)
}

type AddArticleService interface {
	AddArticle(ctx context.Context, title string) (*entity.Article, error)
}
