package clock

import "time"

type Clocker interface {
	Now() time.Time
}

// アプリケーションで実際の時刻を返す time.Now() 関数のラッパ関数
type RealClocker struct{}

func (r RealClocker) Now() time.Time {
	return time.Now()
}

// テスト用に固定の時刻を返す関数
type FixedClocker struct{}

func (fc FixedClocker) Now() time.Time {
	return time.Date(2024, 9, 24, 12, 34, 56, 0, time.UTC)
}
