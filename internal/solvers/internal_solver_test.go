package solvers_test

import (
	"errors"
	"io"
	"testing"

	"log/slog"

	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/configs"
	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/models"
	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/solvers"
	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/testutil/mocks"
	"github.com/stretchr/testify/assert"
)

func TestInternalSolver_Name(t *testing.T) {
	cfg := configs.Hook{GroupName: "test-group"}
	mockClient := &mocks.MockedClient[models.Record]{}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	solver := solvers.NewSolver(cfg, mockClient, logger)
	assert.Equal(t, "test-group", solver.Name(), "Name() should return cfg.GroupName")
}

func TestInternalSolver_Present(t *testing.T) {
	cfg := configs.Hook{GroupName: "test-group"}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	record := models.Record{
		DNSName:      "example.com",
		ResolvedZone: "example.com.",
		ResolvedFQDN: "_acme-challenge.example.com.",
		Key:          "dummy-key",
		Type:         "dns-01",
	}

	t.Run("success", func(t *testing.T) {
		// Tell testify/mock to expect OnAdd with this record and return it + nil error
		mockClient := &mocks.MockedClient[models.Record]{}
		mockClient.On("OnAdd", record).Return(record, nil)
		solver := solvers.NewSolver(cfg, mockClient, logger)

		err := solver.Present(record)
		assert.NoError(t, err, "Present should not fail on successful OnAdd")

		mockClient.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		clientErr := errors.New("client failed")
		mockClient := &mocks.MockedClient[models.Record]{}
		mockClient.On("OnAdd", record).Return(record, clientErr)

		solver := solvers.NewSolver(cfg, mockClient, logger)
		err := solver.Present(record)
		assert.Error(t, err, "Present should fail when OnAdd returns error")
		assert.Contains(t, err.Error(), "Failed to handle the 'Present' event")
		assert.Contains(t, err.Error(), "client failed")

		mockClient.AssertExpectations(t)
	})
}
func TestInternalSolver_CleanUp(t *testing.T) {
	cfg := configs.Hook{GroupName: "test-group"}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	record := models.Record{
		DNSName:      "example.com",
		ResolvedZone: "example.com.",
		ResolvedFQDN: "_acme-challenge.example.com.",
		Key:          "dummy-key",
		Type:         "dns-01",
	}

	t.Run("success", func(t *testing.T) {
		mockClient := &mocks.MockedClient[models.Record]{}
		mockClient.On("OnDelete", record).Return(record, nil)

		solver := solvers.NewSolver(cfg, mockClient, logger)

		err := solver.CleanUp(record)
		assert.NoError(t, err)
		mockClient.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockClient := &mocks.MockedClient[models.Record]{}
		clientErr := errors.New("delete failed")
		mockClient.On("OnDelete", record).Return(record, clientErr)

		solver := solvers.NewSolver(cfg, mockClient, logger)

		err := solver.CleanUp(record)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Failed to handle the 'CleanUp' event")
		assert.Contains(t, err.Error(), "delete failed")
		mockClient.AssertExpectations(t)
	})
}

func TestInternalSolver_Initialize(t *testing.T) {
	cfg := configs.Hook{GroupName: "test-group"}
	mockClient := &mocks.MockedClient[models.Record]{}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	solver := solvers.NewSolver(cfg, mockClient, logger)
	stopCh := make(chan struct{})

	err := solver.Initialize(nil, stopCh)
	assert.NoError(t, err, "Initialize should return no error")
}
