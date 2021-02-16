package uuidx

import (
	"github.com/gofrs/uuid"
	guuid "github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

// fixture
var uuid4 string = "306ca21a-e7da-41ca-8e00-350d6a549ab6"

type UuidWrapperMock struct{}

func (uwm *UuidWrapperMock) NewUuid() uuid.UUID {
	return uuid.Must(uuid.FromString(uuid4))
}

func (uwm *UuidWrapperMock) NewUuidString() string {
	return uuid4
}

func TestUuidWrapper_NewUuid(t *testing.T) {
	u := UuidUtil.NewUuid()
	_, err := guuid.Parse(u.String())
	require.NoError(t, err)
}

func TestUuidWrapper_NewUuidString(t *testing.T) {
	u := UuidUtil.NewUuidString()
	_, err := guuid.Parse(u)
	require.NoError(t, err)
}

// demonstrates how to mock the interface via overriding during tests,
func TestUuidWrapper_Mock(t *testing.T) {
	UuidUtil = &UuidWrapperMock{}
	require.Equal(t, uuid4, UuidUtil.NewUuid().String())
	require.Equal(t, uuid4, UuidUtil.NewUuidString())
}
