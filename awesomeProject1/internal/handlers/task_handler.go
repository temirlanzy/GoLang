package handlers

import (
	"awesomeProject1/internal/models"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
)

type TaskHandler struct {
	tasks  map[int]models.Task
	mu     sync.Mutex
	nextID int
}

func NewTaskHandler() *TaskHandler {
	return &TaskHandler{
		tasks:  make(map[int]models.Task),
		nextID: 1,
	}
}

func (h *TaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		h.handleGet(w, r)
	case http.MethodPost:
		h.handlePost(w, r)
	case http.MethodPatch:
		h.handlePatch(w, r)
	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func (h *TaskHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	h.mu.Lock()
	defer h.mu.Unlock()

	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
			return
		}

		task, ok := h.tasks[id]
		if !ok {
			http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(task)
		return
	}

	var list []models.Task
	for _, t := range h.tasks {
		list = append(list, t)
	}

	json.NewEncoder(w).Encode(list)
}

func (h *TaskHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Title == "" {
		http.Error(w, `{"error":"invalid title"}`, http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	task := models.Task{
		ID:    h.nextID,
		Title: input.Title,
		Done:  false,
	}

	h.tasks[h.nextID] = task
	h.nextID++

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) handlePatch(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	var input struct {
		Done bool `json:"done"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	task, ok := h.tasks[id]
	if !ok {
		http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
		return
	}

	task.Done = input.Done
	h.tasks[id] = task

	json.NewEncoder(w).Encode(map[string]bool{"updated": true})
}
