package auth

import (
	"context"
	"net/http"
	"skillKonnect/app/models"
	"strings"

	"github.com/anthdm/superkit/kit"
)

// Auth object for context
type AuthContextKey struct{}

type AuthPayload struct {
	Authenticated bool
	User          *models.User
}

// Middleware config
type UnifiedAuthConfig struct {
	WebAuthFunc func(*kit.Kit) (*models.User, error)
	APIAuthFunc func(*kit.Kit) (*models.User, error)
	LoginURL    string
}

// Detect JSON / API request
func isJSONRequest(r *http.Request) bool {
	accept := r.Header.Get("Accept")
	ct := r.Header.Get("Content-Type")

	return strings.Contains(accept, "application/json") ||
		strings.Contains(ct, "application/json") ||
		strings.HasPrefix(r.Header.Get("Authorization"), "Bearer")
}

func WithUnifiedAuth(cfg UnifiedAuthConfig, strict bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			kit := &kit.Kit{
				Response: w,
				Request:  r,
			}

			isAPI := isJSONRequest(r)

			var user *models.User
			var err error

			if isAPI {
				// JSON API → Bearer Token
				user, err = cfg.APIAuthFunc(kit)
			} else {
				// Web UI → Cookie session
				user, err = cfg.WebAuthFunc(kit)
			}

			authenticated := (err == nil && user != nil)

			if strict && !authenticated {
				if isAPI {
					// JSON API → return 401
					kit.JSON(http.StatusUnauthorized, map[string]string{
						"error": "unauthorized",
					})
					return
				}

				// Web UI → redirect to login
				kit.Redirect(http.StatusSeeOther, cfg.LoginURL)
				return
			}

			// Inject auth into context
			ctx := context.WithValue(r.Context(), AuthContextKey{}, AuthPayload{
				Authenticated: authenticated,
				User:          user,
			})

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
