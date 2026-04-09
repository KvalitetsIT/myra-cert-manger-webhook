package client_test

import (
	"errors"
	"log/slog"
	"testing"

	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/client"
	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/testutil/mocks"
	myrasec "github.com/Myra-Security-GmbH/myrasec-go/v2"
	"github.com/stretchr/testify/require"
)

func TestClientLogger_OnAdd_Success(t *testing.T) {
	mockClient := new(mocks.MockedClient[myrasec.DNSRecord])
	cl := client.NewClientLogger[myrasec.DNSRecord](mockClient, slog.Default())

	record := myrasec.DNSRecord{Name: "_acme-challenge.example.com.", Value: "key123"}
	mockClient.On("OnAdd", record).Return(record, nil)

	result, err := cl.OnAdd(record)
	require.NoError(t, err)
	require.Equal(t, record, result)

	mockClient.AssertExpectations(t)
}

func TestClientLogger_OnAdd_Error(t *testing.T) {
	mockClient := new(mocks.MockedClient[myrasec.DNSRecord])
	cl := client.NewClientLogger[myrasec.DNSRecord](mockClient, slog.Default())

	record := myrasec.DNSRecord{Name: "_acme-challenge.example.com.", Value: "key123"}
	addErr := errors.New("add failed")
	mockClient.On("OnAdd", record).Return(myrasec.DNSRecord{}, addErr)

	_, err := cl.OnAdd(record)
	require.Error(t, err)
	require.Equal(t, addErr, err)

	mockClient.AssertExpectations(t)
}

func TestClientLogger_OnDelete_Success(t *testing.T) {
	mockClient := new(mocks.MockedClient[myrasec.DNSRecord])
	cl := client.NewClientLogger[myrasec.DNSRecord](mockClient, slog.Default())

	record := myrasec.DNSRecord{Name: "_acme-challenge.example.com.", Value: "key123"}
	mockClient.On("OnDelete", record).Return(record, nil)

	result, err := cl.OnDelete(record)
	require.NoError(t, err)
	require.Equal(t, record, result)

	mockClient.AssertExpectations(t)
}

func TestClientLogger_OnDelete_Error(t *testing.T) {
	mockClient := new(mocks.MockedClient[myrasec.DNSRecord])
	cl := client.NewClientLogger[myrasec.DNSRecord](mockClient, slog.Default())

	record := myrasec.DNSRecord{Name: "_acme-challenge.example.com.", Value: "key123"}
	deleteErr := errors.New("delete failed")
	mockClient.On("OnDelete", record).Return(myrasec.DNSRecord{}, deleteErr)

	_, err := cl.OnDelete(record)
	require.Error(t, err)
	require.Equal(t, deleteErr, err)

	mockClient.AssertExpectations(t)
}
