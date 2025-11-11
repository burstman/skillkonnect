package handlers

import (
	"fmt"
	"log"
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
	//log.Printf("Userid: %s", idStr)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("error conv id in AdminSuspendUser: %+v", err)
	}

	kit.Response.Header().Set("Content-Type", "application/json")
	if err := db.Get().Model(&auth.User{}).Where("id = ?", id).Update("status", "suspended").Error; err != nil {
		return kit.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to suspend user"})
	}

	return kit.JSON(http.StatusOK, map[string]string{"message": "user suspended"})
}

func AdminGetUser(kit *kit.Kit) error {

	ctx := chi.RouteContext(kit.Request.Context())
	if ctx == nil {
		log.Println("Route context is nil ⚠️")
	} else {
		log.Printf("Params keys: %+v", ctx.URLParams.Keys)
		log.Printf("Params values: %+v", ctx.URLParams.Values)
	}

	idStr := chi.URLParam(kit.Request, "id")
	log.Printf("Request path: %s", kit.Request.URL.Path)

	//idStr := chi.RouteContext(kit.Request.Context()).URLParam("id")

	if idStr == "" {
		return kit.JSON(http.StatusBadRequest, map[string]string{
			"error": "missing user ID",
		})
	}
	log.Printf("userid: %s", idStr)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("error conv id in AdminSuspendUser: %+v", err)
	}

	kit.Response.Header().Set("Content-Type", "application/json")
	var user auth.User
	if err := db.Get().First(&user, id).Error; err != nil {
		return kit.JSON(http.StatusNotFound, map[string]string{
			"error": "user not found",
		})
	}

	return kit.JSON(http.StatusOK, user)
}
