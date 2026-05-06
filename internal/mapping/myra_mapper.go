package mapping

import (
	"strings"

	"github.com/KvalitetsIT/cert-manager-webhook-core/pkg/models"
	"github.com/Myra-Security-GmbH/myrasec-go/v2"
)

type MyraMapper struct{}

func (m MyraMapper) ToExternal(record models.Record) myrasec.DNSRecord {
	return myrasec.DNSRecord{
		ID:         0,
		Name:       strings.TrimSuffix(record.ResolvedFQDN, "."),
		Value:      record.Key,
		RecordType: record.Type,
		Active:     true,
		Enabled:    true,
		TTL:        60,
	}
}

func (m MyraMapper) ToInternal(r myrasec.DNSRecord) models.Record {
	return models.Record{
		DNSName: r.Name,
		Key:     r.Value,
		Type:    r.RecordType,
	}
}
