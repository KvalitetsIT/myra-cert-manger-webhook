package main_test

import (
	"fmt"

	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/models"
	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/testutil"
	"github.com/Myra-Security-GmbH/myrasec-go/v2"
)

type mockedClient struct {
	Added   []models.Record
	Deleted []models.Record
	storage *testutil.Storage
}

func NewMockedClient(storage *testutil.Storage) *mockedClient {
	return &mockedClient{
		Added:   []models.Record{},
		Deleted: []models.Record{},
		storage: storage,
	}
}

func (m *mockedClient) OnDelete(record models.Record) (models.Record, error) {
	m.Deleted = append(m.Deleted, record)
	domainID, found := m.storage.GetDomainID(record.ResolvedZone)
	if !found {
		return record, nil
	}

	recordID, found := m.storage.GetRecordID(record.ResolvedFQDN)
	if !found {
		return record, nil
	}

	_, err := m.storage.DeleteRecord(domainID, recordID)
	if err != nil {
		return record, fmt.Errorf("failed to delete from mock storage: %w", err)
	}

	return record, nil
}

func (m *mockedClient) OnAdd(record models.Record) (models.Record, error) {
	m.Added = append(m.Added, record)

	domainID, found := m.storage.GetDomainID(record.ResolvedZone)
	if !found {
		domainID = m.storage.AddDomain(myrasec.Domain{Name: "example.com"})
	}

	m.storage.AddRecord(domainID, myrasec.DNSRecord{
		Name:       record.ResolvedFQDN,
		Value:      record.Key,
		RecordType: "TXT",
	})

	return record, nil
}
