package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/mthatipamula/go-agent-control-plane/internal/agent"
	"github.com/mthatipamula/go-agent-control-plane/internal/api"
	"github.com/mthatipamula/go-agent-control-plane/internal/controller"
	"github.com/mthatipamula/go-agent-control-plane/internal/store"
	"github.com/mthatipamula/go-agent-control-plane/internal/task"
)

func main() {
	taskStore := store.NewTaskStore()
	controlPlane := controller.NewController(taskStore)

	executor := task.NewSimpleExecutor()
	worker := agent.NewAgent("agent-1", taskStore, executor)

	ctx := context.Background()
	go worker.Run(ctx, time.Second)

	handler := api.NewHandler(controlPlane)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.CreateTask(w, r)
			return
		}

		if r.Method == http.MethodGet {
			handler.ListTasks(w, r)
			return
		}

		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	})
	mux.HandleFunc("/api/tasks/", handler.GetTask)

	log.Println("control plane listening on :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
