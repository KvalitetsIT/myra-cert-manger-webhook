package testutil_test

import (
	"log/slog"
	"testing"

	"github.com/KvalitetsIT/cert-manager-webhook-myra/internal/testutil"
	myrasec "github.com/Myra-Security-GmbH/myrasec-go/v2"
	"github.com/stretchr/testify/require"
)

func newTestStorage() *testutil.Storage {
	return testutil.NewStorage(slog.Default())
}

func TestAddAndGetDomain(t *testing.T) {
	s := newTestStorage()

	domain := myrasec.Domain{Name: "example.com"}
	id := s.AddDomain(domain)

	require.Equal(t, 0, id)

	domains := s.GetDomains()
	require.Len(t, domains, 1)
	require.Equal(t, "example.com", domains[0].Name)
	require.Equal(t, id, domains[0].ID)
}

func TestGetDomainID_ExactMatch(t *testing.T) {
	s := newTestStorage()

	id := s.AddDomain(myrasec.Domain{Name: "example.com"})

	got, found := s.GetDomainID("example.com")
	require.True(t, found)
	require.Equal(t, id, got)
}

func TestGetDomainID_TrailingDot(t *testing.T) {
	s := newTestStorage()

	id := s.AddDomain(myrasec.Domain{Name: "example.com"})

	got, found := s.GetDomainID("example.com.")
	require.True(t, found)
	require.Equal(t, id, got)
}

func TestGetDomainID_SubdomainMatch(t *testing.T) {
	s := newTestStorage()

	rootID := s.AddDomain(myrasec.Domain{Name: "example.com"})

	got, found := s.GetDomainID("test.example.com")
	require.True(t, found)
	require.Equal(t, rootID, got)
}

func TestGetDomainID_LongestSuffixMatch(t *testing.T) {
	s := newTestStorage()

	rootID := s.AddDomain(myrasec.Domain{Name: "example.com"})
	subID := s.AddDomain(myrasec.Domain{Name: "sub.example.com"})

	got, found := s.GetDomainID("a.sub.example.com")
	require.True(t, found)
	require.Equal(t, subID, got)

	require.NotEqual(t, rootID, got)
}

func TestGetDomainID_NotFound(t *testing.T) {
	s := newTestStorage()

	_, found := s.GetDomainID("unknown.com")
	require.False(t, found)
}

func TestAddAndGetRecord(t *testing.T) {
	s := newTestStorage()

	domainID := s.AddDomain(myrasec.Domain{Name: "example.com"})

	record := myrasec.DNSRecord{Name: "www"}
	added, err := s.AddRecord(domainID, record)
	require.NoError(t, err)

	got, found := s.GetRecord(domainID, added.ID)
	require.True(t, found)
	require.Equal(t, "www", got.Name)
}

func TestGetRecords(t *testing.T) {
	s := newTestStorage()

	domainID := s.AddDomain(myrasec.Domain{Name: "example.com"})

	_, _ = s.AddRecord(domainID, myrasec.DNSRecord{Name: "a"})
	_, _ = s.AddRecord(domainID, myrasec.DNSRecord{Name: "b"})

	records, err := s.GetRecords(domainID)
	require.NoError(t, err)
	require.Len(t, records, 2)
}

func TestGetRecords_DomainNotFound(t *testing.T) {
	s := newTestStorage()

	records, err := s.GetRecords(999)

	require.Error(t, err)
	require.Empty(t, records)
}

func TestDeleteRecord(t *testing.T) {
	s := newTestStorage()

	domainID := s.AddDomain(myrasec.Domain{Name: "example.com"})
	rec, _ := s.AddRecord(domainID, myrasec.DNSRecord{Name: "www"})

	deleted, err := s.DeleteRecord(domainID, rec.ID)
	require.NoError(t, err)
	require.Equal(t, rec.ID, deleted.ID)

	_, found := s.GetRecord(domainID, rec.ID)
	require.False(t, found)
}

func TestDeleteRecord_NotFound(t *testing.T) {
	s := newTestStorage()

	domainID := s.AddDomain(myrasec.Domain{Name: "example.com"})

	_, err := s.DeleteRecord(domainID, 123)
	require.Error(t, err)
}

func TestGetRecordID(t *testing.T) {
	s := newTestStorage()

	domainID := s.AddDomain(myrasec.Domain{Name: "example.com"})
	rec, _ := s.AddRecord(domainID, myrasec.DNSRecord{Name: "www"})

	id, found := s.GetRecordID("www")
	require.True(t, found)
	require.Equal(t, rec.ID, id)
}

func TestGetRecordID_NotFound(t *testing.T) {
	s := newTestStorage()

	_, found := s.GetRecordID("missing")
	require.False(t, found)
}

func TestCounter(t *testing.T) {
	var c testutil.Counter

	require.Equal(t, 0, c.Next())
	require.Equal(t, 1, c.Next())
	require.Equal(t, 2, c.Next())
}
