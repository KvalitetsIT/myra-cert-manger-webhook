package mocks

import (
	"github.com/stretchr/testify/mock"
)

type MockedMapper[Internal, External any] struct {
	mock.Mock
}

func (m *MockedMapper[Internal, External]) ToExternal(i Internal) External {
	args := m.Called(i)
	return args.Get(0).(External)
}

func (m *MockedMapper[Internal, External]) ToInternal(e External) Internal {
	args := m.Called(e)
	return args.Get(0).(Internal)
}
