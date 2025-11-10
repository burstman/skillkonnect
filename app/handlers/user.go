package handlers

import (
	"fmt"
	"net/http"
	"skillKonnect/app/db"
	"skillKonnect/plugins/auth"
	"strconv"

	"github.com/anthdm/superkit/kit"
	"github.com/go-chi/chi"
)

func AdminListUsers(kit *kit.Kit) error {
	var users []auth.User
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

	if err := db.Get().Model(&auth.User{}).Where("id = ?", id).Update("status", "suspended").Error; err != nil {
		return kit.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to suspend user"})
	}
	return kit.JSON(http.StatusOK, map[string]string{"message": "user suspended"})
}

// func AdminGetUsers(kit *kit.Kit) error {
// 	var users []auth.User
// 	if err := db.Get().Find(&users).Error; err != nil {
// 		return kit.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch users"})
// 	}
// 	return kit.JSON(http.StatusOK, users)
// }
