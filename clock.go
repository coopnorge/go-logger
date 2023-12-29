package logger

import "time"

type functionClock struct {
	now       func() time.Time
	newTicker func(time.Duration) *time.Ticker
}

func (c functionClock) Now() time.Time {
	if c.now != nil {
		return c.now()
	}
	return time.Now()
}

func (c functionClock) NewTicker(duration time.Duration) *time.Ticker {
	if c.newTicker != nil {
		return c.newTicker(duration)
	}
	return time.NewTicker(duration)
}
