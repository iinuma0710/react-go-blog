package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"golang.org/x/sync/errgroup"
)

func main() {
	// コマンドライン引数でポート番号を指定
	if len(os.Args) != 2 {
		log.Printf("need port number\n")
		os.Exit(1)
	}

	// 指定されたポートのリッスンを開始
	p := os.Args[1]
	l, err := net.Listen("tcp", ":"+p)
	if err != nil {
		log.Fatalf("failed to listen port %s: %v", p, err)
	}

	// run 関数にリッスンしているポートの情報を渡して HTTP サーバを立ち上げる
	if err := run(context.Background(), l); err != nil {
		fmt.Printf("failed to terminate server: %v", err)
	}
}

func run(ctx context.Context, l net.Listener) error {
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
