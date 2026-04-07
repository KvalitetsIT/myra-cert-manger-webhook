package adaptors

import (
	"errors"
	"testing"

	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/testutil/mocks"
	"github.com/stretchr/testify/require"
)

type Internal struct {
	Name string
}

type External struct {
	Name string
}

func TestClientAdaptor_OnAdd(t *testing.T) {
	mockClient := new(mocks.MockedClient[External])
	mockMapper := new(mocks.MockedMapper[Internal, External])

	adaptor := &clientAdaptor[Internal, External]{client: mockClient, mapper: mockMapper}

	internalRec := Internal{Name: "test"}
	externalRec := External{Name: "test"}

	// Mapper mocks
	mockMapper.On("ToExternal", internalRec).Return(externalRec)
	mockMapper.On("ToInternal", externalRec).Return(internalRec)

	// Client mock
	mockClient.On("OnAdd", externalRec).Return(externalRec, nil)

	result, err := adaptor.OnAdd(internalRec)
	require.NoError(t, err, "expected no error from OnAdd")
	require.Equal(t, internalRec, result, "internal record should map back correctly")

	mockMapper.AssertExpectations(t)
	mockClient.AssertExpectations(t)
}

func TestClientAdaptor_OnAdd_Error(t *testing.T) {
	mockClient := new(mocks.MockedClient[External])
	mockMapper := new(mocks.MockedMapper[Internal, External])

	adaptor := &clientAdaptor[Internal, External]{client: mockClient, mapper: mockMapper}

	internalRec := Internal{Name: "test"}
	externalRec := External{Name: "test"}
	clientErr := errors.New("client failed")

	mockMapper.On("ToExternal", internalRec).Return(externalRec)
	mockMapper.On("ToInternal", externalRec).Return(internalRec)

	mockClient.On("OnAdd", externalRec).Return(externalRec, clientErr)

	result, err := adaptor.OnAdd(internalRec)
	require.Error(t, err, "expected error from OnAdd when client fails")
	require.Equal(t, clientErr, err, "error should match client error")
	require.Equal(t, internalRec, result, "internal record should still map back on error")

	mockMapper.AssertExpectations(t)
	mockClient.AssertExpectations(t)
}

func TestClientAdaptor_OnDelete(t *testing.T) {
	mockClient := new(mocks.MockedClient[External])
	mockMapper := new(mocks.MockedMapper[Internal, External])

	adaptor := &clientAdaptor[Internal, External]{client: mockClient, mapper: mockMapper}

	internalRec := Internal{Name: "test"}
	externalRec := External{Name: "test"}

	mockMapper.On("ToExternal", internalRec).Return(externalRec)
	mockMapper.On("ToInternal", externalRec).Return(internalRec)

	mockClient.On("OnDelete", externalRec).Return(externalRec, nil)

	result, err := adaptor.OnDelete(internalRec)
	require.NoError(t, err, "expected no error from OnDelete")
	require.Equal(t, internalRec, result, "internal record should map back correctly")

	mockMapper.AssertExpectations(t)
	mockClient.AssertExpectations(t)
}

func TestClientAdaptor_OnDelete_Error(t *testing.T) {
	mockClient := new(mocks.MockedClient[External])
	mockMapper := new(mocks.MockedMapper[Internal, External])

	adaptor := &clientAdaptor[Internal, External]{client: mockClient, mapper: mockMapper}

	internalRec := Internal{Name: "test"}
	externalRec := External{Name: "test"}
	clientErr := errors.New("delete failed")

	mockMapper.On("ToExternal", internalRec).Return(externalRec)
	mockMapper.On("ToInternal", externalRec).Return(internalRec)

	mockClient.On("OnDelete", externalRec).Return(externalRec, clientErr)

	result, err := adaptor.OnDelete(internalRec)
	require.Error(t, err, "expected error from OnDelete when client fails")
	require.Equal(t, clientErr, err, "error should match client error")
	require.Equal(t, internalRec, result, "internal record should still map back on error")

	mockMapper.AssertExpectations(t)
	mockClient.AssertExpectations(t)
}
