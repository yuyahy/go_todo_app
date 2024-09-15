package service

import (
	"context"
	"fmt"

	"github.com/yuyahy/go_todo_app/store"
)

type Login struct {
	DB             store.Queryer
	Repo           UserGetter
	TokenGenerator TokenGenerator
}

func (l *Login) Login(ctx context.Context, name, pw string) (string, error) {
	u, err := l.Repo.GetUser(ctx, l.DB, name)
	// 入力パスワードとハッシュ化して永続化していた文字列を比較し、認証を行う
	if err != nil {
		return "", fmt.Errorf("failed to list: %w", err)
	}
	if err = u.ComparePassword(pw); err != nil {
		return "", fmt.Errorf("wrong password: %w", err)
	}
	// 認証がOKであれば、JWTを発行
	jwt, err := l.TokenGenerator.GenerateToken(ctx, *u)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	return string(jwt), nil
}
