package main

import (
	"fmt"
	"log"

	"github.com/mthatipamula/go-agent-control-plane/internal/controller"
	"github.com/mthatipamula/go-agent-control-plane/internal/store"
)

func main() {
	taskStore := store.NewTaskStore()
	controlPlane := controller.NewController(taskStore)

	t, err := controlPlane.Submit("process order")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("task submitted: id=%s status=%s\n", t.ID, t.Status)
}