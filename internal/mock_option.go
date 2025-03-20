package internal

import "github.com/stretchr/testify/mock"

type MockOption[T any] struct {
	mock.Mock
}

func NewMockOption[T any]() *MockOption[T] {
	return &MockOption[T]{
		Mock: mock.Mock{},
	}
}

func (mO *MockOption[T]) Option(value T) error {
	ret := mO.Called(value)
	return ret.Error(0)
}
