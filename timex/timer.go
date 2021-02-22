package timex

import "time"

func (tt *Timer) UtcNow() time.Time {
	return time.Now().UTC()
}
