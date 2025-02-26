package handler

import (
	"example.com/takehome/model"
	"github.com/stretchr/testify/mock"
)

type MockPointService struct {
	mock.Mock
}

func (m *MockPointService) ProcessReceipt(receipt model.Receipt) (string, error) {
	args := m.Called(receipt)
	return args.String(0), args.Error(1)
}

func (m *MockPointService) GetPoint(id string) (int, error) {
	args := m.Called(id)
	return args.Int(0), args.Error(1)
}
