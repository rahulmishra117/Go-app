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

func CreateItem(item *models.Item) error {
	if item.Name == "" || item.Price <= 0 {
		return errors.New("invalid item data")
	}
	err := repository.CreateItem(item)
	if err == nil {
		config.RedisClient.Del(context.Background(), "all_items")
	}
	return err
}

func GetAllItems() ([]models.Item, error) {
	ctx := context.Background()
	redisKey := "all_items"

	cachedItems, err := config.RedisClient.Get(ctx, redisKey).Result()
	if err == nil {
		var items []models.Item
		json.Unmarshal([]byte(cachedItems), &items)
		fmt.Println("Cache hit for all items")
		return items, nil
	}

	var items []models.Item
	err = repository.GetAllItems(&items)
	if err != nil {
		return nil, err
	}

	itemsJSON, _ := json.Marshal(items)
	config.RedisClient.Set(ctx, redisKey, itemsJSON, 5*time.Minute)

	fmt.Println("Cache miss. Items fetched from DB and cached")
	return items, nil
}

func GetItemByID(id uuid.UUID) (*models.Item, error) {
	ctx := context.Background()
	redisKey := fmt.Sprintf("item:%s", id.String())

	cachedItem, err := config.RedisClient.Get(ctx, redisKey).Result()
	if err == nil {
		var item models.Item
		json.Unmarshal([]byte(cachedItem), &item)
		fmt.Println("Cache hit for item:", id)
		return &item, nil
	}

	var item models.Item
	err = repository.GetItemByID(id, &item)
	if err != nil {
		return nil, errors.New("item not found")
	}

	itemJSON, _ := json.Marshal(item)
	config.RedisClient.Set(ctx, redisKey, itemJSON, 5*time.Minute)

	fmt.Println("Cache miss. Item fetched from DB:", id)
	return &item, nil
}

func UpdateItem(id uuid.UUID, updatedItem *models.Item) error {
	item, err := GetItemByID(id)
	if err != nil {
		return err
	}

	item.Name = updatedItem.Name
	item.Price = updatedItem.Price

	err = repository.UpdateItem(item)
	if err == nil {

		config.RedisClient.Del(context.Background(), fmt.Sprintf("item:%s", id.String()))
		config.RedisClient.Del(context.Background(), "all_items")
	}
	return err
}

func DeleteItem(id uuid.UUID) error {
	err := repository.SoftDeleteItem(id)
	if err == nil {
		config.RedisClient.Del(context.Background(), fmt.Sprintf("item:%s", id.String()))
		config.RedisClient.Del(context.Background(), "all_items")
	}
	return err
}
