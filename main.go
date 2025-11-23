// package main

// import (
// 	"context"
// 	"log"
// 	"net/http"
// 	"os"
// 	"time"

// 	"github.com/jackc/pgx/v5/pgxpool"
// 	"github.com/joho/godotenv"

// 	"github.com/fpswan/anycraft-backend/internal/controller"
// 	"github.com/fpswan/anycraft-backend/internal/middleware"
// 	"github.com/fpswan/anycraft-backend/internal/repository"
// 	"github.com/fpswan/anycraft-backend/internal/service"
// )

// func main() {
// 	_ = godotenv.Load()

// 	dsn := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/anycraft?sslmode=disable")
// 	port := getEnv("PORT", "8080")

// 	ctx := context.Background()
// 	db, err := pgxpool.New(ctx, dsn)
// 	if err != nil {
// 		log.Fatalf("DB connect failed: %v", err)
// 	}
// 	defer db.Close()

// 	repo := repository.NewComposeRepository(db)
// 	svc := service.NewComposeService(repo)
// 	ctrl := controller.NewComposeController(svc)

// 	mux := http.NewServeMux()

// 	// in main.go, before ctrl.RegisterRoutes(mux)
// 	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
// 		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
// 		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
// 		if r.Method == "OPTIONS" {
// 			w.WriteHeader(200)
// 			return
// 		}
// 		http.NotFound(w, r)
// 	})
// 	// then register real routes under a subpath, e.g., /api/v1/... already OK

// 	ctrl.RegisterRoutes(mux)

// 	handler := middleware.Recoverer(
// 		middleware.RequestLogger(&middleware.LoggerOptions{
// 			Writer:  os.Stdout, // or nil to use default
// 			Verbose: false,     // set true to also log headers
// 		})(mux),
// 	)

// 	// srv := &http.Server{
// 	// 	Addr:         ":" + port,
// 	// 	Handler:      mux,
// 	// 	ReadTimeout:  5 * time.Second,
// 	// 	WriteTimeout: 10 * time.Second,
// 	// }

// 	srv := &http.Server{
// 		Addr:              ":8080",
// 		Handler:           handler,
// 		ReadTimeout:       10 * time.Second,
// 		ReadHeaderTimeout: 5 * time.Second,
// 		WriteTimeout:      30 * time.Second,
// 		IdleTimeout:       90 * time.Second,
// 	}

// 	log.Printf("✅ API running at http://localhost:%s", port)
// 	log.Fatal(srv.ListenAndServe())
// }

// func getEnv(key, def string) string {
// 	if v := os.Getenv(key); v != "" {
// 		return v
// 	}
// 	return def
// }

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/fpswan/anycraft-backend/internal/controller"
	"github.com/fpswan/anycraft-backend/internal/middleware"
	"github.com/fpswan/anycraft-backend/internal/repository"
	"github.com/fpswan/anycraft-backend/internal/service"
)

func main() {
	_ = godotenv.Load()

	dsn := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/anycraft?sslmode=disable")
	port := getEnv("PORT", "8080")

	ctx := context.Background()
	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("DB connect failed: %v", err)
	}
	defer db.Close()

	repo := repository.NewComposeRepository(db)
	svc := service.NewComposeService(repo)
	ctrl := controller.NewComposeController(svc)

	mux := http.NewServeMux()

	// ❌ Remove the old "/" handler — it never ran for /api/... and caused the CORS miss

	// Register real routes under /api/v1/...
	ctrl.RegisterRoutes(mux)

	// ---- Middleware chain (outermost first) ----
	handler := middleware.CORS(&middleware.CORSOptions{
		// You can also set CORS_ORIGINS env var like: "http://localhost:3000,https://your.site"
		AllowOrigins:     getEnv("CORS_ORIGINS", "http://localhost:3000"),
		AllowCredentials: false, // set true if you need cookies/auth with specific origins (not "*")
		AllowHeaders:     "Content-Type,Authorization,X-Requested-With",
		AllowMethods:     "GET,POST,OPTIONS",
	})(
		middleware.Recoverer(
			middleware.RequestLogger(&middleware.LoggerOptions{
				Writer:  os.Stdout,
				Verbose: false,
			})(mux),
		),
	)

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       90 * time.Second,
	}

	log.Printf("✅ API running at http://localhost:%s", port)
	log.Fatal(srv.ListenAndServe())
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
