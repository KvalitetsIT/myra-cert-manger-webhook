package client

import (
	"errors"
	"log/slog"
	"testing"

	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/configs"
	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/testutil/mocks"
	myrasec "github.com/Myra-Security-GmbH/myrasec-go/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func createSampleDomain() myrasec.Domain {
	return myrasec.Domain{
		ID:   1,
		Name: "example.com",
	}
}

func sampleDNSRecord() myrasec.DNSRecord {
	return myrasec.DNSRecord{
		Name:       "_acme-challenge.example.com.",
		Value:      "key123",
		RecordType: "TXT",
	}
}

// --- OnAdd ---

func TestMyraClient_OnAdd(t *testing.T) {
	mockAPI := new(mocks.MockedAPI)
	c := &MyraClient{api: mockAPI, cfg: configs.Myra{}, logger: slog.Default()}

	record := sampleDNSRecord()
	domain := createSampleDomain()

	mockAPI.On("ListDomains", mock.Anything).Return([]myrasec.Domain{domain}, nil)
	mockAPI.On("CreateDNSRecord", mock.AnythingOfType("*myrasec.DNSRecord"), domain.ID).Return(&record, nil)

	result, err := c.OnAdd(record)
	require.NoError(t, err)
	require.Equal(t, record.Name, result.Name)
	require.Equal(t, record.Value, result.Value)

	mockAPI.AssertExpectations(t)
}

func TestMyraClient_OnAdd_Error(t *testing.T) {
	mockAPI := new(mocks.MockedAPI)
	c := &MyraClient{api: mockAPI, cfg: configs.Myra{}, logger: slog.Default()}

	record := sampleDNSRecord()
	domain := createSampleDomain()
	clientErr := errors.New("create failed")

	mockAPI.On("ListDomains", mock.Anything).Return([]myrasec.Domain{domain}, nil)
	mockAPI.On("CreateDNSRecord", mock.AnythingOfType("*myrasec.DNSRecord"), domain.ID).Return((*myrasec.DNSRecord)(nil), clientErr)

	_, err := c.OnAdd(record)
	require.Error(t, err)
	require.Contains(t, err.Error(), "create failed")

	mockAPI.AssertExpectations(t)
}

func TestMyraClient_OnAdd_DomainNotFound(t *testing.T) {
	mockAPI := new(mocks.MockedAPI)
	c := &MyraClient{api: mockAPI, cfg: configs.Myra{}, logger: slog.Default()}

	record := sampleDNSRecord()

	mockAPI.On("ListDomains", mock.Anything).Return([]myrasec.Domain{}, nil)

	_, err := c.OnAdd(record)
	require.Error(t, err)
	require.Contains(t, err.Error(), "A domain id is required")

	mockAPI.AssertExpectations(t)
}

// --- OnDelete ---

func TestMyraClient_OnDelete(t *testing.T) {
	mockAPI := new(mocks.MockedAPI)
	c := &MyraClient{api: mockAPI, cfg: configs.Myra{}, logger: slog.Default()}

	record := sampleDNSRecord()
	domain := createSampleDomain()

	mockAPI.On("ListDomains", mock.Anything).Return([]myrasec.Domain{domain}, nil)
	mockAPI.On("ListDNSRecords", domain.ID, mock.Anything).Return([]myrasec.DNSRecord{record}, nil)
	mockAPI.On("DeleteDNSRecord", mock.AnythingOfType("*myrasec.DNSRecord"), domain.ID).Return(&record, nil)

	result, err := c.OnDelete(record)
	require.NoError(t, err)
	require.Equal(t, record.Name, result.Name)
	require.Equal(t, record.Value, result.Value)

	mockAPI.AssertExpectations(t)
}

func TestMyraClient_OnDelete_DomainNotFound(t *testing.T) {
	mockAPI := new(mocks.MockedAPI)
	c := &MyraClient{api: mockAPI, cfg: configs.Myra{}, logger: slog.Default()}

	record := sampleDNSRecord()

	mockAPI.On("ListDomains", mock.Anything).Return([]myrasec.Domain{}, nil)

	_, err := c.OnDelete(record)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Could not aquire domain id")

	mockAPI.AssertExpectations(t)
}

