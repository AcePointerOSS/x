package timex

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// fixture
var ms int64 = 1613453369077

const LAYOUT = "2006-01-02 15:04:05"

type MockTimerMilisecond struct{}

func (mtm *MockTimerMilisecond) Millisecond() int64 {
	return 1613453369077
}

func (mtm *MockTimerMilisecond) UtcNow() time.Time {
	t, _ := time.Parse(LAYOUT, "2009-11-10 23:00:00")
	return t
}

func TestTimer_Millisecond(t *testing.T) {
	now := time.Now().UnixNano()
	millis := now / 1000000
	// unixnano returns microseconds, milliseconds should execute fast enough so they should be the same.
	require.Equal(t, millis, TimeUtil.Millisecond())
}

// demonstrates how to mock the interface in any calling package.
func TestTimer_MockInterface(t *testing.T) {
	TimeUtil = &MockTimerMilisecond{}
	require.Equal(t, ms, TimeUtil.Millisecond())
}
