package main

import (
	"awesomeProject1/internal/handlers"
	"awesomeProject1/internal/middleware"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	taskHandler := handlers.NewTaskHandler()

	mux.Handle("/tasks", middleware.LoggingMiddleware(
		middleware.APIKeyMiddleware(taskHandler),
	))

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
