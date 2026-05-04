package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"todo-api/internal/httpapi"
	"todo-api/internal/todo"
)

func main() {
	addr := getenv("ADDR", ":8000")

	repo := todo.NewInMemoryRepo()
	svc := todo.NewService(repo, time.Now)
	srv := httpapi.NewServer(svc)

	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, srv.Handler()); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
