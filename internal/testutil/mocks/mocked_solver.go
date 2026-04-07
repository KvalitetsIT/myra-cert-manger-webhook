package mocks

import (
	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/models"
	"github.com/stretchr/testify/mock"
	"k8s.io/client-go/rest"
)

type MockedSolver struct {
	mock.Mock
}

func (m *MockedSolver) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockedSolver) Present(r models.Record) error {
	args := m.Called(r)
	return args.Error(0)
}

func (m *MockedSolver) CleanUp(r models.Record) error {
	args := m.Called(r)
	return args.Error(0)
}

func (m *MockedSolver) Initialize(cfg *rest.Config, stopCh <-chan struct{}) error {
	args := m.Called(cfg, stopCh)
	return args.Error(0)
}
