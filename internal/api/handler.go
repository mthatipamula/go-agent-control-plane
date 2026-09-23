package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/mthatipamula/go-agent-control-plane/internal/controller"
	"github.com/mthatipamula/go-agent-control-plane/internal/task"
)

type Handler struct {
	controller *controller.Controller
}

func NewHandler(controller *controller.Controller) *Handler {
	return &Handler{
		controller: controller,
	}
}

type createTaskRequest struct {
	Payload string `json:"payload"`
}

type taskResponse struct {
	ID             string  `json:"id"`
	Payload        string  `json:"payload"`
	Status         string  `json:"status"`
	AgentID        string  `json:"agentId,omitempty"`
	Attempt        int     `json:"attempt"`
	Version        int64   `json:"version"`
	FencingToken   int64   `json:"fencingToken"`
	LeaseExpiresAt *string `json:"leaseExpiresAt,omitempty"`
	CreatedAt      string  `json:"createdAt"`
	UpdatedAt      string  `json:"updatedAt"`
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request createTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(request.Payload) == "" {
		http.Error(w, "payload is required", http.StatusBadRequest)
		return
	}

	t, err := h.controller.Submit(request.Payload)
	if err != nil {
		http.Error(w, "failed to create task", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, taskResponse{
		ID:           t.ID,
		Payload:      t.Payload,
		Status:       string(t.Status),
		Attempt:      t.Attempt,
		Version:      t.Version,
		FencingToken: t.FencingToken,
		CreatedAt:    t.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    t.UpdatedAt.Format(time.RFC3339),
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(value)
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	taskID := strings.TrimPrefix(r.URL.Path, "/api/tasks/")

	if taskID == "" {
		http.Error(w, "task id is required", http.StatusBadRequest)
		return
	}

	t, err := h.controller.Get(taskID)
	if err != nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, toTaskResponse(t))
}

func (h *Handler) ListTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tasks, err := h.controller.List()
	if err != nil {
		http.Error(w, "failed to list tasks", http.StatusInternalServerError)
		return
	}

	response := make([]taskResponse, 0, len(tasks))

	for _, t := range tasks {
		response = append(response, toTaskResponse(t))
	}

	writeJSON(w, http.StatusOK, response)
}

func toTaskResponse(t task.Task) taskResponse {
	var leaseExpiresAt *string

	if t.LeaseExpiresAt != nil {
		value := t.LeaseExpiresAt.Format(time.RFC3339)
		leaseExpiresAt = &value
	}

	return taskResponse{
		ID:             t.ID,
		Payload:        t.Payload,
		Status:         string(t.Status),
		AgentID:        t.AgentID,
		Attempt:        t.Attempt,
		Version:        t.Version,
		FencingToken:   t.FencingToken,
		LeaseExpiresAt: leaseExpiresAt,
		CreatedAt:      t.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      t.UpdatedAt.Format(time.RFC3339),
	}
}
