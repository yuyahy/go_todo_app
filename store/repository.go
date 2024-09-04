package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/yuyahy/go_todo_app/config"
)

func New(ctx context.Context, cfg *config.Config) (*sqlx.DB, func(), error) {
	// sqlx.Connextを使うと内部でpingする
	db, err := sql.Open("mysql",
		fmt.Sprintf(
			// parseTime=truがないと、time.Timeのフィールドに正しい値が設定できない
			"%s:%s@tcp(%s:%d)/%s?parseTime=true",
			cfg.DBUser, cfg.DBPassword,
			cfg.DBHost, cfg.DBPort,
			cfg.DBName,
		),
	)
	if err != nil {
		return nil, nil, err
	}
	// Openは接続テストは行われない
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, func() { _ = db.Close() }, err
	}
	xdb := sqlx.NewDb(db, "mysql")
	return xdb, func() { _ = db.Close() }, nil
}
