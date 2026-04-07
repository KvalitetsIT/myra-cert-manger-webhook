package testutil

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/Myra-Security-GmbH/myrasec-go/v2"
)

type Storage struct {
	domains        *AtomicMap[int, myrasec.Domain]                     // domainID -> domain
	records        *AtomicMap[int, *AtomicMap[int, myrasec.DNSRecord]] // domainID -> recordID -> record
	domainIds      *AtomicMap[string, int]                             // domainName -> domainID
	recordIds      *AtomicMap[string, int]                             // recordName -> recordID
	domain_counter *Counter
	record_counter *Counter
}

// NewStorage initializes a Storage instance
func NewStorage(logger *slog.Logger) *Storage {
	dc := Counter(0)
	rc := Counter(0)
	return &Storage{
		domains:        NewAtomicMap[int, myrasec.Domain]("domains"),
		records:        NewAtomicMap[int, *AtomicMap[int, myrasec.DNSRecord]]("records"),
		domainIds:      NewAtomicMap[string, int]("domainId"),
		recordIds:      NewAtomicMap[string, int]("recordIds"),
		domain_counter: &dc,
		record_counter: &rc,
	}
}

// GetDomains returns all domains
func (s *Storage) GetDomains() []*myrasec.Domain {
	result := []*myrasec.Domain{}
	for _, domain := range s.domains.GetAll() {
		result = append(result, &domain)
	}
	return result
}

// GetRecords returns all records for a given domainID
func (s *Storage) GetRecords(domainID int) ([]myrasec.DNSRecord, error) {
	recordMap, found := s.records.Get(domainID)
	if !found {
		// collect existing domain IDs
		var ids []int
		for id := range s.records.GetAll() {
			ids = append(ids, id)
		}

		return []myrasec.DNSRecord{}, fmt.Errorf(
			"domainID %d not found. existing domainIDs: %v",
			domainID,
			ids,
		)
	}
	return recordMap.GetAll(), nil
}

// GetRecord returns a record in a certain domain
func (s *Storage) GetRecord(domainID, recordId int) (myrasec.DNSRecord, bool) {
	domain, found := s.records.Get(domainID)

	if !found {
		return myrasec.DNSRecord{}, false
	}
	record, found := domain.Get(recordId)

	if !found {
		return myrasec.DNSRecord{}, false
	}
	return record, true
}

// GetDomainID returns the domainID for a given domain name
func (s *Storage) GetDomainID(name string) (int, bool) {
	// 1. Clean the input (remove trailing dot if present)
	query := strings.TrimSuffix(name, ".")

	// 2. Try an exact match first for performance
	if id, found := s.domainIds.Get(query); found {
		return id, true
	}

	// 3. Fallback: Search for the parent domain (suffix matching)
	// We iterate through all registered domain names
	var bestMatchID int
	var bestMatchName string
	found := false

	// GetAll() on domainIds returns map[string]int
	for domainName, id := range s.domainIds.entries {
		// Check if query is a subdomain (e.g., "test.example.com" ends with ".example.com")
		if strings.HasSuffix(query, "."+domainName) {
			// In case of overlapping domains (e.g., "example.com" and "sub.example.com")
			// we want the longest match.
			if len(domainName) > len(bestMatchName) {
				bestMatchName = domainName
				bestMatchID = id
				found = true
			}
		}
	}

	if !found {
		return -1, false
	}

	return bestMatchID, found
}

// GetRecordID returns the recordID for a given record name
func (s *Storage) GetRecordID(name string) (int, bool) {

	if id, found := s.recordIds.Get(name); found {
		return id, true
	}

	var bestMatchID int
	var bestMatchName string
	found := false

	for recordName, id := range s.recordIds.entries {
		if name == recordName {
			if len(recordName) > len(bestMatchName) {
				bestMatchName = recordName
				bestMatchID = id
				found = true
			}
		}
	}
	return bestMatchID, found
}

// AddDomain adds a new domain and returns its ID
func (s *Storage) AddDomain(domain myrasec.Domain) int {
	domainID := s.domain_counter.Next()
	domain.ID = domainID
	s.domains.Set(domainID, domain)
	s.domainIds.Set(domain.Name, domainID)

	// Initialize empty record map for this domain
	s.records.Set(domainID, NewAtomicMap[int, myrasec.DNSRecord](fmt.Sprintf("domain/%d/records", domainID)))

	return domainID
}

// AddRecord adds a DNS record to a domain
func (s *Storage) AddRecord(domainID int, record myrasec.DNSRecord) (myrasec.DNSRecord, error) {
	record.ID = s.record_counter.Next()

	recordMap, found := s.records.Get(domainID)
	if !found {
		recordMap = NewAtomicMap[int, myrasec.DNSRecord](fmt.Sprintf("domain/%d/records", domainID))
		s.records.Set(domainID, recordMap)
	}

	recordMap.Set(record.ID, record)
	s.recordIds.Set(record.Name, record.ID)

	return record, nil
}

// DeleteRecord deletes a record from a domain
func (s *Storage) DeleteRecord(domainID, recordID int) (myrasec.DNSRecord, error) {
	recordMap, found := s.records.Get(domainID)
	if !found {
		return myrasec.DNSRecord{}, fmt.Errorf("domainID %d not found", domainID)
	}

	record, deleted := recordMap.Delete(recordID)
	if !deleted {
		return myrasec.DNSRecord{}, fmt.Errorf("recordID %d not found in domain %d", recordID, domainID)
	}

	// Remove record index
	s.recordIds.Delete(record.Name)

	return record, nil
}

// Counter for generating incremental IDs
type Counter int

func (c *Counter) Next() int {
	current := int(*c)
	*c = *c + 1
	return current
}
