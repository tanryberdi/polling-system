package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"polling-system/domain"
	"polling-system/ports"
)

type HTTPHandler struct {
	pollService ports.PollService
}

func NewHTTPHandler(pollService ports.PollService) *HTTPHandler {
	return &HTTPHandler{pollService: pollService}
}

func (h *HTTPHandler) CreatePollHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var poll domain.Poll
	err := json.NewDecoder(r.Body).Decode(&poll)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.pollService.CreatePoll(poll)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(poll)
}

func (h *HTTPHandler) VoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var vote domain.Vote
	err := json.NewDecoder(r.Body).Decode(&vote)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if PollID is provided
	if vote.PollID == "" {
		http.Error(w, "Missing poll_id", http.StatusBadRequest)
		return
	}

	err = h.pollService.Vote(vote)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *HTTPHandler) ResultsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pollID := r.PathValue("id")
	if pollID == "" {
		http.Error(w, "Missing poll_id parameter", http.StatusBadRequest)
		return
	}

	result, err := h.pollService.GetResults(pollID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *HTTPHandler) PollUpdatesHandler(w http.ResponseWriter, r *http.Request) {
	pollID := r.PathValue("id")
	if pollID == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {
		result, err := h.pollService.GetResults(pollID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "result for pollID=%v: %v\n\n", pollID, result.Results)
		w.(http.Flusher).Flush()

		time.Sleep(3 * time.Second)
	}
}
