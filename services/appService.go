package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rahulmishra/go-crud-app/config"
	"github.com/rahulmishra/go-crud-app/models"
	"github.com/rahulmishra/go-crud-app/repository"
)

// CreateItem creates a new item and clears cache
func CreateItem(item *models.Item) error {
	if item.Name == "" || item.Price <= 0 {
		return errors.New("invalid item data")
	}
	err := repository.CreateItem(item)
	if err == nil {
		// Invalidate cache after creating a new item
		config.RedisClient.Del(context.Background(), "all_items")
	}
	return err
}

// GetAllItems retrieves all items with Redis caching
func GetAllItems() ([]models.Item, error) {
	ctx := context.Background()
	redisKey := "all_items"

	// Check Redis cache first
	cachedItems, err := config.RedisClient.Get(ctx, redisKey).Result()
	if err == nil {
		var items []models.Item
		json.Unmarshal([]byte(cachedItems), &items)
		fmt.Println("Cache hit for all items")
		return items, nil
	}

	// Fetch from DB if cache is empty
	var items []models.Item
	err = repository.GetAllItems(&items)
	if err != nil {
		return nil, err
	}

	// Cache the result in Redis for 5 minutes
	itemsJSON, _ := json.Marshal(items)
	config.RedisClient.Set(ctx, redisKey, itemsJSON, 5*time.Minute)

	fmt.Println("Cache miss. Items fetched from DB and cached")
	return items, nil
}

// GetItemByID retrieves an item by ID with Redis caching
func GetItemByID(id uuid.UUID) (*models.Item, error) {
	ctx := context.Background()
	redisKey := fmt.Sprintf("item:%s", id.String())

	// Try to get item from Redis cache
	cachedItem, err := config.RedisClient.Get(ctx, redisKey).Result()
	if err == nil {
		var item models.Item
		json.Unmarshal([]byte(cachedItem), &item)
		fmt.Println("Cache hit for item:", id)
		return &item, nil
	}

	// Fetch from DB if not found in cache
	var item models.Item
	err = repository.GetItemByID(id, &item)
	if err != nil {
		return nil, errors.New("item not found")
	}

	// Cache the item in Redis
	itemJSON, _ := json.Marshal(item)
	config.RedisClient.Set(ctx, redisKey, itemJSON, 5*time.Minute)

	fmt.Println("Cache miss. Item fetched from DB:", id)
	return &item, nil
}

// UpdateItem updates an existing item and clears cache
func UpdateItem(id uuid.UUID, updatedItem *models.Item) error {
	item, err := GetItemByID(id)
	if err != nil {
		return err
	}

	// Update fields
	item.Name = updatedItem.Name
	item.Price = updatedItem.Price

	// Update in DB
	err = repository.UpdateItem(item)
	if err == nil {
		// Invalidate cache
		config.RedisClient.Del(context.Background(), fmt.Sprintf("item:%s", id.String()))
		config.RedisClient.Del(context.Background(), "all_items")
	}
	return err
}

// DeleteItem performs a soft delete and clears cache
func DeleteItem(id uuid.UUID) error {
	err := repository.SoftDeleteItem(id)
	if err == nil {
		// Remove item from cache
		config.RedisClient.Del(context.Background(), fmt.Sprintf("item:%s", id.String()))
		config.RedisClient.Del(context.Background(), "all_items")
	}
	return err
}
