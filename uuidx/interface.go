package uuidx

import (
	"github.com/gofrs/uuid"
)

type UuidInterface interface {
	NewUuid() uuid.UUID
	NewUuidString() string
}

type UuidWrapper struct{}

var (
	UuidUtil UuidInterface
)

func init() {
	UuidUtil = &UuidWrapper{}
}
