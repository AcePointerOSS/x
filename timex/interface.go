package timex

// this is presented as an interface, allowing us to mock the value during tests.
type TimeInterface interface {
	Millisecond() int64
}

type Timer struct{}

var (
	TimeUtil TimeInterface
)

func init() {
	TimeUtil = &Timer{}
}
