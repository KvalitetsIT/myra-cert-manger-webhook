package service

import (
	"log/slog"

	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/client"
	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/client/adaptors"
	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/configs"
	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/mapping"
	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/solvers"
)

type ServiceFactory struct {
	cfg    configs.Configuration
	logger *slog.Logger
}

func NewServiceFactory(cfg configs.Configuration, logger *slog.Logger) *ServiceFactory {
	return &ServiceFactory{
		cfg:    cfg,
		logger: logger,
	}
}

func (s ServiceFactory) CreateDefault() (*service, error) {
	if myra, err := client.NewMyraClient(s.cfg.Myra); err != nil {
		return nil, err
	} else {
		clientAdapter := adaptors.NewMyraClientAdaptor(myra)
		clientLogger := client.NewClientLogger(clientAdapter, s.logger)
		solver := solvers.NewSolver(s.cfg.Hook, clientLogger, s.logger)
		solverAdaptor := solvers.NewSolverAdaptor(solver, mapping.NewCertManagerMapper())
		return newService(s.cfg.Hook, solverAdaptor, s.logger), nil
	}

}
