package client

import (
	"testing"

	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/configs"
	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/testutil/mocks"
	"github.com/Myra-Security-GmbH/myrasec-go/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func createSampleDomain() myrasec.Domain {
	return myrasec.Domain{
		ID:   1,
		Name: "example.com",
	}
}

/*
	func TestMyraClient_OnAdd(t *testing.T) {
		mockAPI := new(mocks.MockedAPI)
		client := &MyraClient{api: mockAPI, cfg: configs.Myra{}}

		record := sampleDNSRecord()
		domain := createSampleDomain()

		mockAPI.On("ListDomains", mock.Anything).Return([]myrasec.Domain{domain}, nil)
		mockAPI.On("CreateDNSRecord", mock.AnythingOfType("*myrasec.DNSRecord"), domain.ID).Return(&record, nil)

		result, err := client.OnAdd(record)
		require.NoError(t, err, "expected no error from OnAdd")
		require.Equal(t, record.Name, result.Name, "DNS record Name mismatch")
		require.Equal(t, record.Value, result.Value, "DNS record Value mismatch")

		mockAPI.AssertExpectations(t)
	}

	func TestMyraClient_OnAdd_Error(t *testing.T) {
		mockAPI := new(mocks.MockedAPI)
		client := &MyraClient{api: mockAPI, cfg: configs.Myra{}}

		record := sampleDNSRecord()
		domain := createSampleDomain()
		clientErr := errors.New("create failed")

		mockAPI.On("ListDomains", mock.Anything).Return([]myrasec.Domain{domain}, nil)
		mockAPI.On("CreateDNSRecord", mock.AnythingOfType("*myrasec.DNSRecord"), domain.ID).Return(nil, clientErr)

		_, err := client.OnAdd(record)
		require.Error(t, err, "expected an error from OnAdd when client fails")
		require.Contains(t, err.Error(), "create failed", "error message should contain 'create failed'")

		mockAPI.AssertExpectations(t)
	}

	func TestMyraClient_OnDelete(t *testing.T) {
		mockAPI := new(mocks.MockedAPI)
		client := &MyraClient{api: mockAPI, cfg: configs.Myra{}}

		record := sampleDNSRecord()
		domain := createSampleDomain()

		mockAPI.On("ListDomains", mock.Anything).Return([]myrasec.Domain{domain}, nil)
		mockAPI.On("ListDNSRecords", domain.ID, mock.Anything).Return([]myrasec.DNSRecord{record}, nil)
		mockAPI.On("DeleteDNSRecord", mock.AnythingOfType("*myrasec.DNSRecord"), domain.ID).Return(&record, nil)

		result, err := client.OnDelete(record)
		require.NoError(t, err, "expected no error from OnDelete")
		require.Equal(t, record.Name, result.Name, "DNS record Name mismatch after deletion")
		require.Equal(t, record.Value, result.Value, "DNS record Value mismatch after deletion")

		mockAPI.AssertExpectations(t)
	}

	func TestMyraClient_OnDelete_Error(t *testing.T) {
		mockAPI := new(mocks.MockedAPI)
		client := &MyraClient{api: mockAPI, cfg: configs.Myra{}}

		record := sampleDNSRecord()
		domain := createSampleDomain()
		clientErr := errors.New("delete failed")

		mockAPI.On("ListDomains", mock.Anything).Return([]myrasec.Domain{domain}, nil)
		mockAPI.On("ListDNSRecords", domain.ID, record.Name).Return([]myrasec.DNSRecord{record}, nil)
		mockAPI.On("DeleteDNSRecord", &record, domain.ID).Return(nil, clientErr)

		_, err := client.OnDelete(record)
		require.Error(t, err, "expected an error from OnDelete when client fails")
		errorMsg := "Could not aquire domain id associated with the record. A domain id is required in order to delete"
		require.Contains(t, err.Error(), errorMsg, fmt.Sprintf("error message should contain '%s'", errorMsg))

		mockAPI.AssertExpectations(t)
	}

	func TestMyraClient_get_domain_id(t *testing.T) {
		mockAPI := new(mocks.MockedAPI)
		client := &MyraClient{api: mockAPI, cfg: configs.Myra{}}

		domains := []myrasec.Domain{
			{ID: 1, Name: "example.com"},
			{ID: 2, Name: "sub.example.com"},
		}

		mockAPI.On("ListDomains", mock.Anything).Return(domains, nil)

		domainID, err := client.get_domain_id("test.example.com")
		require.NoError(t, err, "expected no error when domain exists")
		require.Equal(t, 1, domainID, "expected domainID 1 for 'test.example.com'")

		subID, err := client.get_domain_id("sub.example.com")
		require.NoError(t, err, "expected no error for exact subdomain match")
		require.Equal(t, 2, subID, "expected domainID 2 for 'sub.example.com'")

		_, err = client.get_domain_id("notfound.com")
		require.Error(t, err, "expected error for non-existent domain")
		require.Contains(t, err.Error(), "Could not derive domain id", "error message should mention domain id derivation failure")
	}
*/
func TestMyraClient_get_record_id(t *testing.T) {
	mockAPI := new(mocks.MockedAPI)
	client := &MyraClient{api: mockAPI, cfg: configs.Myra{}}

	domainID := 1
	records := []myrasec.DNSRecord{
		{ID: 101, Name: "_acme-challenge.example.com."},
		{ID: 102, Name: "_acme-challenge.sub.example.com."},
	}

	mockAPI.On("ListDNSRecords", domainID, mock.Anything).Return(records, nil)

	recordID, err := client.get_record_id(domainID, "_acme-challenge.example.com.")
	require.NoError(t, err, "expected no error when record exists")
	require.Equal(t, 101, recordID, "expected recordID 101 for '_acme-challenge.example.com.'")

	_, err = client.get_record_id(domainID, "_missing.example.com.")
	require.Error(t, err, "expected error for non-existent record")
	require.Contains(t, err.Error(), "Could not derive record id", "error message should mention record id derivation failure")
}

func sampleDNSRecord() myrasec.DNSRecord {
	return myrasec.DNSRecord{
		Name:       "_acme-challenge.example.com.",
		Value:      "key123",
		RecordType: "TXT",
	}
}
