package repository

import (
	"github.com/google/uuid"
	"github.com/rahulmishra/go-crud-app/models"
	"github.com/stretchr/testify/mock"
)

// MockAppRepository is a mock for repository functions
type MockAppRepository struct {
	mock.Mock
}

func (m *MockAppRepository) CreateItem(item *models.Item) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MockAppRepository) GetAllItems(items *[]models.Item) error {
	args := m.Called(items)
	return args.Error(0)
}

func (m *MockAppRepository) GetItemByID(id uuid.UUID, item *models.Item) error {
	args := m.Called(id, item)
	return args.Error(0)
}

func (m *MockAppRepository) UpdateItem(item *models.Item) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MockAppRepository) SoftDeleteItem(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}
