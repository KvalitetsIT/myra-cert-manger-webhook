package testutil

import (
	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/models"
	"github.com/Myra-Security-GmbH/myrasec-go/v2"
)

func MakeRecord() models.Record {
	return models.Record{
		Action:            "present",
		Type:              "dns-01",
		DNSName:           "_acme-challenge.example.com",
		Key:               "key123",
		ResourceNamespace: "default",
		ResolvedFQDN:      "_acme-challenge.example.com.",
		ResolvedZone:      "example.com.",
	}
}

func MakeDNSRecordFromRecord(r models.Record) myrasec.DNSRecord {
	return myrasec.DNSRecord{
		Name:       r.ResolvedFQDN, // include trailing dot
		Value:      r.Key,
		RecordType: r.Type,
		Active:     true,
		Enabled:    true,
		TTL:        60,
	}
}
