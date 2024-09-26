package handler

import (
	"net/http"
	"time"

	"github.com/iinuma0710/react-go-blog/backend/entity"
	"github.com/iinuma0710/react-go-blog/backend/store"
	"github.com/jmoiron/sqlx"
)

type ListArticle struct {
	DB   *sqlx.DB
	Repo *store.Repository
}

type article struct {
	ID        entity.ArticleID     `json:"id"`
	Title     string               `json:"title"`
	Status    entity.ArticleStatus `json:"status"`
	CreatedAt time.Time            `json:"created_at"`
}

func (la *ListArticle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	articles, err := la.Repo.ListArticles(ctx, la.DB)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	rsp := []article{}
	for _, a := range articles {
		rsp = append(rsp, article{
			ID:        a.ID,
			Title:     a.Title,
			Status:    a.Status,
			CreatedAt: a.CreatedAt,
		})
	}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
