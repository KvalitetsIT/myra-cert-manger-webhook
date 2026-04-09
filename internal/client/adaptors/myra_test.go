package adaptors_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/client/adaptors"
	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/testutil"
	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/testutil/mocks"
	myrasec "github.com/Myra-Security-GmbH/myrasec-go/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMyraClientAdaptor_OnAdd(t *testing.T) {
	mockClient := new(mocks.MockedClient[myrasec.DNSRecord])
	adaptor := adaptors.NewMyraClientAdaptor(mockClient)

	record := testutil.MakeRecord()
	dnsRecord := testutil.MakeDNSRecordFromRecord(record)

	mockClient.On("OnAdd", mock.MatchedBy(func(r myrasec.DNSRecord) bool {
		return r.Name == strings.TrimSuffix(record.ResolvedFQDN, ".") && r.Value == record.Key
	})).Return(dnsRecord, nil)

	result, err := adaptor.OnAdd(record)
	require.NoError(t, err)
	require.Equal(t, dnsRecord.Name, result.DNSName)
	require.Equal(t, record.Key, result.Key)

	mockClient.AssertExpectations(t)
}

func TestMyraClientAdaptor_OnAdd_Error(t *testing.T) {
	mockClient := new(mocks.MockedClient[myrasec.DNSRecord])
	adaptor := adaptors.NewMyraClientAdaptor(mockClient)

	record := testutil.MakeRecord()
	dnsRecord := testutil.MakeDNSRecordFromRecord(record)
	clientErr := errors.New("client failed")

	mockClient.On("OnAdd", mock.MatchedBy(func(r myrasec.DNSRecord) bool {
		return r.Name == strings.TrimSuffix(record.ResolvedFQDN, ".") && r.Value == record.Key
	})).Return(dnsRecord, clientErr)

	_, err := adaptor.OnAdd(record)
	require.Error(t, err)
	require.Equal(t, clientErr, err)

	mockClient.AssertExpectations(t)
}

func TestMyraClientAdaptor_OnDelete(t *testing.T) {
	mockClient := new(mocks.MockedClient[myrasec.DNSRecord])
	adaptor := adaptors.NewMyraClientAdaptor(mockClient)

	record := testutil.MakeRecord()
	dnsRecord := testutil.MakeDNSRecordFromRecord(record)

	mockClient.On("OnDelete", mock.MatchedBy(func(r myrasec.DNSRecord) bool {
		return r.Name == strings.TrimSuffix(record.ResolvedFQDN, ".") && r.Value == record.Key
	})).Return(dnsRecord, nil)

	result, err := adaptor.OnDelete(record)
	require.NoError(t, err)
	require.Equal(t, dnsRecord.Name, result.DNSName)
	require.Equal(t, record.Key, result.Key)

	mockClient.AssertExpectations(t)
}

func TestMyraClientAdaptor_OnDelete_Error(t *testing.T) {
	mockClient := new(mocks.MockedClient[myrasec.DNSRecord])
	adaptor := adaptors.NewMyraClientAdaptor(mockClient)

	record := testutil.MakeRecord()
	dnsRecord := testutil.MakeDNSRecordFromRecord(record)
	clientErr := errors.New("delete failed")

	mockClient.On("OnDelete", mock.MatchedBy(func(r myrasec.DNSRecord) bool {
		return r.Name == strings.TrimSuffix(record.ResolvedFQDN, ".") && r.Value == record.Key
	})).Return(dnsRecord, clientErr)

	_, err := adaptor.OnDelete(record)
	require.Error(t, err)
	require.Equal(t, clientErr, err)

	mockClient.AssertExpectations(t)
}
