package timex

import "time"

func (tt *Timer) Millisecond() int64 {
	now := time.Now()
	nanos := now.UnixNano()
	millis := nanos / 1000000
	return millis
}
