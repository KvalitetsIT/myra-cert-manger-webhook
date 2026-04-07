package client

import myrasec "github.com/Myra-Security-GmbH/myrasec-go/v2"

type MyraAPI interface {
	ListDomains(map[string]string) ([]myrasec.Domain, error)
	ListDNSRecords(int, map[string]string) ([]myrasec.DNSRecord, error)
	CreateDNSRecord(record *myrasec.DNSRecord, domainId int) (*myrasec.DNSRecord, error)
	DeleteDNSRecord(record *myrasec.DNSRecord, domainId int) (*myrasec.DNSRecord, error)
}
