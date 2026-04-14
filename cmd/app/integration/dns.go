package main_test

import (
	"fmt"
	"log/slog"

	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/testutil"
	"github.com/miekg/dns"
)

type DnsHandler struct {
	store  *testutil.Storage // fqdn -> key
	logger *slog.Logger
}

func NewDnsHandler(store *testutil.Storage, logger *slog.Logger) *DnsHandler {
	return &DnsHandler{
		store:  store,
		logger: logger,
	}
}

func (e *DnsHandler) ServeDNS(w dns.ResponseWriter, req *dns.Msg) {
	e.logger.Info("ServeDNS", slog.Any("message", req))
	msg := new(dns.Msg)
	msg.SetReply(req)
	switch req.Opcode {
	case dns.OpcodeQuery:
		for _, q := range msg.Question {
			if err := e.addDNSAnswer(q, msg, req); err != nil {
				msg.SetRcode(req, dns.RcodeServerFailure)
				break
			}
		}
	}
	w.WriteMsg(msg)
}

func (e *DnsHandler) addDNSAnswer(q dns.Question, msg *dns.Msg, req *dns.Msg) error {
	switch q.Qtype {
	// Always return loopback for any A query
	case dns.TypeA:
		return e.addRecord(
			msg,
			fmt.Sprintf("%s 5 IN A 127.0.0.1", q.Name),
		)

	// TXT records are the only important record for ACME dns-01 challenges
	case dns.TypeTXT:

		var fqdm = q.Name

		domainId, found := e.store.GetDomainID(fqdm)
		if !found {
			msg.SetRcode(req, dns.RcodeNameError)
			return nil
		}

		recordId, found := e.store.GetRecordID(fqdm)
		if !found {
			msg.SetRcode(req, dns.RcodeNameError)
			return nil
		}

		record, found := e.store.GetRecord(domainId, recordId)
		if !found {
			msg.SetRcode(req, dns.RcodeNameError)
			return nil
		}
		return e.addRecord(
			msg,
			fmt.Sprintf("%s 5 IN TXT %s", q.Name, record.Value),
		)

	// NS and SOA are for authoritative lookups, return obviously invalid data
	case dns.TypeNS:
		return e.addRecord(
			msg,
			fmt.Sprintf("%s 5 IN NS ns.example-acme-webook.invalid.", q.Name),
		)
	case dns.TypeSOA:
		return e.addRecord(
			msg,
			fmt.Sprintf("%s 5 IN SOA %s 20 5 5 5 5", "ns.example-acme-webook.invalid.", "ns.example-acme-webook.invalid."),
		)
	default:
		return fmt.Errorf("unimplemented record type %v", q.Qtype)
	}
}

func (e *DnsHandler) addRecord(msg *dns.Msg, record string) error {
	rr, err := dns.NewRR(record)
	if err != nil {
		return err
	}
	msg.Answer = append(msg.Answer, rr)
	return nil
}
