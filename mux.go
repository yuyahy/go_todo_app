package main

import "net/http"

// どのようなハンドラー実装をどんなパスURLで公開するかをルーティングする
// ※戻り値を*http.ServeMuxではなく、http.Handlerインターフェースにする事で内部実装に依存しない関数シグネチャになる
func NewMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		// 静的解析のエラーを回避するために明示的に戻り値を捨てている
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	return mux
}
