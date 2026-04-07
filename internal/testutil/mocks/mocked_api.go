package mocks

import (
	myrasec "github.com/Myra-Security-GmbH/myrasec-go/v2"
	"github.com/stretchr/testify/mock"
)

// MockedAPI implements MyraAPI for unit testing
type MockedAPI struct {
	mock.Mock
}

func (m *MockedAPI) ListDomains(params map[string]string) ([]myrasec.Domain, error) {
	args := m.Called(params)
	return args.Get(0).([]myrasec.Domain), args.Error(1)
}

func (m *MockedAPI) ListDNSRecords(domainId int, params map[string]string) ([]myrasec.DNSRecord, error) {
	args := m.Called(domainId, params)
	return args.Get(0).([]myrasec.DNSRecord), args.Error(1)
}

func (m *MockedAPI) CreateDNSRecord(record *myrasec.DNSRecord, domainId int) (*myrasec.DNSRecord, error) {
	args := m.Called(record, domainId)
	return args.Get(0).(*myrasec.DNSRecord), args.Error(1)
}

func (m *MockedAPI) DeleteDNSRecord(record *myrasec.DNSRecord, domainId int) (*myrasec.DNSRecord, error) {
	args := m.Called(record, domainId)
	return args.Get(0).(*myrasec.DNSRecord), args.Error(1)
}