func TestMyraClient_OnDelete_RecordAlreadyGone(t *testing.T) {
	mockAPI := new(mocks.MockedAPI)
	c := &MyraClient{api: mockAPI, cfg: configs.Myra{}, logger: slog.Default()}

	record := sampleDNSRecord()
	domain := createSampleDomain()

	mockAPI.On("ListDomains", mock.Anything).Return([]myrasec.Domain{domain}, nil)
	mockAPI.On("ListDNSRecords", domain.ID, mock.Anything).Return([]myrasec.DNSRecord{}, nil)

	result, err := c.OnDelete(record)
	require.NoError(t, err)
	require.Equal(t, record.Name, result.Name)

	mockAPI.AssertExpectations(t)
}

func TestMyraClient_OnDelete_DeleteFails(t *testing.T) {
	mockAPI := new(mocks.MockedAPI)
	c := &MyraClient{api: mockAPI, cfg: configs.Myra{}, logger: slog.Default()}

	record := sampleDNSRecord()
	domain := createSampleDomain()
	deleteErr := errors.New("delete failed")

	mockAPI.On("ListDomains", mock.Anything).Return([]myrasec.Domain{domain}, nil)
	mockAPI.On("ListDNSRecords", domain.ID, mock.Anything).Return([]myrasec.DNSRecord{record}, nil)
	mockAPI.On("DeleteDNSRecord", mock.AnythingOfType("*myrasec.DNSRecord"), domain.ID).Return((*myrasec.DNSRecord)(nil), deleteErr)

	_, err := c.OnDelete(record)
	require.Error(t, err)
	require.Contains(t, err.Error(), "delete failed")

	mockAPI.AssertExpectations(t)
}

// --- get_record_id ---

func TestMyraClient_get_record_id(t *testing.T) {
	mockAPI := new(mocks.MockedAPI)
	c := &MyraClient{api: mockAPI, cfg: configs.Myra{}, logger: slog.Default()}

	domainID := 1
	records := []myrasec.DNSRecord{
		{ID: 101, Name: "_acme-challenge.example.com."},
		{ID: 102, Name: "_acme-challenge.sub.example.com."},
	}

	mockAPI.On("ListDNSRecords", domainID, mock.Anything).Return(records, nil)

	recordID, err := c.get_record_id(domainID, "_acme-challenge.example.com.")
	require.NoError(t, err)
	require.Equal(t, 101, recordID)

	_, err = c.get_record_id(domainID, "_missing.example.com.")
	require.Error(t, err)
	require.Contains(t, err.Error(), "Could not derive record id")
}

// --- get_domain_id ---

func TestMyraClient_get_domain_id(t *testing.T) {
	mockAPI := new(mocks.MockedAPI)
	c := &MyraClient{api: mockAPI, cfg: configs.Myra{}, logger: slog.Default()}

	domains := []myrasec.Domain{
		{ID: 1, Name: "example.com"},
		{ID: 2, Name: "other.dk"},
	}

	mockAPI.On("ListDomains", mock.Anything).Return(domains, nil)

	id, err := c.get_domain_id("_acme-challenge.example.com.")
	require.NoError(t, err)
	require.Equal(t, 1, id)

	id, err = c.get_domain_id("_acme-challenge.other.dk.")
	require.NoError(t, err)
	require.Equal(t, 2, id)

	_, err = c.get_domain_id("notfound.io.")
	require.Error(t, err)
	require.Contains(t, err.Error(), "Could not derive domain id")

	mockAPI.AssertExpectations(t)
}

// --- extractTopDomain ---

func TestExtractTopDomain(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		wantErr  bool
	}{
		{"_acme-challenge.example.com.", "example.com", false},
		{"www.sub.example.com.", "example.com", false},
		{"example.com.", "example.com", false},
		{"_acme-challenge.example.com", "example.com", false},
		{"single", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := extractTopDomain(tt.input)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, result)
			}
		})
	}
}
