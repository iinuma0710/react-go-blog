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
	ID        ArticleID     `json:"id" db:"id"`
	Title     string        `json:"title" db:"title"`
	Status    ArticleStatus `json:"status" db:"status"`
	CreatedAt time.Time     `json:"crated_at" db:"created_at"`
}

type Articles []*Article
