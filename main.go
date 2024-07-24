// main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type Poll struct {
	ID       string   `json:"id"`
	Question string   `json:"question"`
	Options  []string `json:"options"`
}

type Vote struct {
	PollID string `json:"poll_id"`
	Option string `json:"option"`
}

type PollResult struct {
	Poll    Poll           `json:"poll"`
	Results map[string]int `json:"results"`
}

var (
	polls     = make(map[string]*Poll)
	votes     = make(map[string]map[string]int)
	pollMutex sync.RWMutex
	voteMutex sync.RWMutex
)

func main() {
	http.HandleFunc("/create_poll", createPollHandler)
	http.HandleFunc("/vote", voteHandler)
	http.HandleFunc("/results/{id}", resultsHandler)
	http.HandleFunc("/poll_updates/{id}", pollUpdatesHandler)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createPollHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var poll Poll
	err := json.NewDecoder(r.Body).Decode(&poll)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pollMutex.Lock()
	polls[poll.ID] = &poll
	pollMutex.Unlock()

	voteMutex.Lock()
	votes[poll.ID] = make(map[string]int)
	voteMutex.Unlock()

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(poll)
}

func voteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var vote Vote
	err := json.NewDecoder(r.Body).Decode(&vote)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	voteMutex.Lock()
	if _, ok := votes[vote.PollID]; !ok {
		voteMutex.Unlock()
		http.Error(w, "Poll not found", http.StatusNotFound)
		return
	}
	votes[vote.PollID][vote.Option]++
	voteMutex.Unlock()

	w.WriteHeader(http.StatusOK)
}

func resultsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pollID := r.PathValue("id")
	if pollID == "" {
		http.Error(w, "Missing poll_id parameter", http.StatusBadRequest)
		return
	}

	pollMutex.RLock()
	poll, ok := polls[pollID]
	pollMutex.RUnlock()

	if !ok {
		http.Error(w, "Poll not found", http.StatusNotFound)
		return
	}

	voteMutex.RLock()
	results := votes[pollID]
	voteMutex.RUnlock()

	pollResult := PollResult{
		Poll:    *poll,
		Results: results,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(pollResult)
}

func pollUpdatesHandler(w http.ResponseWriter, r *http.Request) {
	pollID := r.PathValue("id")
	if pollID == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {
		voteMutex.RLock()
		results := votes[pollID]
		voteMutex.RUnlock()

		_, _ = fmt.Fprintf(w, "data: %v\n\n", results)
		w.(http.Flusher).Flush()

		// Simulate delay between updates
		// In a real application, you'd use channels to notify of changes
		// and avoid busy-waiting
		time.Sleep(3 * time.Second)
	}
}
