package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/yuyahy/go_todo_app/config"
)

// httpサーバーを起動する
func run(ctx context.Context) error {
	cfg, err := config.New()
	if err != nil {
		return err
	}
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen port %d: %v", cfg.Port, err)
	}
	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("start with: %v", url)
	mux, cleanup, err := NewMux(ctx, cfg)
	if err != nil {
		return err
	}
	defer cleanup()
	s := NewServer(l, mux)
	return s.Run(ctx)
}

func main() {
	// httpサーバーを起動する
	if err := run(context.Background()); err != nil {
		log.Printf("failed to terminate server : %v", err)
		os.Exit(1)
	}
}
