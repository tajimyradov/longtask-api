package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"longtask-api/api"
	"longtask-api/task"
)

func main() {
	manager := task.NewManager()
	go manager.StartWorkerPool(3)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/tasks", api.CreateTaskHandler(manager))
	r.Get("/tasks/{id}", api.GetTaskHandler(manager))

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}
