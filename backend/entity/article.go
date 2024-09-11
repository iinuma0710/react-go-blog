package entity

import "time"

type ArticleID int64
type ArticleStatus string

const (
	ArticleDraft     ArticleStatus = "draft"
	ArticlePublished ArticleStatus = "published"
	ArticleWithdrawn ArticleStatus = "withdrawn"
)

type Article struct {
	ID        ArticleID     `json:"id"`
	Title     string        `json:"title"`
	Status    ArticleStatus `json:"status"`
	CreatedAt time.Time     `json:"crated_at"`
}

type Articles []*Article
