package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/rahulmishra/go-crud-app/config"
	"github.com/rahulmishra/go-crud-app/models"
)

func CreateItem(item *models.Item) error {
	result := config.DB.Create(item)
	return result.Error
}

func GetAllItems(items *[]models.Item) error {
	result := config.DB.Find(items)
	return result.Error
}

func GetItemByID(id uuid.UUID, item *models.Item) error {
	return config.DB.Where("id = ?", id).First(item).Error
}

func UpdateItem(item *models.Item) error {
	result := config.DB.Save(item)
	return result.Error
}

func DeleteItem(id uint) error {
	result := config.DB.Delete(&models.Item{}, id)
	return result.Error
}
func SoftDeleteItem(id uuid.UUID) error {
	if config.DB == nil {
		return errors.New("database is not initialized")
	}
	return config.DB.Where("id = ?", id).Delete(&models.Item{}).Error
}
