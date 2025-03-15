package repository

import (
	"github.com/google/uuid"
	"github.com/rahulmishra/go-crud-app/config"
	"github.com/rahulmishra/go-crud-app/models"
)

// CreateItem - Adds a new item to the database
func CreateItem(item *models.Item) error {
	result := config.DB.Create(item)
	return result.Error
}

// GetAllItems - Retrieves all items
func GetAllItems(items *[]models.Item) error {
	result := config.DB.Find(items)
	return result.Error
}

// GetItemByID - Fetches a single item by ID
func GetItemByID(id uuid.UUID, item *models.Item) error {
	return config.DB.Where("id = ?", id).First(item).Error
}

// UpdateItem - Updates an existing item
func UpdateItem(item *models.Item) error {
	result := config.DB.Save(item)
	return result.Error
}

// DeleteItem - Deletes an item by ID
func DeleteItem(id uint) error {
	result := config.DB.Delete(&models.Item{}, id)
	return result.Error
}
func SoftDeleteItem(id uuid.UUID) error {
	return config.DB.Where("id = ?", id).Delete(&models.Item{}).Error
}
