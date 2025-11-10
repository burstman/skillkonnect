package auth

import (
	"github.com/anthdm/superkit/kit"
	"github.com/go-chi/chi/v5"
)

func InitializeRoutes(router chi.Router) {
	authConfig := kit.AuthenticationConfig{
		AuthFunc:    AuthenticateUser,
		RedirectURL: "/api/admin/login",
	}

	//router.Get("/email/verify", kit.Handler(HandleEmailVerify))
	//router.Post("/resend-email-verification", kit.Handler(HandleResendVerificationCode))

	router.Group(func(auth chi.Router) {
		auth.Use(kit.WithAuthentication(authConfig, false))
		auth.Get("/api/admin/login", kit.Handler(HandleLoginIndex))
		auth.Post("/api/admin/login", kit.Handler(HandleLoginCreate))
		auth.Delete("/api/admin/logout", kit.Handler(HandleLoginDelete))

		//auth.Get("/signup", kit.Handler(HandleSignupIndex))
		//auth.Post("/signup", kit.Handler(HandleSignupCreate))

	})

	router.Group(func(auth chi.Router) {
		auth.Use(kit.WithAuthentication(authConfig, true))
		auth.Get("/", kit.Handler(HandleLoginIndex))
		auth.Get("/profile", kit.Handler(HandleProfileShow))
		auth.Put("/profile", kit.Handler(HandleProfileUpdate))
	})
}
