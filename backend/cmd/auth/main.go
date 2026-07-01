package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/aarontrelstad/api/internal/db"
	"github.com/aarontrelstad/api/internal/handlers"
	internalmiddleware "github.com/aarontrelstad/api/internal/middleware"
	"github.com/aarontrelstad/api/internal/services"
	"github.com/aarontrelstad/api/pkg/httputil"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

func main() {
	dbURL := "postgres://dev:dev@127.0.0.1:5432/reddit?sslmode=disable"

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	queries := db.New(conn)
	authService := services.NewAuthService(queries)
	authHandler := handlers.NewAuthHandler(authService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(httputil.CORSMiddleware)

	r.Options("/*", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	r.Post("/auth/register", authHandler.Register)
	r.Post("/auth/login", authHandler.Login)
	r.Post("/auth/logout", authHandler.Logout)
	r.Post("/auth/refresh", authHandler.Refresh)

	r.Group(func(r chi.Router) {
		r.Use(internalmiddleware.Auth)
		r.Get("/auth/me", authHandler.Me)
	})

	log.Println("Running locally on :8080")
	http.ListenAndServe(":8080", r)
}
