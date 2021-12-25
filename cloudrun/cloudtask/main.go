package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
)

func main() {
	ctx := context.Background()

	log.Print("starting server...")
	http.HandleFunc("/", handler)

	tasksClient, err := cloudtasks.NewClient(ctx)
	if err != nil {
		panic(err)
	}
	h := Handlers{
		tasks: tasksClient,
	}
	http.HandleFunc("/tasks/serve", h.ServeTask)
	http.HandleFunc("/tasks/add", h.AddTask)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	go func() {
		trap := make(chan os.Signal, 1)
		signal.Notify(trap, syscall.SIGTERM)
		s := <-trap
		fmt.Printf("Received shutdown signal %s\n", s)
		fmt.Printf("Shutdown gracefully....\n")
		if err := h.tasks.Close(); err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Shutdown finish!\n")
	}()

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello!\n")
}
