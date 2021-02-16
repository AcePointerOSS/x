package uuidx

import (
	"github.com/gofrs/uuid"
)

func (u *UuidWrapper) NewUuid() uuid.UUID {
	return uuid.Must(uuid.NewV4())
}
func (u *UuidWrapper) NewUuidString() string {
	return uuid.Must(uuid.NewV4()).String()
}
