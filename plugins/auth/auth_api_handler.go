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

// @Summary Health check
// @Description Check if the server is running
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/v1/health [get]
func HandleHealthCheck(kit *kit.Kit) error {
	return kit.JSON(http.StatusOK, map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// @Summary Admin login
// @Description Login with email and password; returns token and user information
// @Tags auth, public
// @Accept json
// @Produce json
// @Param credentials body object true "Login payload"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/admin/login [post]
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

// @Summary Client login
// @Description Login as client with email and password; returns token and user information
// @Tags auth, public
// @Accept json
// @Produce json
// @Param credentials body object true "Login payload with email and password"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/client/login [post]
func HandleClientLoginCreate(kit *kit.Kit) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(kit.Request.Body).Decode(&req); err != nil {
		return kit.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid json",
		})
	}

	// Trim spaces
	req.Email = strings.TrimSpace(req.Email)

	// Find user by email
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

	// Only clients can use this endpoint
	if user.Role != "client" {
		return kit.JSON(http.StatusForbidden, map[string]string{
			"error": "client only",
		})
	}

	// Check if suspended
	if user.Suspended {
		return kit.JSON(http.StatusForbidden, map[string]string{
			"error": "account suspended",
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
			"id":         user.ID,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"role":       user.Role,
		},
	})
}

// @Summary Logout
// @Description Invalidate the current API token (Bearer)
// @Tags auth
// @Param Authorization header string true "Bearer token"
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/admin/logout [delete]
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

// @Summary Current user
// @Description Returns information about the authenticated user
// @Tags auth
// @Param Authorization header string true "Bearer token"
// @Produce json
// @Success 200 {object} models.UserSwagger
// @Failure 401 {object} map[string]string
// @Router /api/v1/auth/me [get]
// HandleApiAuthMe returns information about the currently authenticated user.
func HandleApiAuthMe(kit *kit.Kit) error {
	// Read authentication from unified middleware
	payload, ok := kit.Request.Context().Value(models.AuthContextKey{}).(models.AuthPayload)
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
