package handlers

import (
	"fmt"
	"net/http"
	"skillKonnect/app/db"
	"skillKonnect/app/models"
	"strconv"

	"github.com/anthdm/superkit/kit"
	"github.com/go-chi/chi/v5"
)

func AdminListUsers(kit *kit.Kit) error {
	var users []models.User
	if err := db.Get().Find(&users).Error; err != nil {
		return kit.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch users"})
	}
	kit.Response.Header().Set("Content-Type", "application/json")
	return kit.JSON(http.StatusOK, users)
}

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

func AdminGetUser(kit *kit.Kit) error {

	idStr := chi.URLParam(kit.Request, "id")
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
