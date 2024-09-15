package auth

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/yuyahy/go_todo_app/clock"
	"github.com/yuyahy/go_todo_app/entity"
	"github.com/yuyahy/go_todo_app/store"
	"github.com/yuyahy/go_todo_app/testutil/fixture"
)

func TestEmbed(t *testing.T) {
	want := []byte("-----BEGIN PUBLIC KEY-----")
	if !bytes.Contains(rawPubKey, want) {
		t.Errorf("want %s, but got %s", want, rawPubKey)
	}
	want = []byte("-----BEGIN PRIVATE KEY-----")
	if !bytes.Contains(rawPrivKey, want) {
		t.Errorf("want %s, but got %s", want, rawPubKey)
	}
}

func TestJWTer_Generate(t *testing.T) {
	c := clock.FixedClocker{}
	ctx := context.Background()
	moq := &StoreMock{}
	wantID := entity.UserID(20)
	u := fixture.User(&entity.User{ID: wantID})
	moq.SaveFunc = func(ctx context.Context, key string, userID entity.UserID) error {
		if userID != wantID {
			t.Errorf("want %d, but got %d", wantID, userID)
		}
		return nil
	}
	sut, err := NewJWTer(moq, c)
	if err != nil {
		t.Fatal(err)
	}
	got, err := sut.GenerateToken(ctx, *u)
	if err != nil {
		t.Fatalf("not want err: %v", err)
	}
	if len(got) == 0 {
		t.Errorf("toke is empty")
	}
}

func TestJWter_GetToken(t *testing.T) {
	t.Parallel()
	// 時刻を固定
	c := clock.FixedClocker{}
	want, err := jwt.NewBuilder().
		JwtID(uuid.New().String()).
		Issuer(`github.com/yuyahy/go_todo_app`).
		Subject("access_token").
		IssuedAt(c.Now()).
		Expiration(c.Now().Add(30*time.Minute)).
		Claim(RoleKey, "test").
		Claim(UserNameKey, "test_user").
		Build()

	if err != nil {
		t.Fatal(err)
	}

	pkey, err := jwk.ParseKey(rawPrivKey, jwk.WithPEM(true))
	if err != nil {
		t.Fatal(err)
	}
	// 署名済みのアクセストークンを生成
	signed, err := jwt.Sign(want, jwt.WithKey(jwa.RS256, pkey))
	if err != nil {
		t.Fatal(err)
	}
	userID := entity.UserID(20)

	ctx := context.Background()
	// 共有メモリへの保存処理をmock
	moq := &StoreMock{}
	moq.LoadFunc = func(ctx context.Context, key string) (entity.UserID, error) {
		return userID, nil
	}
	sut, err := NewJWTer(moq, c)
	if err != nil {
		t.Fatal(err)
	}

	// 生成したアクセストークンはリクエストのAuthorization ヘッダーに埋め込む
	req := httptest.NewRequest(http.MethodGet, `https://github.com/yuyahy`, nil)
	req.Header.Set(`Authorization`, fmt.Sprintf(`Bearer %s`, signed))
	got, err := sut.GetToken(ctx, req)
	if err != nil {
		t.Fatalf("want no error, but got %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetToken() got = %v, want %v", got, want)
	}
}

type FixedTomorrowClocker struct{}

func (c FixedTomorrowClocker) Now() time.Time {
	return clock.FixedClocker{}.Now().Add(24 * time.Hour)
}

func TestJWTer_GetToken_NG(t *testing.T) {
	t.Parallel()

	c := clock.FixedClocker{}
	tok, err := jwt.NewBuilder().
		JwtID(uuid.New().String()).
		Issuer(`github.com/yuyahy/go_todo_app`).
		Subject("access_token").
		IssuedAt(c.Now()).
		Expiration(c.Now().Add(30*time.Minute)).
		Claim(RoleKey, "test").
		Claim(UserNameKey, "test_user").
		Build()
	if err != nil {
		t.Fatal(err)
	}
	pkey, err := jwk.ParseKey(rawPrivKey, jwk.WithPEM(true))
	if err != nil {
		t.Fatal(err)
	}
	signed, err := jwt.Sign(tok, jwt.WithKey(jwa.RS256, pkey))
	if err != nil {
		t.Fatal(err)
	}

	type moq struct {
		userID entity.UserID
		err    error
	}
	tests := map[string]struct {
		c   clock.Clocker
		moq moq
	}{
		"expire": {
			// tokenのexpire時間より未来の時間を返す
			c: FixedTomorrowClocker{},
		},
		"notFoundInStore": {
			// Redis上にデータが残っていなかったケース
			c: clock.FixedClocker{},
			moq: moq{
				err: store.ErrNotFound,
			},
		},
	}
	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			moq := &StoreMock{}
			moq.LoadFunc = func(ctx context.Context, key string) (entity.UserID, error) {
				return tt.moq.userID, tt.moq.err
			}
			sut, err := NewJWTer(moq, tt.c)
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(
				http.MethodGet,
				`https://github.com/yuyahy`,
				nil,
			)
			req.Header.Set(`Authorization`, fmt.Sprintf(`Bearer %s`, signed))
			got, err := sut.GetToken(ctx, req)
			// トークンが期限切れ
			if err == nil {
				t.Errorf("want error but got nil")
			}
			// Redisにトークンが残っていない
			if got != nil {
				t.Errorf("want nil, but got %v", got)
			}
		})
	}
}
