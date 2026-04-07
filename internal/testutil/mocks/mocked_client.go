package mocks

import "github.com/stretchr/testify/mock"

type MockedClient[T any] struct {
	mock.Mock
}

func (m *MockedClient[T]) OnAdd(record T) (T, error) {
	args := m.Called(record)
	return args.Get(0).(T), args.Error(1)
}

func (m *MockedClient[T]) OnDelete(record T) (T, error) {
	args := m.Called(record)
	return args.Get(0).(T), args.Error(1)
}
