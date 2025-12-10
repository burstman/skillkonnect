// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"skillKonnect/app"
	"skillKonnect/app/db"
	"skillKonnect/app/models"
	_ "skillKonnect/docs"
	"skillKonnect/public"

	"github.com/anthdm/superkit/kit"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	// Additional imports for db and models
	// "skillKonnect/app/db"
	// "skillKonnect/app/models"
)

func main() {
	kit.Setup()
	// --- AUTOMATIC DB MIGRATION ---
	// Run GORM AutoMigrate for all main models before starting server
	dbConn := db.Get()
	err := dbConn.AutoMigrate(
		&models.User{},
		&models.Skill{},
		&models.Category{},
		&models.Session{},
		&models.WorkerProfile{},
	)
	if err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	router := chi.NewMux()

	app.InitializeMiddleware(router)

	if kit.IsDevelopment() {
		router.Handle("/public/*", disableCache(staticDev()))
	} else if kit.IsProduction() {
		router.Handle("/public/*", staticProd())
	}

	kit.UseErrorHandler(app.ErrorHandler)
	router.NotFound(kit.Handler(app.NotFoundHandler))

	app.InitializeRoutes(router)
	app.RegisterEvents()
	if kit.IsDevelopment() {
		router.Get("/swagger/*", httpSwagger.WrapHandler)
	}

	listenAddr := os.Getenv("HTTP_LISTEN_ADDR")
	url := "http://localhost:8080"
	if kit.IsProduction() {
		router.Get("/swagger/*", httpSwagger.WrapHandler)
		url = fmt.Sprintf("http://localhost%s", listenAddr)
	}

	router.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir("public/uploads"))))

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
	// Optionally load environment variables from a .env file if it exists.
	// Missing .env is not fatal; runtime variables can be provided via Docker or the host env.
	if info, statErr := os.Stat(".env"); statErr == nil && !info.IsDir() {
		if err := godotenv.Load(); err != nil {
			log.Printf("warning: could not load .env: %v", err)
		}
	} else if statErr != nil && !os.IsNotExist(statErr) {
		log.Printf("warning: could not stat .env: %v", statErr)
	}
}
