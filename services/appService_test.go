package services

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rahulmishra/go-crud-app/config"
	"github.com/rahulmishra/go-crud-app/models"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type MockRedisClient struct {
	mock.Mock
	redis.Cmdable
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	cmd := redis.NewStringCmd(ctx)
	cmd.SetVal(args.String(0))
	cmd.SetErr(args.Error(1))
	return cmd
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	m.Called(ctx, key, value, expiration)
	return redis.NewStatusCmd(ctx)
}

func (m *MockRedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	args := []interface{}{ctx}
	for _, key := range keys {
		args = append(args, key)
	}
	m.Called(args...)
	return redis.NewIntCmd(ctx)
}

// Mock App Repository
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

func TestCreateItem(t *testing.T) {
	mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	config.SetDB(mockDB)

	err = mockDB.AutoMigrate(&models.Item{})
	if err != nil {
		t.Fatalf("Failed to migrate database schema: %v", err)
	}

	mockItem := &models.Item{ID: uuid.New(), Name: "Test Item", Price: 10.0}

	err = config.DB.Create(mockItem).Error
	assert.NoError(t, err)

	var retrievedItem models.Item
	err = config.DB.First(&retrievedItem, "id = ?", mockItem.ID).Error
	assert.NoError(t, err)

	assert.Equal(t, mockItem.ID, retrievedItem.ID)
	assert.Equal(t, "Test Item", retrievedItem.Name)
	assert.Equal(t, 10.0, retrievedItem.Price)
}

func TestGetAllItems_CacheHit(t *testing.T) {
	mockRedis := new(MockRedisClient)
	config.RedisClient = mockRedis

	mockUUID := uuid.Must(uuid.NewRandom())
	mockItems := []models.Item{{ID: mockUUID, Name: "Item1", Price: 20}}

	cachedData, _ := json.Marshal(mockItems)
	mockRedis.On("Get", mock.Anything, "all_items").Return(string(cachedData), nil)

	items, err := GetAllItems()

	assert.NoError(t, err)
	assert.Equal(t, len(items), 1)
	assert.Equal(t, items[0].Name, "Item1")

	mockRedis.AssertExpectations(t)
}

func TestGetItemByID_CacheHit(t *testing.T) {
	mockRedis := new(MockRedisClient)
	itemID := uuid.Must(uuid.NewRandom())
	mockItem := models.Item{ID: itemID, Name: "Item1", Price: 20}

	cachedData, _ := json.Marshal(mockItem)
	mockRedis.On("Get", mock.Anything, "item:"+itemID.String()).Return(string(cachedData), nil)

	config.RedisClient = mockRedis

	item, err := GetItemByID(itemID)

	assert.NoError(t, err)
	assert.Equal(t, item.Name, "Item1")

	mockRedis.AssertExpectations(t)
}

func TestUpdateItem(t *testing.T) {
	mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	config.SetDB(mockDB)

	err = mockDB.AutoMigrate(&models.Item{})
	if err != nil {
		t.Fatalf("Failed to migrate database schema: %v", err)
	}

	mockRedis := new(MockRedisClient)
	itemID := uuid.New()

	mockItem := &models.Item{ID: itemID, Name: "Updated Item1", Price: 20}
	mockDB.Create(mockItem)

	updatedItem := &models.Item{Name: "Updated Item1", Price: 30}

	mockRedis.On("Get", mock.Anything, fmt.Sprintf("item:%s", itemID.String())).Return("cached_data", nil)
	mockRedis.On("Del", mock.Anything, fmt.Sprintf("item:%s", itemID.String())).Return(redis.NewIntCmd(context.Background()))
	mockRedis.On("Del", mock.Anything, "all_items").Return(redis.NewIntCmd(context.Background()))

	config.RedisClient = mockRedis

	// Debug before update
	var beforeUpdate models.Item
	_ = mockDB.First(&beforeUpdate, "id = ?", itemID).Error
	fmt.Println("Before Update (Test):", beforeUpdate.Name, beforeUpdate.Price)

	err = UpdateItem(itemID, updatedItem)
	assert.NoError(t, err)

	// Debug after update
	var retrievedItem models.Item
	err = mockDB.First(&retrievedItem, "id = ?", itemID).Error
	fmt.Println("After Update (Test):", retrievedItem.Name, retrievedItem.Price)

	assert.NoError(t, err)
	assert.Equal(t, "Updated Item1", retrievedItem.Name)

	mockRedis.AssertExpectations(t)
}

func TestDeleteItem(t *testing.T) {
	mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}

	config.DB = mockDB

	err = mockDB.AutoMigrate(&models.Item{})
	if err != nil {
		t.Fatalf("Failed to migrate database schema: %v", err)
	}

	itemID := uuid.New()
	mockItem := &models.Item{ID: itemID, Name: "Test Item", Price: 50}
	mockDB.Create(mockItem)

	mockRedis := new(MockRedisClient)
	config.RedisClient = mockRedis

	mockRedis.On("Del", mock.Anything, fmt.Sprintf("item:%s", itemID.String())).Return(nil)
	mockRedis.On("Del", mock.Anything, "all_items").Return(nil)

	err = DeleteItem(itemID)
	assert.NoError(t, err, "DeleteItem should not return an error")

	var retrievedItem models.Item
	err = mockDB.First(&retrievedItem, "id = ?", itemID).Error
	assert.Error(t, err, "Item should not be found after soft delete")

	mockRedis.AssertExpectations(t)
}
