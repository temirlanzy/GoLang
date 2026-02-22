package app

import (
	"context"
	"log"
	"net/http"
	"practice3/internal/handlers"
	"practice3/internal/middleware"
	"practice3/internal/repository"
	"practice3/internal/repository/_postgres"
	"practice3/internal/usecase"
	"practice3/pkg/modules"
	"time"
)

func Run() {
	ctx := context.Background()

	cfg := &modules.PostgreConfig{
		Host:        "localhost",
		Port:        "5432",
		Username:    "postgres",
		Password:    "postgres",
		DBName:      "mydb",
		SSLMode:     "disable",
		ExecTimeout: 5 * time.Second,
	}

	db := _postgres.NewPGXDialect(ctx, cfg)
	repos := repository.NewRepositories(db)
	use := usecase.NewUserUsecase(repos.UserRepository)
	handler := handlers.NewUserHandler(use)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	mux.Handle("/users", middleware.LoggingMiddleware(
		middleware.APIKeyMiddleware(http.HandlerFunc(handler.GetUsers)),
	))

	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", mux)
}
