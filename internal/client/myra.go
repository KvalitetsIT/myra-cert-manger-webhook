package client

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/configs"
	myrasec "github.com/Myra-Security-GmbH/myrasec-go/v2"
)

type MyraClient struct {
	api    MyraAPI
	cfg    configs.Myra
	logger *slog.Logger
}

func NewMyraClient(cfg configs.Myra, logger *slog.Logger) (*MyraClient, error) {
	api, err := myrasec.NewWithToken(cfg.Api.Token)

	if err != nil {
		api, err = myrasec.New(cfg.Api.Key, cfg.Api.Secret)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Myra client; %w", err)
		}
	}

	// Defaults to "https://apiv2.myracloud.com/%s"
	if cfg.Api.URL != "" {
		api.BaseURL = strings.TrimRight(cfg.Api.URL, "/") + "/%s"
	}

	return &MyraClient{
		api:    api,
		cfg:    cfg,
		logger: logger,
	}, nil
}

// Deletes the given domain in the Myra Dns Server
// Returns the deleted record otherwise error
func (c *MyraClient) OnDelete(record myrasec.DNSRecord) (myrasec.DNSRecord, error) {

	domainId, err := c.get_domain_id(record.Name)
	if err != nil {
		return myrasec.DNSRecord{}, fmt.Errorf("Could not aquire domain id associated with the record. A domain id is required in order to delete '%s'; %w", record.Name, err)
	}

	recordId, err := c.get_record_id(domainId, record.Name)
	if err != nil {
		// If the record isn't found, it might have been "cleaned up" already
		// Ignoring the error and returning record
		return record, nil
	}

	record.ID = recordId // This can be done since the record won't have an id from cert manager

	c.logger.Info("deleting record", slog.Any("record", record), slog.Any("domainid", domainId))
	deletedRecord, err := c.api.DeleteDNSRecord(&record, domainId)

	if err != nil {
		c.logger.Error("Failed to delete record", slog.Any("record", record), slog.Any("domainId", domainId))
		return myrasec.DNSRecord{}, err
	}
	c.logger.Info("deleted record", slog.Any("deleted record", deletedRecord), slog.Any("domainId", domainId))

	return *deletedRecord, err
}

// Adds the given domain in the Myra Dns Server and returns the record
func (c *MyraClient) OnAdd(record myrasec.DNSRecord) (myrasec.DNSRecord, error) {
	record.RecordType = "TXT"

	domain_id, err := c.get_domain_id(record.Name)
	if err != nil {
		return myrasec.DNSRecord{}, fmt.Errorf("A domain id is required in order to add the record; %w", err)
	}
	createdRecord, err := c.api.CreateDNSRecord(&record, domain_id)
	if err != nil {
		return myrasec.DNSRecord{}, fmt.Errorf("Could not create dns record; %w", err)
	}
	return *createdRecord, err
}

func (c *MyraClient) get_record_id(domainId int, recordName string) (int, error) {
	records, err := c.api.ListDNSRecords(domainId, nil)
	if err != nil {
		return -1, fmt.Errorf("Was not able to fetch records; %w", err)
	}

	recordId, err := extractId[myrasec.DNSRecord](records, recordName, func(r myrasec.DNSRecord) (int, string) { return r.ID, r.Name })
	if err != nil {
		return -1, fmt.Errorf("Could not derive record id from the fetched records; %s", err)
	}
	return recordId, nil
}

// Get the domain id with the provided record name
func (c *MyraClient) get_domain_id(recordName string) (int, error) {
	domains, err := c.api.ListDomains(nil)
	if err != nil {
		return -1, fmt.Errorf("Could not fetch domains; %w", err)
	}

	name, err := extractTopDomain(recordName)
	if err != nil {
		return 0, fmt.Errorf("Could not extract domain name; %w", err)
	}

	domain_id, err := extractId[myrasec.Domain](domains, name, func(r myrasec.Domain) (int, string) { return r.ID, r.Name })
	if err != nil {
		return -1, fmt.Errorf("Could not derive domain id from the fetched domains; %s", err)
	}

	return domain_id, nil
}

// extractTopDomain returns the top-level domain from a fully qualified domain name (FQDN).
// Example:
//
//	"cert-manager-dns01-tests.example.com." -> "example.com"
//	"www.sub.example.co.uk." -> "co.uk"
func extractTopDomain(fqdn string) (string, error) {
	// Remove trailing dot if present
	fqdn = strings.TrimSuffix(fqdn, ".")

	// Split into labels
	labels := strings.Split(fqdn, ".")
	if len(labels) < 2 {
		return "", fmt.Errorf("domain must have at least two labels, got: %s", fqdn)
	}

	// Take the last two labels as the top-level domain
	topDomain := strings.Join(labels[len(labels)-2:], ".")
	return topDomain, nil
}

// extractTopDomain returns the id of the domain associated with the name given.
func extractId[T any](entries []T, name string, extractor func(T) (int, string)) (int, error) {
	for _, entry := range entries {

		id, n := extractor(entry)

		domainFound := n == name
		if domainFound {
			return id, nil
		}
	}

	return 0, fmt.Errorf("Failed to extract id of entry with name '%s' - Did not find a matching domain name", name)
}
