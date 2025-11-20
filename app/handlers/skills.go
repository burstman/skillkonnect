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

// @Summary List skills
// @Description Get all skills
// @Tags admin, skills
// @Produce json
// @Success 200 {array} models.SkillSwagger
// @Failure 500 {object} map[string]string
// @Router /api/v1/admin/skills [get]
func AdminListSkills(kit *kit.Kit) error {
	var skills []models.Skill

	if err := db.Get().Preload("Category").Find(&skills).Error; err != nil {
		return kit.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch skills",
		})
	}

	return kit.JSON(http.StatusOK, skills)
}

// @Summary Create skill
// @Description Create a new skill
// @Tags admin, skills
// @Accept json
// @Produce json
// @Param skill body object true "Skill payload"
// @Success 201 {object} models.SkillSwagger
// @Failure 400 {object} map[string]string
// @Router /api/v1/admin/skills [post]
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

// @Summary Delete skill
// @Description Delete a skill by ID
// @Tags admin, skills
// @Param id path int true "Skill ID"
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/v1/admin/skills/{id} [delete]
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

// @Summary Update skill
// @Description Update a skill by ID
// @Tags admin, skills
// @Param id path int true "Skill ID"
// @Accept json
// @Produce json
// @Param skill body object true "Skill update payload"
// @Success 200 {object} models.SkillSwagger
// @Failure 400 {object} map[string]string
// @Router /api/v1/admin/skills/{id} [put]
func AdminUpdateSkill(kit *kit.Kit) error {
	idStr := chi.URLParam(kit.Request, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return kit.JSON(http.StatusBadRequest, map[string]string{"error": "invalid skill id"})
	}

	var input struct {
		Name        *string `json:"name,omitempty"`
		Description *string `json:"description,omitempty"`
		CategoryID  *uint   `json:"category_id,omitempty"`
	}
	if err := json.NewDecoder(kit.Request.Body).Decode(&input); err != nil {
		return kit.JSON(http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
	}

	var skill models.Skill
	if err := db.Get().First(&skill, id).Error; err != nil {
		return kit.JSON(http.StatusNotFound, map[string]string{"error": "skill not found"})
	}

	if input.Name != nil {
		skill.Name = *input.Name
	}
	if input.Description != nil {
		skill.Description = *input.Description
	}
	if input.CategoryID != nil {
		skill.CategoryID = *input.CategoryID
	}

	if err := db.Get().Save(&skill).Error; err != nil {
		return kit.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update skill"})
	}

	return kit.JSON(http.StatusOK, skill)
}
