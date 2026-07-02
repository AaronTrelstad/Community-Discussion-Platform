package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

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
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://dev:dev@127.0.0.1:5432/agentsandbox?sslmode=disable"
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	queries := db.New(conn)

	authService := services.NewAuthService(queries)
	teamService := services.NewTeamService(queries)

	authHandler := handlers.NewAuthHandler(authService)
	teamHandler := handlers.NewTeamHandler(teamService)

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

	r.Get("/teams", teamHandler.ListTeams)

	r.Group(func(r chi.Router) {
		r.Use(internalmiddleware.Auth)

		r.Get("/auth/me", authHandler.Me)

		r.Post("/teams", teamHandler.CreateTeam)
		r.Put("/teams/{id}", teamHandler.UpdateTeam)
		r.Get("/teams/{id}", teamHandler.GetTeam)
		r.Post("/teams/{id}/join", teamHandler.JoinTeam)
		r.Delete("/teams/{id}/leave", teamHandler.LeaveTeam)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Running on :%s", port)
	http.ListenAndServe(":"+port, r)
}
