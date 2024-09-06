package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/iinuma0710/react-go-blog/backend/config"
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
		log.Fatalf("failed to listen port %d: %v", cfg.BckendPort, err)
	}

	// サーバの URL を表示
	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("start with: %v", url)

	// ルーティングの設定を取得
	mux := NewMux()

	// Server 型のインスタンスを生成し、HTTP サーバを起動
	s := NewServer(l, mux)
	return s.Run(ctx)
}
