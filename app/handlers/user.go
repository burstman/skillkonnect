package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"skillKonnect/app/db"
	"skillKonnect/app/models"
	"strconv"

	"github.com/anthdm/superkit/kit"
	"github.com/go-chi/chi/v5"
)

func WebAdminListUsers(kit *kit.Kit) error {
	var users []models.User
	if err := db.Get().Find(&users).Error; err != nil {
		return kit.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch users"})
	}
	kit.Response.Header().Set("Content-Type", "application/json")
	return kit.JSON(http.StatusOK, users)
}

// @Summary List users
// @Description Get all users
// @Tags admin, users
// @Produce json
// @Success 200 {array} models.UserSwagger
// @Failure 500 {object} map[string]string
// @Router /api/admin/users [get]
func ApiAdminListUsers(kit *kit.Kit) error {
	var users []models.User
	if err := db.Get().Find(&users).Error; err != nil {
		return kit.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch users"})
	}
	kit.Response.Header().Set("Content-Type", "application/json")
	return kit.JSON(http.StatusOK, users)
}

// @Summary Suspend user
// @Description Suspend a user by ID
// @Tags admin, users
// @Param id path int true "User ID"
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/admin/users/{id}/suspend [put]
func AdminSuspendUser(kit *kit.Kit) error {
	idStr := chi.URLParam(kit.Request, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("error conv id in AdminSuspendUser: %+v", err)
	}

	kit.Response.Header().Set("Content-Type", "application/json")
	if err := db.Get().Model(&models.User{}).Where("id = ?", id).Update("suspended", true).Error; err != nil {
		return kit.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to suspend user: " + err.Error()})
	}

	return kit.JSON(http.StatusOK, map[string]string{"message": "user suspended"})
}

// @Summary Activate user
// @Description Activate a user by ID
// @Tags admin, users
// @Param id path int true "User ID"
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/admin/users/{id}/activate [put]
func AdminActivateUser(kit *kit.Kit) error {
	idStr := chi.URLParam(kit.Request, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("error conv id in AdminActivateUser: %+v", err)
	}

	kit.Response.Header().Set("Content-Type", "application/json")
	if err := db.Get().Model(&models.User{}).Where("id = ?", id).Update("suspended", false).Error; err != nil {
		return kit.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to activate user: " + err.Error()})
	}

	return kit.JSON(http.StatusOK, map[string]string{"message": "user activated"})
}

// @Summary Get user
// @Description Get a user by ID
// @Tags admin, users
// @Param id path int true "User ID"
// @Produce json
// @Success 200 {object} models.UserSwagger
// @Failure 404 {object} map[string]string
// @Router /api/admin/users/{id} [get]
func AdminGetUser(kit *kit.Kit) error {

	idStr := chi.URLParam(kit.Request, "id")
	// @Success 200 {object} models.UserSwagger
	//log.Printf("Request path: %s", kit.Request.URL.Path)

	if idStr == "" {
		return kit.JSON(http.StatusBadRequest, map[string]string{
			"error": "missing user ID",
		})
	}
	//log.Printf("userid: %s", idStr)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("error conv id in AdminSuspendUser: %+v", err)
	}

	kit.Response.Header().Set("Content-Type", "application/json")
	var user models.User
	if err := db.Get().First(&user, id).Error; err != nil {
		return kit.JSON(http.StatusNotFound, map[string]string{
			"error": "user not found",
		})
	}

	return kit.JSON(http.StatusOK, user)
}

// @Summary Update user
// @Description Update a user by ID
// @Tags admin, users
// @Param id path int true "User ID"
// @Accept json
// @Produce json
// @Param user body object true "User update payload"
// @Success 200 {object} models.UserSwagger
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/admin/users/{id} [put]
func AdminUpdateUser(kit *kit.Kit) error {
	idStr := chi.URLParam(kit.Request, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return kit.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user id"})
	}

	var input struct {
		Email     *string `json:"email,omitempty"`
		FirstName *string `json:"first_name,omitempty"`
		LastName  *string `json:"last_name,omitempty"`
		Role      *string `json:"role,omitempty"`
		Suspended *bool   `json:"suspended,omitempty"`
	}
	if err := json.NewDecoder(kit.Request.Body).Decode(&input); err != nil {
		return kit.JSON(http.StatusBadRequest, map[string]string{"error": "invalid json"})
	}

	var user models.User
	if err := db.Get().First(&user, id).Error; err != nil {
		return kit.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
	}

	// Update fields if provided
	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.FirstName != nil {
		user.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		user.LastName = *input.LastName
	}
	if input.Role != nil {
		user.Role = *input.Role
	}
	if input.Suspended != nil {
		user.Suspended = *input.Suspended
	}

	if err := db.Get().Save(&user).Error; err != nil {
		return kit.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update user"})
	}

	return kit.JSON(http.StatusOK, user)
}

// @Summary Delete user
// @Description Soft delete a user by ID
// @Tags admin, users
// @Param id path int true "User ID"
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/admin/users/{id} [delete]
func AdminDeleteUser(kit *kit.Kit) error {
	idStr := chi.URLParam(kit.Request, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return kit.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user id"})
	}

	var user models.User
	if err := db.Get().First(&user, id).Error; err != nil {
		return kit.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
	}

	if err := db.Get().Delete(&user).Error; err != nil {
		return kit.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete user"})
	}

	return kit.JSON(http.StatusOK, map[string]string{"message": "user deleted"})
}
