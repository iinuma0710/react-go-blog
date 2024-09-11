package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/iinuma0710/react-go-blog/backend/handler"
	"github.com/iinuma0710/react-go-blog/backend/store"
)

func NewMux() http.Handler {
	mux := chi.NewRouter()

	// ヘルスチェック用のエンドポイント
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		// 静的解析ツールのエラー回避のため、明示的に戻り値を捨てる
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})

	// 記事を追加するためのエンドポイント
	v := validator.New()
	aa := &handler.AddArticle{Store: store.Articles, Validator: v}
	mux.Post("/articles", aa.ServeHTTP)

	// 記事一覧を取得するためのエンドポイント
	la := &handler.ListArticle{Store: store.Articles}
	mux.Get("/articles", la.ServeHTTP)

	return mux
}
