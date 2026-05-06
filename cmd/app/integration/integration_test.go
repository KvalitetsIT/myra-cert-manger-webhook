package main_test

import (
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/KvalitetsIT/cert-manager-webhook-core/pkg/mapping"
	"github.com/KvalitetsIT/cert-manager-webhook-core/pkg/solvers"

	"github.com/KvalitetsIT/cert-manager-webhook-core/pkg/client"
	"github.com/KvalitetsIT/cert-manager-webhook-core/pkg/configs"
	"github.com/KvalitetsIT/cert-manager-webhook-core/pkg/logging"
	"github.com/KvalitetsIT/cert-manager-webhook-core/pkg/models"

	myraClient "github.com/KvalitetsIT/cert-manager-webhook-myra/internal/client"
	"github.com/KvalitetsIT/cert-manager-webhook-myra/internal/client/adaptors"
	"github.com/KvalitetsIT/cert-manager-webhook-myra/internal/testutil"

	"github.com/miekg/dns"

	"github.com/cert-manager/cert-manager/pkg/acme/webhook"
	acmetest "github.com/cert-manager/cert-manager/test/acme"
)

var (
	zone   = os.Getenv("TEST_ZONE_NAME")
	logger *slog.Logger
)

func init() {
	logger = logging.NewJSONLogger()
	slog.SetDefault(logger)
}

// This test spins up both the mocked dns and http server
// These two servers combined will mimic the myra api
func TestRunsSuite(t *testing.T) {
	dns_port := "59351"
	myra := NewMyraApiMock("8080", dns_port, logger)
	myra.ListenAndServe()
	solver := get_solver(get_real_client())
	runTest(t, solver, dns_port)
}

func TestRunsFakeSuite(t *testing.T) {
	store := testutil.NewStorage(logger)
	dns_port := "59352"
	mocked_dns := dns.Server{
		Addr:    fmt.Sprintf(":%s", dns_port),
		Net:     "udp",
		Handler: NewDnsHandler(store, logger),
	}
	go func() {
		if err := mocked_dns.ListenAndServe(); err != nil {
			logger.Error("DNS Server failed", slog.Any("error", err))
		}
	}()
	solver := get_solver(NewMockedClient(store))
	runTest(t, solver, dns_port)
}

func runTest(t *testing.T, solver webhook.Solver, dns_port string) {
	fixture := acmetest.NewFixture(solver,
		acmetest.SetResolvedZone(zone),
		acmetest.SetAllowAmbientCredentials(false),
		acmetest.SetManifestPath("testdata/cert-manager-webhook-myra"),
		acmetest.SetDNSServer(fmt.Sprintf("localhost:%s", dns_port)),
		acmetest.SetUseAuthoritative(false), // <- do not contact system dns
		// Excluded since the method does not exist
		// acmetest.SetBinariesPath("_test/kubebuilder/bin"),
	)

	// need to uncomment and RunConformance delete runBasic and runExtended
	// once https://github.com/cert-manager/cert-manager/pull/4835 is merged
	//fixture.RunConformance(t)
	fixture.RunBasic(t)
	fixture.RunExtended(t)
}

func get_solver(c client.Client[models.Record]) webhook.Solver {
	cfg := configs.Hook{
		GroupName: "test-group",
		Port:      59351,
	}
	clientLogger := client.NewClientLogger(c, logger)
	solver := solvers.NewSolver(cfg, clientLogger, logger)
	adaptor := solvers.NewSolverAdaptor(solver, mapping.NewCertManagerMapper())
	return adaptor
}

func get_real_client() client.Client[models.Record] {
	cfg := configs.Myra{
		Api: configs.Api{
			URL:    "http://localhost:8080",
			Key:    "dummy-key",
			Secret: "dummy-secret",
			Token:  "dummy-token",
		},
	}

	if myraClient, err := myraClient.NewMyraClient(cfg, logger); err != nil {
		logger.Error("Could not create myra client", slog.Any("error", err))
		panic(err)
	} else {
		clientAdapter := adaptors.NewMyraClientAdaptor(myraClient)
		return clientAdapter
	}
}
