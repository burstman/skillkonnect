package auth

import (
	"skillKonnect/app/handlers"
	"skillKonnect/app/models"

	"github.com/anthdm/superkit/kit"
	"github.com/go-chi/chi/v5"
)

// Wrapper to adapt WebUIAuthFunc (returns kit.Auth) to UnifiedAuthConfig (expects *models.User)
func webUIAuthFuncAdapter(k *kit.Kit) (*models.User, error) {
	auth, err := WebUIAuthFunc(k)
	if err != nil {
		return nil, err
	}
	if auth == nil {
		return nil, nil
	}
	// kit.Auth is an interface, cast to *models.AuthPayload to get the user
	if payload, ok := auth.(*models.AuthPayload); ok && payload.Authenticated {
		return payload.User, nil
	}
	return nil, nil
}

// Wrapper to adapt APIAuthFunc (returns kit.Auth) to UnifiedAuthConfig (expects *models.User)
func apiAuthFuncAdapter(k *kit.Kit) (*models.User, error) {
	auth, err := APIAuthFunc(k)
	if err != nil {
		return nil, err
	}
	if auth == nil {
		return nil, nil
	}
	// kit.Auth is an interface, cast to *models.AuthPayload to get the user
	if payload, ok := auth.(*models.AuthPayload); ok && payload.Authenticated {
		return payload.User, nil
	}
	return nil, nil
}

func InitializeRoutes(router chi.Router) {

	//----------------------------------------------------------------------
	// 1) WEB UI AUTH
	//----------------------------------------------------------------------
	webAuthConfig := UnifiedAuthConfig{
		WebAuthFunc: webUIAuthFuncAdapter,
		APIAuthFunc: apiAuthFuncAdapter,
		LoginURL:    "",
	}

	// Public Web UI login
	router.Group(func(r chi.Router) {
		r.Use(WithUnifiedAuth(webAuthConfig, false)) // allow unauth
		r.Get("/web/admin/login", kit.Handler(HandleLoginIndex))
		r.Post("/web/admin/login", kit.Handler(HandleLoginCreate))
		r.Delete("/web/admin/logout", kit.Handler(HandleLoginDelete))
	})

	// Protected Web UI routes
	router.Group(func(r chi.Router) {
		r.Use(WithUnifiedAuth(webAuthConfig, true))
		r.Get("/", kit.Handler(HandleLoginIndex))
		r.Get("/profile", kit.Handler(HandleProfileShow))
		r.Put("/profile", kit.Handler(HandleProfileUpdate))
	})

	//----------------------------------------------------------------------
	// 2) API AUTH (JSON)
	//----------------------------------------------------------------------
	apiAuthConfig := UnifiedAuthConfig{
		WebAuthFunc: webUIAuthFuncAdapter,
		APIAuthFunc: apiAuthFuncAdapter,
		LoginURL:    "", // no HTML redirects for API
	}

	// Versioned API routes under /api/v1
	router.Route("/api/v1", func(v1 chi.Router) {
		// Health check endpoint (no auth required)
		v1.Get("/health", kit.Handler(HandleHealthCheck))

		// Public API login (NO AUTH)
		v1.Post("/admin/login", kit.Handler(HandleApiLoginCreate))
		v1.Post("/client/login", kit.Handler(HandleClientLoginCreate))

		// Authenticated routes (any logged-in user)
		v1.Group(func(api chi.Router) {
			api.Use(WithUnifiedAuth(apiAuthConfig, true))
			api.Get("/auth/me", kit.Handler(HandleApiAuthMe))
		})

		// User-specific routes (for authenticated users)
		v1.Group(func(api chi.Router) {
			api.Use(WithUnifiedAuth(apiAuthConfig, true))
			api.Route("/user", func(r chi.Router) {
				r.Get("/dashboard", kit.Handler(handlers.UserDashboardStats))
			})
		})

		// Client-specific routes (for authenticated clients)
		v1.Group(func(api chi.Router) {
			api.Use(WithUnifiedAuth(apiAuthConfig, true))
			api.Use(RequireRoleAPI("client"))
			api.Route("/client", func(r chi.Router) {
				r.Get("/dashboard", kit.Handler(handlers.ClientDashboardStats))
				r.Get("/workers", kit.Handler(handlers.ClientListWorkers))
				// Add more client endpoints here
			})
		})

		// Worker-specific routes (for authenticated workers)
		v1.Group(func(api chi.Router) {
			api.Use(WithUnifiedAuth(apiAuthConfig, true))
			api.Use(RequireRoleAPI("worker"))
			api.Route("/worker", func(r chi.Router) {
				r.Get("/dashboard", kit.Handler(handlers.WorkerDashboardStats))
				// Add more worker endpoints here
			})
		})

		// Protected admin routes
		v1.Group(func(api chi.Router) {
			api.Use(WithUnifiedAuth(apiAuthConfig, true))
			api.Use(RequireAdminAPI)

			api.Delete("/admin/logout", kit.Handler(HandleApiLoginDelete))

			api.Route("/admin", func(r chi.Router) {
				// USERS
				r.Route("/users", func(r chi.Router) {
					r.Get("/", kit.Handler(handlers.ApiAdminListUsers))
					r.Get("/{id}", kit.Handler(handlers.AdminGetUser))
					r.Put("/{id}/suspend", kit.Handler(handlers.AdminSuspendUser))
					r.Put("/{id}/activate", kit.Handler(handlers.AdminActivateUser))
					r.Put("/{id}", kit.Handler(handlers.AdminUpdateUser))
					r.Delete("/{id}", kit.Handler(handlers.AdminDeleteUser))
				})

				// CATEGORIES
				r.Route("/categories", func(r chi.Router) {
					r.Get("/", kit.Handler(handlers.AdminListCategories))
					r.Post("/", kit.Handler(handlers.AdminCreateCategory))
					r.Delete("/{id}", kit.Handler(handlers.AdminDeleteCategory))
					r.Put("/{id}", kit.Handler(handlers.AdminUpdateCategory))
				})

				// SKILLS
				r.Route("/skills", func(r chi.Router) {
					r.Get("/", kit.Handler(handlers.AdminListSkills))
					r.Post("/", kit.Handler(handlers.AdminCreateSkill))
					r.Put("/{id}", kit.Handler(handlers.AdminUpdateSkill))
					r.Delete("/{id}", kit.Handler(handlers.AdminDeleteSkill))
				})

				// STATUS
				r.Route("/status", func(r chi.Router) {
					r.Get("/dashboard", kit.Handler(handlers.AdminDashboardStats))
				})
			})
		})

		// Move the `/upload` route directly under `/api/v1` to avoid double prefix
		v1.Post("/upload", kit.Handler(handlers.AdminUploadFile))
	})
}

// @Summary Upload a file
// @Description Upload a file to the server. Allowed types: .png, .jpg, .jpeg, .webp, .pdf, .docx
// @Tags Admin
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Success 200 {object} map[string]string "{\"filename\":\"example.png\",\"url\":\"/uploads/example.png\"}"
// @Failure 400 {object} map[string]string "{\"error\":\"invalid_form\"}"
// @Failure 500 {object} map[string]string "{\"error\":\"file_save_error\"}"
// @Security BearerAuth
// @Router /api/v1/upload [post]
