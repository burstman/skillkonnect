package auth

import (
	"skillKonnect/app/handlers"

	"github.com/anthdm/superkit/kit"
	"github.com/go-chi/chi/v5"
)

func InitializeRoutes(router chi.Router) {

	//----------------------------------------------------------------------
	// 1) WEB UI AUTH
	//----------------------------------------------------------------------
	authWebConfig := kit.AuthenticationConfig{
		AuthFunc:    WebUIAuthFunc,
		RedirectURL: "/web/admin/login",
	}

	// Public Web UI login
	router.Group(func(r chi.Router) {
		r.Use(kit.WithAuthentication(authWebConfig, false)) // allow unauth
		r.Get("/web/admin/login", kit.Handler(HandleLoginIndex))
		r.Post("/web/admin/login", kit.Handler(HandleLoginCreate))
		r.Delete("/web/admin/logout", kit.Handler(HandleLoginDelete))
	})

	// Protected Web UI routes
	router.Group(func(r chi.Router) {
		r.Use(kit.WithAuthentication(authWebConfig, true))
		r.Get("/", kit.Handler(HandleLoginIndex))
		r.Get("/profile", kit.Handler(HandleProfileShow))
		r.Put("/profile", kit.Handler(HandleProfileUpdate))
	})

	//----------------------------------------------------------------------
	// 2) API AUTH (JSON)
	//----------------------------------------------------------------------
	authApiConfig := kit.AuthenticationConfig{
		AuthFunc:    APIAuthFunc,
		RedirectURL: "", // no HTML redirects for API
	}

	// Public API login (NO AUTH) - in its own group
	router.Group(func(r chi.Router) {
		r.Post("/api/admin/login", kit.Handler(HandleApiLoginCreate))
	})

	// Protected API routes
	router.Group(func(api chi.Router) {
		api.Use(kit.WithAuthentication(authApiConfig, true))
		api.Use(RequireAdmin)

		api.Delete("/api/admin/logout", kit.Handler(HandleApiLoginDelete))

		// ADMIN API
		api.Route("/api/admin", func(r chi.Router) {

			// USERS
			r.Route("/users", func(r chi.Router) {
				r.Get("/", kit.Handler(handlers.WebAdminListUsers))
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
