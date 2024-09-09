package testutil

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
)

// テスト用にDBに接続するヘルパー関数
func OpenDBForTest(t *testing.T) *sqlx.DB {
	t.Helper()

	port := 33306
	// 環境変数:"CI"はGitHub Actions上でのみ定義されている想定
	if _, defined := os.LookupEnv("CI"); defined {
		port = 3306
	}
	db, err := sql.Open("mysql", fmt.Sprintf("todo:todo@tcp(127.0.0.1:%d)/todo?parseTime=true", port))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return sqlx.NewDb(db, "mysql")
}
