package envx

import (
	"os"
)

type Environment interface {
	Getenv(key string) string
	InitMap()
	Setenv(key string, val string) error
}

type Env struct{}

func (e *Env) Getenv(key string) string {
	return os.Getenv(key)
}

// This is so that we can call these functions in our mocks, we dont implement them here
// but its implemented in our mock package.
func (e *Env) InitMap() {}

func (e *Env) Setenv(key string, val string) error {
	return os.Setenv(key, val)
}

var (
	OSEnv Environment
)

func init() {
	OSEnv = &Env{}
}
