package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"longtask-api/task"
)

type createTaskRequest struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

func CreateTaskHandler(manager *task.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		var t task.Task
		switch req.Type {
		case "long_task":
			t = task.NewLongTask(req.Payload)
		default:
			http.Error(w, "unsupported task type", http.StatusBadRequest)
			return
		}

		manager.EnqueueTask(t)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"id":     t.ID(),
			"status": "queued",
		})
	}
}

func GetTaskHandler(manager *task.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		taskResult, ok := manager.GetTask(id)
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(taskResult)
	}
}
