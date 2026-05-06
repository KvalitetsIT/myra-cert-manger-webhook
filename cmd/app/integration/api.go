package main_test

import (
	"log/slog"
	"net/http"

	"github.com/KvalitetsIT/cert-manager-webhook-myra/internal/testutil"
	"github.com/Myra-Security-GmbH/myrasec-go/v2"
	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
	"github.com/miekg/dns"
)

type myraApiMock struct {
	dns    *dns.Server
	http   *http.Server
	store  *testutil.Storage
	logger *slog.Logger
}

func NewMyraApiMock(http_port string, dns_port string, logger *slog.Logger) *myraApiMock {

	storage := testutil.NewStorage(logger)

	// Pre-populate domains (from your previous hardcoded JSON)
	initialDomains := []string{
		"google.com",
		"facebook.com",
		"platform.dk",
		"stjerne.dk",
		"example.com",
		"skyen.dk",
	}

	for _, domain := range initialDomains {
		// Add domain with empty record map

		domain := myrasec.Domain{
			Created:     &types.DateTime{},
			Modified:    &types.DateTime{},
			Name:        domain,
			AutoUpdate:  false,
			AutoDNS:     false,
			Paused:      false,
			PausedUntil: &types.DateTime{},
			Reversed:    false,
		}

		storage.AddDomain(domain)
	}

	return &myraApiMock{
		store: storage,
		http: &http.Server{
			Addr:    ":" + http_port,
			Handler: newMyraHttpHandler(storage, logger),
		},
		dns: &dns.Server{
			Addr:    ":" + dns_port,
			Net:     "udp",
			Handler: NewDnsHandler(storage, logger),
		},
		logger: logger,
	}
}

func (e *myraApiMock) ListenAndServe() {
	go func() {
		if err := e.http.ListenAndServe(); err != nil {
			e.logger.Error("HTTP Server failed", slog.Any("error", err))
		}
	}()

	go func() {
		if err := e.dns.ListenAndServe(); err != nil {
			e.logger.Error("DNS Server failed", slog.Any("error", err))
		}
	}()
}
