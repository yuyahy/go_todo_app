package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/sync/errgroup"
)

// httpサーバーを起動する
func run(ctx context.Context) error {
	s := &http.Server{
		Addr: ":18080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}
	eg, ctx := errgroup.WithContext(ctx)
	// 別goroutineでHTTPサーバーを起動する
	eg.Go(func() error {
		// http.ErrSeverClosedは
		// http.Server.Shutdown()が正常に終了した事を示すので異常ではない
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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
	// httpサーバーを起動する
	if err := run(context.Background()); err != nil {
		log.Printf("failed to terminate server : %v", err)
	}
}
