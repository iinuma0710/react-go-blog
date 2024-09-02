package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/iinuma0710/react-go-blog/backend/config"
	"golang.org/x/sync/errgroup"
)

func main() {
	// run 関数を呼び出す
	if err := run(context.Background()); err != nil {
		fmt.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	// 環境変数で指定された設定値を取得
	cfg, err := config.New()
	if err != nil {
		return err
	}

	// 環境変数 BACKEND_PORT で設定されたポートのリッスンを開始
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.BckendPort))
	if err != nil {
		log.Fatalf("failed to listen port %s: %v", cfg.BckendPort, err)
	}

	// サーバの URL を表示
	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("start with: %v", url)

	s := &http.Server{
		// Addr フィールドは指定しない
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s!\n", r.URL.Path[1:])
		}),
	}

	eg, ctx := errgroup.WithContext(ctx)

	// eg.Go メソッドで HTTP サーバを起動するゴルーチンを立ち上げる
	eg.Go(func() error {
		// Server メソッドで HTTP サーバを立ち上げ
		if err := s.Serve(l); err != nil && err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})

	// チャネルからの終了通知を待機
	<-ctx.Done()
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}

	// Go メソッドで起動したゴルーチンの終了を待つ
	return eg.Wait()
}
