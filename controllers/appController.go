package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rahulmishra/go-crud-app/models"
	"github.com/rahulmishra/go-crud-app/services"
)

// CreateItem godoc
// @Summary Create a new item
// @Description Adds a new item to the database
// @Tags Items
// @Accept json
// @Produce json
// @Param item body models.Item true "Item Data"
// @Success 201 {object} models.Item
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /items/ [post]
func CreateItem(c *gin.Context) {
	var item models.Item
	item.ID = uuid.New()
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := services.CreateItem(&item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// GetAllItems godoc
// @Summary Get all items
// @Description Retrieves a list of all items
// @Tags Items
// @Produce json
// @Success 200 {array} models.Item
// @Failure 500 {object} map[string]string
// @Router /items/ [get]
func GetAllItems(c *gin.Context) {
	items, err := services.GetAllItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetItemByID godoc
// @Summary Get an item by ID
// @Description Retrieves an item by its ID
// @Tags Items
// @Produce json
// @Param id path string true "Item ID (UUID)"
// @Success 200 {object} models.Item
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /items/{id} [get]
func GetItemByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	item, err := services.GetItemByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// UpdateItem godoc
// @Summary Update an item
// @Description Updates an existing item
// @Tags Items
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Param item body models.Item true "Updated Item Data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /items/{id} [put]
func UpdateItem(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	var updatedItem models.Item
	if err := c.ShouldBindJSON(&updatedItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedItem.ID = id
	if err := services.UpdateItem(id, &updatedItem); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Item updated successfully"})
}

// DeleteItem godoc
// @Summary Delete an item
// @Description Deletes an item by ID
// @Tags Items
// @Produce json
// @Param id path string true "Item ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /items/{id} [delete]
func DeleteItem(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	if err := services.DeleteItem(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}
