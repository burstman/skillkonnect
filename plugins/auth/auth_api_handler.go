package auth

import (
	"encoding/json"
	"net/http"
	"skillKonnect/app/db"
	"skillKonnect/app/models"
	"strings"

	"time"

	"github.com/anthdm/superkit/kit"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HandleApiLoginCreate(kit *kit.Kit) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(kit.Request.Body).Decode(&req); err != nil {
		return kit.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid json",
		})
	}

	// Find user
	var user models.User
	if err := db.Get().Where("email = ?", req.Email).First(&user).Error; err != nil {
		return kit.JSON(http.StatusUnauthorized, map[string]string{
			"error": "invalid credentials",
		})
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return kit.JSON(http.StatusUnauthorized, map[string]string{
			"error": "invalid credentials",
		})
	}

	// Only admin can use API
	if user.Role != "admin" {
		return kit.JSON(http.StatusForbidden, map[string]string{
			"error": "admin only",
		})
	}

	// Create token session
	token := uuid.New().String()

	session := models.Session{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}
	db.Get().Create(&session)

	return kit.JSON(http.StatusOK, map[string]any{
		"token":      token,
		"expires_at": session.ExpiresAt,
		"user": map[string]any{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func HandleApiLoginDelete(kit *kit.Kit) error {
	token := kit.Request.Header.Get("Authorization")
	if token == "" || !strings.HasPrefix(token, "Bearer ") {
		return kit.JSON(http.StatusUnauthorized, map[string]string{
			"error": "missing token",
		})
	}

	rawToken := strings.TrimPrefix(token, "Bearer ")

	if err := db.Get().Where("token = ?", rawToken).Delete(&models.Session{}).Error; err != nil {
		return kit.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to logout",
		})
	}

	return kit.JSON(http.StatusOK, map[string]string{
		"message": "logged out",
	})
}

// HandleApiAuthMe returns information about the currently authenticated user.
func HandleApiAuthMe(kit *kit.Kit) error {
	// Read authentication from unified middleware
	payload, ok := kit.Request.Context().Value(AuthContextKey{}).(AuthPayload)
	if !ok || !payload.Authenticated || payload.User == nil {
		return kit.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	user := payload.User

	return kit.JSON(http.StatusOK, map[string]any{
		"id":         user.ID,
		"email":      user.Email,
		"role":       user.Role,
		"first_name": user.FirstName,
	})
}
