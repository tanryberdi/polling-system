package main

import (
	"log"
	"net/http"

	"polling-system/adapters/handlers"
	"polling-system/adapters/repositories"
	"polling-system/adapters/services"
)

func main() {
	repo := repositories.NewMemoryRepository()
	pollService := services.NewPollService(repo)
	handler := handlers.NewHTTPHandler(pollService)

	http.HandleFunc("/create_poll", handler.CreatePollHandler)
	http.HandleFunc("/vote", handler.VoteHandler)
	http.HandleFunc("/results/{id}", handler.ResultsHandler)
	http.HandleFunc("/poll_updates/{id}", handler.PollUpdatesHandler)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
