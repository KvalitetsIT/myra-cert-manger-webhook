package service

import (
	"log/slog"

	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/configs"
	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/solvers"
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/cmd"
)

type service struct {
	solver *solvers.SolverAdaptor
	cfg    configs.Hook
	logger *slog.Logger
}

func newService(cfg configs.Hook, listener *solvers.SolverAdaptor, logger *slog.Logger) *service {
	return &service{
		listener,
		cfg,
		logger,
	}
}

func (self *service) Start() {
	cmd.RunWebhookServer(self.cfg.GroupName, self.solver)
}
