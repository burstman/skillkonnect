package app

import (
	"log/slog"
	"skillKonnect/app/handlers"
	"skillKonnect/app/views/errors"
	"skillKonnect/plugins/auth"

	"github.com/anthdm/superkit/kit"
	"github.com/anthdm/superkit/kit/middleware"
	"github.com/go-chi/chi/v5"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

// Define your global middleware
func InitializeMiddleware(router *chi.Mux) {
	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.Recoverer)
	router.Use(middleware.WithRequest)
}

// Define your routes in here
func InitializeRoutes(router *chi.Mux) {
	// Authentication plugin
	//
	// By default the auth plugin is active, to disable the auth plugin
	// you will need to pass your own handler in the `AuthFunc`` field
	// of the `kit.AuthenticationConfig`.
	//  authConfig := kit.AuthenticationConfig{
	//      AuthFunc: YourAuthHandler,
	//      RedirectURL: "/login",
	//  }
	auth.InitializeRoutes(router)
	webcfg := kit.AuthenticationConfig{
		AuthFunc:    auth.WebUIAuthFunc,
		RedirectURL: "/api/admin/login",
	}

	// apiCfg := kit.AuthenticationConfig{
	// 	AuthFunc:    auth.APIAuthFunc,
	// 	RedirectURL: "", // not used for JSON
	// }

	// Routes that "might" have an authenticated user
	router.Group(func(app chi.Router) {
		app.Use(kit.WithAuthentication(webcfg, true)) // strict set to false
		app.Use(auth.RequireAdmin)
		// Routes that "must" have an authenticated user or else they
		// will be redirected to the configured redirectURL, set in the
		// AuthenticationConfig.
		app.Route("/web/admin", func(r chi.Router) {
			// Dashboard
			r.Get("/dashboard", kit.Handler(handlers.AdminDashboard))

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

			// // Workers
			// r.Route("/workers", func(r chi.Router) {
			// 	r.Get("/", kit.Handler(handlers.AdminListWorkers))
			// 	r.Get("/{id}", kit.Handler(handlers.AdminGetWorker))
			// 	r.Put("/{id}/approve", kit.Handler(handlers.AdminApproveWorker))
			// 	r.Put("/{id}/ban", kit.Handler(handlers.AdminBanWorker))
			// })
		})

	})

}

// NotFoundHandler that will be called when the requested path could
// not be found.
func NotFoundHandler(kit *kit.Kit) error {
	return kit.Render(errors.Error404())
}

// ErrorHandler that will be called on errors return from application handlers.
func ErrorHandler(kit *kit.Kit, err error) {
	slog.Error("internal server error", "err", err.Error(), "path", kit.Request.URL.Path)
	kit.Render(errors.Error500())
}
