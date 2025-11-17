package auth

import (
	"skillKonnect/app/handlers"

	"github.com/anthdm/superkit/kit"
	"github.com/go-chi/chi/v5"
)

func InitializeRoutes(router chi.Router) {
	authWebConfig := kit.AuthenticationConfig{
		AuthFunc:    WebUIAuthFunc,
		RedirectURL: "/web/admin/login",
	}

	authApiConfig := kit.AuthenticationConfig{
		AuthFunc:    APIAuthFunc,
		RedirectURL: "", // not used for JSON
	}

	//router.Get("/email/verify", kit.Handler(HandleEmailVerify))
	//router.Post("/resend-email-verification", kit.Handler(HandleResendVerificationCode))

	router.Group(func(auth chi.Router) {
		auth.Use(kit.WithAuthentication(authWebConfig, false))
		auth.Get("/web/admin/login", kit.Handler(HandleLoginIndex))
		auth.Post("/web/admin/login", kit.Handler(HandleLoginCreate))
		auth.Delete("/web/admin/logout", kit.Handler(HandleLoginDelete))

		//auth.Get("/signup", kit.Handler(HandleSignupIndex))
		//auth.Post("/signup", kit.Handler(HandleSignupCreate))

	})

	router.Group(func(auth chi.Router) {
		auth.Use(kit.WithAuthentication(authWebConfig, true))
		auth.Get("/", kit.Handler(HandleLoginIndex))
		auth.Get("/profile", kit.Handler(HandleProfileShow))
		auth.Put("/profile", kit.Handler(HandleProfileUpdate))
	})
	router.Group(func(api chi.Router) {
		api.Use(kit.WithAuthentication(authApiConfig, true))
		api.Use(RequireAdmin)

		api.Route("/api/admin", func(r chi.Router) {

			r.Route("/login", func(r chi.Router) {
				//r.Get("/", kit.Handler(HandleLoginIndex))
				r.Post("/", kit.Handler(HandleApiLoginCreate))
				r.Delete("/logout", kit.Handler(HandleApiLoginDelete))
			})

			// Users management
			r.Route("/users", func(r chi.Router) {
				r.Get("/", kit.Handler(handlers.AdminListUsers))
				r.Get("/{id}", kit.Handler(handlers.AdminGetUser))
				r.Put("/{id}/suspend", kit.Handler(handlers.AdminSuspendUser))
				r.Put("/{id}/activate", kit.Handler(handlers.AdminActivateUser))
			})

			// Skills and Categories
			r.Route("/categories", func(r chi.Router) {
				r.Get("/", kit.Handler(handlers.AdminListCategories))
				r.Post("/", kit.Handler(handlers.AdminCreateCategory))
				r.Delete("/{id}", kit.Handler(handlers.AdminDeleteCategory))
			})

			r.Route("/skills", func(r chi.Router) {
				r.Get("/", kit.Handler(handlers.AdminListSkills))
				r.Post("/", kit.Handler(handlers.AdminCreateSkill))
				r.Delete("/{id}", kit.Handler(handlers.AdminDeleteSkill))
			})

			// Protected API routes can be added here

		})

	})
}
