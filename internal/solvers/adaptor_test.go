package solvers_test

import (
	"testing"

	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/models"
	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/solvers"
	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/testutil/mocks"
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/rest"
)

func TestSolverAdaptor_Name(t *testing.T) {
	mockedSolver := new(mocks.MockedSolver)
	mockedMapper := new(mocks.MockedMapper[models.Record, v1alpha1.ChallengeRequest])

	mockedSolver.On("Name").Return("test-solver")

	adaptor := solvers.NewSolverAdaptor(
		mockedSolver,
		mockedMapper,
	)

	name := adaptor.Name()

	require.Equal(t, "test-solver", name)
	mockedSolver.AssertExpectations(t)
}

func TestSolverAdaptor_Present(t *testing.T) {
	mockedSolver := new(mocks.MockedSolver)
	mockedMapper := new(mocks.MockedMapper[models.Record, v1alpha1.ChallengeRequest])

	ch := v1alpha1.ChallengeRequest{}
	record := models.Record{}

	mockedMapper.On("ToInternal", ch).Return(record)
	mockedSolver.On("Present", record).Return(nil)

	adaptor := solvers.NewSolverAdaptor(
		mockedSolver,
		mockedMapper,
	)

	err := adaptor.Present(&ch)

	require.NoError(t, err)
	mockedMapper.AssertExpectations(t)
	mockedSolver.AssertExpectations(t)
}

func TestSolverAdaptor_Present_Error(t *testing.T) {
	mockedSolver := new(mocks.MockedSolver)
	mockedMapper := new(mocks.MockedMapper[models.Record, v1alpha1.ChallengeRequest])

	ch := v1alpha1.ChallengeRequest{}
	record := models.Record{}

	mockedMapper.On("ToInternal", ch).Return(record)
	mockedSolver.On("Present", record).Return(assertErr())

	adaptor := solvers.NewSolverAdaptor(
		mockedSolver,
		mockedMapper,
	)

	err := adaptor.Present(&ch)

	require.Error(t, err)
}

func TestSolverAdaptor_CleanUp(t *testing.T) {
	mockSolver := new(mocks.MockedSolver)
	mockMapper := new(mocks.MockedMapper[models.Record, v1alpha1.ChallengeRequest])

	ch := v1alpha1.ChallengeRequest{}
	record := models.Record{}

	mockMapper.On("ToInternal", ch).Return(record)
	mockSolver.On("CleanUp", record).Return(nil)

	adaptor := solvers.NewSolverAdaptor(
		mockSolver,
		mockMapper,
	)

	err := adaptor.CleanUp(&ch)

	require.NoError(t, err)
}

func TestSolverAdaptor_Initialize(t *testing.T) {
	mockSolver := new(mocks.MockedSolver)
	mockMapper := new(mocks.MockedMapper[models.Record, v1alpha1.ChallengeRequest])

	cfg := &rest.Config{}

	stopCh := make(<-chan struct{})
	mockSolver.On("Initialize", cfg, stopCh).Return(nil)

	adaptor := solvers.NewSolverAdaptor(
		mockSolver,
		mockMapper,
	)

	err := adaptor.Initialize(cfg, stopCh)

	require.NoError(t, err)
	mockSolver.AssertExpectations(t)
}

func assertErr() error {
	return &testError{}
}

type testError struct{}

func (e *testError) Error() string { return "test error" }
