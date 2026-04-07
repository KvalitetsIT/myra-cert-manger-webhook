package solvers

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/client"
	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/configs"
	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/models"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)



type InternalSolver struct {
	cfg              configs.Hook
	kubernetesClient *kubernetes.Clientset
	client           client.Client[models.Record]
	logger           *slog.Logger
	sync.RWMutex
}

func NewSolver(cfg configs.Hook, publisher client.Client[models.Record], logger *slog.Logger) *InternalSolver {
	solver := &InternalSolver{
		cfg:              cfg,
		kubernetesClient: &kubernetes.Clientset{},
		client:           publisher,
		logger:           logger,
	}
	return solver
}

func (s *InternalSolver) Name() string {
	return s.cfg.GroupName
}

// The 'Present' is triggered by the cert-manager when a ACME challenge is initialized
func (s *InternalSolver) Present(record models.Record) error {
	// Create the record for the DNS provider
	_, err := s.client.OnAdd(record)
	if err != nil {
		return createError("Present", err)
	}
	return nil
}

// A clean up event is triggered when the record created in the DNS has been successfully validated by the Certificate Authority (CA)
func (s *InternalSolver) CleanUp(record models.Record) error {
	// Delete a record from the DNS provider's console
	_, err := s.client.OnDelete(record)
	if err != nil {
		return createError("CleanUp", err)
	}
	return nil
}

func (s *InternalSolver) Initialize(kubeClientConfig *rest.Config, stopCh <-chan struct{}) error {
	return nil
}

func createError(event string, err error) error {
	return fmt.Errorf("Failed to handle the '%s' event; %w", event, err)
}
