package mapping_test

import (
	"testing"

	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/mapping"
	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/models"
	myrasec "github.com/Myra-Security-GmbH/myrasec-go/v2"
	"github.com/stretchr/testify/require"
)

func TestMyraMapper_ToExternal(t *testing.T) {
	mapper := &mapping.MyraMapper{}

	record := models.Record{
		DNSName:      "_acme-challenge.example.com.",
		Key:          "key123",
		Type:         "dns-01",
		ResolvedFQDN: "_acme-challenge.example.com.",
	}

	expected := myrasec.DNSRecord{
		ID:         0,
		Name:       record.ResolvedFQDN,
		Value:      record.Key,
		RecordType: record.Type,
		Active:     true,
		Enabled:    true,
		TTL:        60,
	}

	result := mapper.ToExternal(record)

	require.Equal(t, expected, result)
}

func TestMyraMapper_ToInternal(t *testing.T) {
	mapper := &mapping.MyraMapper{}

	dnsRecord := myrasec.DNSRecord{
		ID:         1,
		Name:       "_acme-challenge.example.com.",
		Value:      "key123",
		RecordType: "dns-01",
		Active:     true,
		Enabled:    true,
		TTL:        60,
	}

	expected := models.Record{
		DNSName: "_acme-challenge.example.com.",
		Key:     "key123",
		Type:    "dns-01",
	}

	result := mapper.ToInternal(dnsRecord)

	require.Equal(t, expected, result)
}
