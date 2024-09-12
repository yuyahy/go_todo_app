package fixture

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/yuyahy/go_todo_app/entity"
)

// テスト用の*entity.Userを準備するfixture
func User(u *entity.User) *entity.User {
	// テスト用User
	result := &entity.User{
		ID:       entity.UserID(rand.Int()),
		Name:     "yuyahy" + strconv.Itoa(rand.Int())[:5],
		Password: "password",
		Role:     "admin",
		Created:  time.Now(),
		Modified: time.Now(),
	}
	if u == nil {
		return result
	}
	// 引数のUserのフィールドに何らかの値が設定されていた場合は、
	// その値でフィールドを更新する
	if u.ID != 0 {
		result.ID = u.ID
	}
	if u.Name != "" {
		result.Name = u.Name
	}
	if u.Password != "" {
		result.Password = u.Password
	}
	if u.Role != "" {
		result.Role = u.Role
	}
	if !u.Created.IsZero() {
		result.Created = u.Created
	}
	if !u.Modified.IsZero() {
		result.Modified = u.Modified
	}
	return result
}
