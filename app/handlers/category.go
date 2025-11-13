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

func AdminListCategories(kit *kit.Kit) error {
	var categories []models.Category

	if err := db.Get().Find(&categories).Error; err != nil {
		return kit.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch categories",
		})
	}

	return kit.JSON(http.StatusOK, categories)
}

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
