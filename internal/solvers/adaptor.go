package solvers

import (
	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/mapping"
	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/models"
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"k8s.io/client-go/rest"
)

// The purpose of this is to adapt from by calling the appropiate mappers
type SolverAdaptor struct {
	solver Solver[models.Record]
	mapper mapping.Mapper[models.Record, v1alpha1.ChallengeRequest]
}

func NewSolverAdaptor(solver Solver[models.Record], mapper mapping.Mapper[models.Record, v1alpha1.ChallengeRequest]) *SolverAdaptor {
	return &SolverAdaptor{
		solver: solver,
		mapper: mapper,
	}
}

func (s *SolverAdaptor) Name() string {
	return s.solver.Name()
}

func (s *SolverAdaptor) Present(ch *v1alpha1.ChallengeRequest) error {
	record := s.mapper.ToInternal(*ch)
	return s.solver.Present(record)
}

func (s *SolverAdaptor) CleanUp(ch *v1alpha1.ChallengeRequest) error {
	record := s.mapper.ToInternal(*ch)
	return s.solver.CleanUp(record)
}

func (s *SolverAdaptor) Initialize(kubeClientConfig *rest.Config, stopCh <-chan struct{}) error {
	return s.solver.Initialize(kubeClientConfig, stopCh)
}
