package handler

import (
	"net/http"

	"github.com/iinuma0710/react-go-blog/backend/entity"
)

type ListArticle struct {
	Service ListArticlesService
}

type article struct {
	ID     entity.ArticleID     `json:"id"`
	Title  string               `json:"title"`
	Status entity.ArticleStatus `json:"status"`
}

func (la *ListArticle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	articles, err := la.Service.ListArticles(ctx)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	rsp := []article{}
	for _, a := range articles {
		rsp = append(rsp, article{
			ID:     a.ID,
			Title:  a.Title,
			Status: a.Status,
		})
	}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
