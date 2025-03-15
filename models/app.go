package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Item struct {
	ID    *uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id,omitempty"`
	Name  string     `json:"name"`
	Price float64    `json:"price"`
}

func (item *Item) BeforeCreate(tx *gorm.DB) (err error) {
	if item.ID == nil {
		newUUID := uuid.New()
		item.ID = &newUUID
	}
	return
}
