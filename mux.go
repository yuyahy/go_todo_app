package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/yuyahy/go_todo_app/handler"
	"github.com/yuyahy/go_todo_app/store"
)

// どのようなハンドラー実装をどんなパスURLで公開するかをルーティングする
// ※戻り値を*http.ServeMuxではなく、http.Handlerインターフェースにする事で内部実装に依存しない関数シグネチャになる
func NewMux() http.Handler {
	mux := chi.NewRouter()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	v := validator.New()
	// AddTaskとListTaskで永続化情報を共有するためにstore.Tasksを渡す
	at := &handler.AddTask{Store: store.Tasks, Validator: v}
	mux.Post("/tasks", at.ServeHTTP)
	lt := &handler.ListTask{Store: store.Tasks}
	mux.Get("/tasks", lt.ServeHTTP)
	return mux
}
