package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/iinuma0710/react-go-blog/backend/clock"
	"github.com/iinuma0710/react-go-blog/backend/config"
	"github.com/iinuma0710/react-go-blog/backend/handler"
	"github.com/iinuma0710/react-go-blog/backend/service"
	"github.com/iinuma0710/react-go-blog/backend/store"
)

func NewMux(ctx context.Context, cfg *config.Config) (http.Handler, func(), error) {
	mux := chi.NewRouter()

	// ヘルスチェック用のエンドポイント
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		// 静的解析ツールのエラー回避のため、明示的に戻り値を捨てる
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})

	v := validator.New()

	// データベースに接続
	db, cleanup, err := store.New(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}

	// store.Repository 型のインスタンスを生成
	r := store.Repository{Clocker: clock.RealClocker{}}

	// 記事を追加するためのエンドポイント
	aa := &handler.AddArticle{
		Service:   &service.AddArticle{DB: db, Repo: &r},
		Validator: v,
	}
	mux.Post("/articles", aa.ServeHTTP)

	// 記事一覧を取得するためのエンドポイント
	la := &handler.ListArticle{
		Service: &service.ListArticle{DB: db, Repo: &r},
	}
	mux.Get("/articles", la.ServeHTTP)

	return mux, cleanup, nil
}
