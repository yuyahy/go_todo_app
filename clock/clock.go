package clock

import "time"

type Clocker interface {
	Now() time.Time
}

// アプリで実際に使用する時刻情報
type RealClocker struct{}

func (r RealClocker) Now() time.Time {
	return time.Now()
}

// テストで使用する固定時刻情報
type FixedClocker struct{}

func (fc FixedClocker) Now() time.Time {
	return time.Date(2022, 5, 10, 12, 34, 56, 0, time.UTC)
}
