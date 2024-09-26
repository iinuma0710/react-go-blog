package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/iinuma0710/react-go-blog/backend/entity"
	"github.com/iinuma0710/react-go-blog/backend/store"
	"github.com/jmoiron/sqlx"
)

type AddArticle struct {
	DB        *sqlx.DB
	Repo      *store.Repository
	Validator *validator.Validate
}

func (aa *AddArticle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// リクエストボディに必要な項目を構造体として定義
	var b struct {
		Title  string               `json:"title" validate:"required"`
		Status entity.ArticleStatus `json:"status" validate:"required,oneof=draft published withdrawn"`
	}

	// リクエストのボディをでコード
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	// リクエストをでコードした JSON に必要な項目が含まれていることを確認
	if err := aa.Validator.Struct(b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	// 新しい Article 型の値を作成
	t := &entity.Article{
		Title:  b.Title,
		Status: b.Status,
	}

	// データベースに記事を追加
	err := aa.Repo.AddArticle(ctx, aa.DB, t)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	// 登録された記事の ID とステータスコードを含めてレスポンスを返す
	rsp := struct {
		ID entity.ArticleID `json:"id"`
	}{ID: t.ID}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
