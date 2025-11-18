package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"skillKonnect/app"
	_ "skillKonnect/docs"
	"skillKonnect/public"

	"github.com/anthdm/superkit/kit"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	kit.Setup()
	router := chi.NewMux()

	app.InitializeMiddleware(router)

	if kit.IsDevelopment() {
		router.Handle("/public/*", disableCache(staticDev()))
	} else if kit.IsProduction() {
		router.Handle("/public/*", staticProd())
	}

	kit.UseErrorHandler(app.ErrorHandler)
	// Try to set Chi's NotFound handler on the concrete mux if available.
	// This avoids registering a wildcard route which can conflict with
	// other route registrations and cause method-not-allowed (405) errors.
	// Use the router's NotFound setter (method) to register the handler.
	// Using a wildcard route with HandleFunc can interfere with method
	// resolution and produce 405 responses, so prefer the explicit NotFound
	// setter when available.
	router.NotFound(kit.Handler(app.NotFoundHandler))

	app.InitializeRoutes(router)
	app.RegisterEvents()
	// Serve Swagger UI only in development
	if kit.IsDevelopment() {
		router.Get("/swagger/*", httpSwagger.WrapHandler)
	}

	listenAddr := os.Getenv("HTTP_LISTEN_ADDR")
	// In development link the full Templ proxy url.
	url := "http://localhost:7331"
	if kit.IsProduction() {
		router.Get("/swagger/*", httpSwagger.WrapHandler)
		url = fmt.Sprintf("http://localhost%s", listenAddr)
	}

	fmt.Printf("application running in %s at %s\n", kit.Env(), url)

	http.ListenAndServe(listenAddr, router)
}

func staticDev() http.Handler {
	return http.StripPrefix("/public/", http.FileServerFS(os.DirFS("public")))
}

func staticProd() http.Handler {
	return http.StripPrefix("/public/", http.FileServerFS(public.AssetsFS))
}

func disableCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
