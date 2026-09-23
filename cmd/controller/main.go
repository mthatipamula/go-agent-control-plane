package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mthatipamula/go-agent-control-plane/internal/agent"
	"github.com/mthatipamula/go-agent-control-plane/internal/api"
	"github.com/mthatipamula/go-agent-control-plane/internal/controller"
	"github.com/mthatipamula/go-agent-control-plane/internal/store"
	"github.com/mthatipamula/go-agent-control-plane/internal/task"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://control_plane:control_plane@localhost:5432/control_plane?sslmode=disable"
	}

	db, err := store.OpenPostgres(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	taskStore := store.NewPostgresTaskStore(db)
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

	if err := http.ListenAndServe(":8080", enableCORS(mux)); err != nil {
		log.Fatal(err)
	}
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
