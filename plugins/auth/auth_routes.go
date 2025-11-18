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
		LoginURL:    "/web/admin/login",
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

	// Public API login (NO AUTH) - in its own group
	router.Group(func(r chi.Router) {
		r.Post("/api/admin/login", kit.Handler(HandleApiLoginCreate))
	})

	// Protected API routes
	router.Group(func(api chi.Router) {
		api.Use(WithUnifiedAuth(apiAuthConfig, true))
		api.Use(RequireAdminAPI)

		api.Delete("/api/admin/logout", kit.Handler(HandleApiLoginDelete))

		// ADMIN API
		api.Route("/api/admin", func(r chi.Router) {

			// USERS
			r.Route("/users", func(r chi.Router) {
				r.Get("/", kit.Handler(handlers.ApiAdminListUsers))
				r.Get("/{id}", kit.Handler(handlers.AdminGetUser))
				r.Put("/{id}/suspend", kit.Handler(handlers.AdminSuspendUser))
				r.Put("/{id}/activate", kit.Handler(handlers.AdminActivateUser))
			})

			// CATEGORIES
			r.Route("/categories", func(r chi.Router) {
				r.Get("/", kit.Handler(handlers.AdminListCategories))
				r.Post("/", kit.Handler(handlers.AdminCreateCategory))
				r.Delete("/{id}", kit.Handler(handlers.AdminDeleteCategory))
			})

			// SKILLS
			r.Route("/skills", func(r chi.Router) {
				r.Get("/", kit.Handler(handlers.AdminListSkills))
				r.Post("/", kit.Handler(handlers.AdminCreateSkill))
				r.Delete("/{id}", kit.Handler(handlers.AdminDeleteSkill))
			})

		})
	})
}
