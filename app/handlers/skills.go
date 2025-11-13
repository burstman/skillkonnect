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

func AdminListSkills(kit *kit.Kit) error {
	var skills []models.Skill

	if err := db.Get().Preload("Category").Find(&skills).Error; err != nil {
		return kit.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch skills",
		})
	}

	return kit.JSON(http.StatusOK, skills)
}

func AdminCreateSkill(kit *kit.Kit) error {
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		CategoryID  uint   `json:"category_id"`
	}

	if err := json.NewDecoder(kit.Request.Body).Decode(&input); err != nil {
		return kit.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid JSON body",
		})
	}

	if input.Name == "" {
		return kit.JSON(http.StatusBadRequest, map[string]string{
			"error": "name is required",
		})
	}

	skill := models.Skill{
		Name:        strings.ToLower(input.Name),
		Description: input.Description,
		CategoryID:  input.CategoryID,
	}

	if err := db.Get().Create(&skill).Error; err != nil {
		return kit.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to create skill: " + err.Error(),
		})
	}

	return kit.JSON(http.StatusCreated, skill)
}

func AdminDeleteSkill(kit *kit.Kit) error {
	idStr := chi.URLParam(kit.Request, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return kit.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid skill ID",
		})
	}

	if err := db.Get().Delete(&models.Skill{}, id).Error; err != nil {
		return kit.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to delete skill",
		})
	}

	return kit.JSON(http.StatusOK, map[string]string{
		"message": "skill deleted successfully",
	})
}
