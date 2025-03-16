package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Item struct {
	ID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name  string    `json:"name"`
	Price float64   `json:"price"`
}

func (item *Item) BeforeCreate(tx *gorm.DB) (err error) {
	if item.ID == (uuid.UUID{}) {
		item.ID = uuid.New()
	}
	return
}
