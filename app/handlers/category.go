package handlers

import (
	"encoding/json"
	"net/http"
	"skillKonnect/app/db"
	"skillKonnect/app/models"
	"strconv"
	"strings"

	"github.com/anthdm/superkit/kit"
	"github.com/go-chi/chi/v5"
)

// @Summary List categories
// @Description Get all categories
// @Tags admin, categories
// @Produce json
// @Success 200 {array} models.CategorySwagger
// @Failure 500 {object} map[string]string
// @Router /api/v1/admin/categories [get]
func AdminListCategories(kit *kit.Kit) error {
	var categories []models.Category

	if err := db.Get().Find(&categories).Error; err != nil {
		return kit.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch categories",
		})
	}

	return kit.JSON(http.StatusOK, categories)
}

// @Summary Create category
// @Description Create a new category
// @Tags admin, categories
// @Accept json
// @Produce json
// @Param category body object true "Category payload"
// @Success 201 {object} models.CategorySwagger
// @Failure 400 {object} map[string]string
// @Router /api/v1/admin/categories [post]
func AdminCreateCategory(kit *kit.Kit) error {
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(kit.Request.Body).Decode(&input); err != nil {
		return kit.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid JSON",
		})
	}

	if input.Name == "" {
		return kit.JSON(http.StatusBadRequest, map[string]string{
			"error": "name is required",
		})
	}

	category := models.Category{
		Name:        strings.ToLower(input.Name),
		Description: input.Description,
	}

	if err := db.Get().Create(&category).Error; err != nil {
		return kit.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to create category: " + err.Error(),
		})
	}

	return kit.JSON(http.StatusCreated, category)
}

// @Summary Delete category
// @Description Delete a category by ID
// @Tags admin, categories
// @Param id path int true "Category ID"
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/v1/admin/categories/{id} [delete]
func AdminDeleteCategory(kit *kit.Kit) error {
	idStr := chi.URLParam(kit.Request, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return kit.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid category ID",
		})
	}

	if err := db.Get().Delete(&models.Category{}, id).Error; err != nil {
		return kit.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to delete category",
		})
	}

	return kit.JSON(http.StatusOK, map[string]string{
		"message": "category deleted successfully",
	})
}

// @Summary Update category
// @Description Update a category by ID
// @Tags admin, categories
// @Param id path int true "Category ID"
// @Accept json
// @Produce json
// @Param category body object true "Category update payload"
// @Success 200 {object} models.CategorySwagger
// @Failure 400 {object} map[string]string
// @Router /api/v1/admin/categories/{id} [put]
func AdminUpdateCategory(kit *kit.Kit) error {
	idStr := chi.URLParam(kit.Request, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return kit.JSON(http.StatusBadRequest, map[string]string{"error": "invalid category id"})
	}

	var input struct {
		Name        *string `json:"name,omitempty"`
		Description *string `json:"description,omitempty"`
	}
	if err := json.NewDecoder(kit.Request.Body).Decode(&input); err != nil {
		return kit.JSON(http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
	}

	var category models.Category
	if err := db.Get().First(&category, id).Error; err != nil {
		return kit.JSON(http.StatusNotFound, map[string]string{"error": "category not found"})
	}

	if input.Name != nil {
		category.Name = *input.Name
	}
	if input.Description != nil {
		category.Description = *input.Description
	}

	if err := db.Get().Save(&category).Error; err != nil {
		return kit.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update category"})
	}

	return kit.JSON(http.StatusOK, category)
}
