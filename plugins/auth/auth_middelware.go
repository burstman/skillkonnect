package auth

import (
	"context"
	"net/http"
	"strings"

	"skillKonnect/app/models"

	"github.com/anthdm/superkit/kit"
)

type AuthContextKey struct{}

type AuthPayload struct {
	Authenticated bool
	User          *models.User
}

// Unified config for API + Web UI
type UnifiedAuthConfig struct {
	WebAuthFunc func(*kit.Kit) (*models.User, error)
	APIAuthFunc func(*kit.Kit) (*models.User, error)
	LoginURL    string
}

// Detect JSON / API request
func isJSONRequest(r *http.Request) bool {
	accept := r.Header.Get("Accept")
	ct := r.Header.Get("Content-Type")
	auth := r.Header.Get("Authorization")

	return strings.Contains(accept, "application/json") ||
		strings.Contains(ct, "application/json") ||
		strings.HasPrefix(auth, "Bearer ")
}

func WithUnifiedAuth(cfg UnifiedAuthConfig, strict bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			kit := &kit.Kit{Response: w, Request: r}

			isAPI := isJSONRequest(r)

			var user *models.User
			var err error

			if isAPI {
				user, err = cfg.APIAuthFunc(kit)
			} else {
				user, err = cfg.WebAuthFunc(kit)
			}

			auth := (err == nil && user != nil)

			// Strict mode → block unauthorized requests
			if strict && !auth {

				if isAPI {
					// JSON API → return 401 JSON
					_ = kit.JSON(http.StatusUnauthorized, map[string]string{
						"error": "unauthorized",
					})
					return
				}

				// WEB → redirect to login
				kit.Redirect(http.StatusSeeOther, cfg.LoginURL)
				return
			}

			// Inject into context
			ctx := context.WithValue(r.Context(), AuthContextKey{}, AuthPayload{
				Authenticated: auth,
				User:          user,
			})

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
