package client

import (
	myrasec "github.com/Myra-Security-GmbH/myrasec-go/v2"
	"github.com/stretchr/testify/mock"
)

type MockMyraClient struct {
	mock.Mock
}

func (m *MockMyraClient) OnAdd(r myrasec.DNSRecord) (myrasec.DNSRecord, error) {
	args := m.Called(r)
	return args.Get(0).(myrasec.DNSRecord), args.Error(1)
}

func (m *MockMyraClient) OnDelete(r myrasec.DNSRecord) (myrasec.DNSRecord, error) {
	args := m.Called(r)
	return args.Get(0).(myrasec.DNSRecord), args.Error(1)
}
