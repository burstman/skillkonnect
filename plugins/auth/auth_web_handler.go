package auth

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"skillKonnect/app/db"
	"skillKonnect/app/models"
	"strconv"
	"strings"
	"time"

	"github.com/anthdm/superkit/kit"
	v "github.com/anthdm/superkit/validate"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	userSessionName = "user-session"
)

var authSchema = v.Schema{
	"email":    v.Rules(v.Email),
	"password": v.Rules(v.Required),
}

func HandleLoginIndex(kit *kit.Kit) error {
	if kit.Auth().Check() {
		redirectURL := kit.Getenv("SUPERKIT_AUTH_REDIRECT_AFTER_LOGIN", "/profile")
		return kit.Redirect(http.StatusSeeOther, redirectURL)
	}
	return kit.Render(LoginIndex(LoginIndexPageData{}))
}

func HandleLoginCreate(kit *kit.Kit) error {
	var values LoginFormValues
	errors, ok := v.Request(kit.Request, &values, authSchema)
	if !ok {
		return kit.Render(LoginForm(values, errors))
	}

	var user models.User
	err := db.Get().Find(&user, "email = ?", values.Email).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			errors.Add("credentials", "invalid credentials")
			return kit.Render(LoginForm(values, errors))
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(values.Password))
	if err != nil {
		log.Printf("err in CompareHashAndPassword at HandleLoginCreate: %v", err)
		errors.Add("credentials", "invalid credentials")
		return kit.Render(LoginForm(values, errors))
	}

	skipVerify := kit.Getenv("SUPERKIT_AUTH_SKIP_VERIFY", "false")
	if skipVerify != "true" {
		if !user.EmailVerifiedAt.Valid {
			errors.Add("verified", "please verify your email")
			return kit.Render(LoginForm(values, errors))
		}
	}

	sessionExpiryStr := kit.Getenv("SUPERKIT_AUTH_SESSION_EXPIRY_IN_HOURS", "48")
	sessionExpiry, err := strconv.Atoi(sessionExpiryStr)
	if err != nil {
		sessionExpiry = 48
	}
	session := models.Session{
		UserID:    user.ID,
		Token:     uuid.New().String(),
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(sessionExpiry)),
	}
	log.Printf("Creating session token %s with expiry %s", session.Token, session.ExpiresAt)
	if err = db.Get().Create(&session).Error; err != nil {
		return err
	}

	sess := kit.GetSession(userSessionName)
	sess.Values["sessionToken"] = session.Token
	sess.Save(kit.Request, kit.Response)
	redirectURL := kit.Getenv("SUPERKIT_AUTH_REDIRECT_AFTER_LOGIN", "/profile")

	return kit.Redirect(http.StatusSeeOther, redirectURL)
}

func HandleLoginDelete(kit *kit.Kit) error {
	sess := kit.GetSession(userSessionName)
	defer func() {
		sess.Values = map[any]any{}
		sess.Save(kit.Request, kit.Response)
	}()
	err := db.Get().Delete(&models.Session{}, "token = ?", sess.Values["sessionToken"]).Error
	if err != nil {
		return err
	}
	return kit.Redirect(http.StatusSeeOther, "/")
}

func HandleEmailVerify(kit *kit.Kit) error {
	tokenStr := kit.Request.URL.Query().Get("token")
	if len(tokenStr) == 0 {
		return kit.Render(EmailVerificationError("invalid verification token"))
	}

	token, err := jwt.ParseWithClaims(
		tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
			return []byte(os.Getenv("SUPERKIT_SECRET")), nil
		}, jwt.WithLeeway(5*time.Second))
	if err != nil {
		return kit.Render(EmailVerificationError("invalid verification token"))
	}
	if !token.Valid {
		return kit.Render(EmailVerificationError("invalid verification token"))
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return kit.Render(EmailVerificationError("invalid verification token"))
	}
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return kit.Render(EmailVerificationError("Email verification token expired"))
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return kit.Render(EmailVerificationError("Email verification token expired"))
	}

	var user models.User
	err = db.Get().First(&user, userID).Error
	if err != nil {
		return err
	}

	if user.EmailVerifiedAt.Time.After(time.Time{}) {
		return kit.Render(EmailVerificationError("Email already verified"))
	}

	now := sql.NullTime{Time: time.Now(), Valid: true}
	user.EmailVerifiedAt = now
	err = db.Get().Save(&user).Error
	if err != nil {
		return err
	}

	return kit.Redirect(http.StatusSeeOther, "/login")
}

func AuthenticateUser(kit *kit.Kit) (models.ExtendedAuth, error) {
	auth := models.AuthPayload{}
	sess := kit.GetSession(userSessionName)
	token, ok := sess.Values["sessionToken"]
	if !ok {
		return auth, nil
	}

	var session models.Session
	err := db.Get().
		Preload("User").
		Where("token = ? AND expires_at > ?", token, time.Now()).
		First(&session).Error
	if err != nil || session.ID == 0 {
		return auth, nil
	}

	return models.AuthPayload{
		Authenticated: true,
		User:          &session.User,
		Token:         session.Token,
	}, nil
}

func WebUIAuthFunc(kit *kit.Kit) (kit.Auth, error) {
	auth, err := AuthenticateUser(kit)
	if err != nil {
		return &models.AuthPayload{Authenticated: false}, nil
	}

	if !auth.Check() {
		return &models.AuthPayload{Authenticated: false}, nil
	}

	return &models.AuthPayload{
		Authenticated: true,
		User:          auth.GetUser(), // if your AuthenticateUser returns a user
		Token:         "",
	}, nil
}

func APIAuthFunc(kit *kit.Kit) (kit.Auth, error) {
	header := kit.Request.Header.Get("Authorization")
	if header == "" || !strings.HasPrefix(header, "Bearer ") {
		return &models.AuthPayload{Authenticated: false}, nil
	}

	token := strings.TrimPrefix(header, "Bearer ")

	var session models.Session
	if err := db.Get().Where("token = ? AND expires_at > ?", token, time.Now()).
		First(&session).Error; err != nil {
		return &models.AuthPayload{Authenticated: false}, nil
	}

	var user models.User
	db.Get().First(&user, session.UserID)

	return &models.AuthPayload{
		Authenticated: true,
		User:          &user,
		Token:         token,
	}, nil
}

func RequireWebAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		kit := &kit.Kit{
			Response: w,
			Request:  r,
		}

		// Retrieve session token
		sess := kit.GetSession(userSessionName)
		token, ok := sess.Values["sessionToken"]
		if !ok {
			http.Error(w, "unauthorized: missing session", http.StatusUnauthorized)
			return
		}

		// Load session and user
		var session models.Session
		err := db.Get().
			Preload("User").
			Where("token = ? AND expires_at > ?", token, time.Now()).
			First(&session).Error
		if err != nil || session.ID == 0 {
			http.Error(w, "unauthorized: invalid or expired session", http.StatusUnauthorized)
			return
		}

		user := session.User
		log.Printf("Authenticated user: %s (%s)", user.Email, user.Role)

		// Check admin role
		if user.Role != "admin" {
			http.Error(w, "forbidden: admin access only", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func RequireAdminAPI(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		payload, ok := r.Context().Value(models.AuthContextKey{}).(models.AuthPayload)
		if !ok || !payload.Authenticated {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if payload.User.Role != "admin" {
			http.Error(w, "forbidden: admin only", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
