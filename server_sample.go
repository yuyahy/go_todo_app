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

// httpサーバーを起動する
func run(ctx context.Context, l net.Listener) error {
	s := &http.Server{
		// 引数で受け取ったnet.Listenerを利用するので、
		// Addrフィールドは指定しない
		// Addr: ":18080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}
	eg, ctx := errgroup.WithContext(ctx)
	// 別goroutineでHTTPサーバーを起動する
	eg.Go(func() error {
		// http.ErrSeverClosedは
		// http.Server.Shutdown()が正常に終了した事を示すので異常ではない
		if err := s.Serve(l); err != nil && err != http.ErrServerClosed {
			log.Printf("failed to close : %v", err)
			return err
		}
		return nil
	})

	// チャネルから終了通知を待機する
	<-ctx.Done()
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown : %v", err)
	}
	return eg.Wait()
}

func main() {
	// 引数でポート番号を指定する
	if len(os.Args) != 2 {
		log.Printf("need port number\n")
		os.Exit(1)
	}
	p := os.Args[1]
	l, err := net.Listen("tcp", ":"+p)
	if err != nil {
		log.Fatalf("failed to listen port %s: %v", p, err)
	}

	// httpサーバーを起動する
	if err := run(context.Background(), l); err != nil {
		log.Printf("failed to terminate server : %v", err)
		os.Exit(1)
	}
}
