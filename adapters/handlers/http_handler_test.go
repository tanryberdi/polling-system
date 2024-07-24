package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"polling-system/domain"
	"polling-system/mocks"
)

func TestCreatePollHandler(t *testing.T) {
	mockService := mocks.NewMockPollService()
	handler := NewHTTPHandler(mockService)

	poll := domain.Poll{
		ID:       "1",
		Question: "Test question?",
		Options:  []string{"Option 1", "Option 2"},
	}

	body, _ := json.Marshal(poll)
	req := httptest.NewRequest("POST", "/create_poll", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.CreatePollHandler(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	var responsePoll domain.Poll
	err := json.Unmarshal(rr.Body.Bytes(), &responsePoll)
	if err != nil {
		t.Fatal(err)
	}

	if responsePoll.ID != poll.ID {
		t.Errorf("handler returned unexpected body: got %v want %v",
			responsePoll.ID, poll.ID)
	}
}

func TestVoteHandler(t *testing.T) {
	mockService := mocks.NewMockPollService()
	handler := NewHTTPHandler(mockService)

	// Create a poll first
	_ = mockService.CreatePoll(domain.Poll{ID: "1", Question: "Test?", Options: []string{"Option 1", "Option 2"}})

	vote := domain.Vote{
		PollID: "1",
		Option: "Option 1",
	}

	body, _ := json.Marshal(vote)
	req := httptest.NewRequest("POST", "/vote", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.VoteHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestVoteMultipleHandler(t *testing.T) {
	mockService := mocks.NewMockPollService()
	handler := NewHTTPHandler(mockService)

	// Create test polls
	_ = mockService.CreatePoll(domain.Poll{ID: "1", Question: "Test 1?", Options: []string{"Option 1", "Option 2"}})
	_ = mockService.CreatePoll(domain.Poll{ID: "2", Question: "Test 2?", Options: []string{"Option A", "Option B"}})

	votes := []domain.Vote{
		{PollID: "1", Option: "Option 1"},
		{PollID: "2", Option: "Option B"},
	}

	body, _ := json.Marshal(votes)
	req := httptest.NewRequest("POST", "/vote_multiple", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.VoteMultipleHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Verify votes were recorded
	result1, _ := mockService.GetResults("1")
	result2, _ := mockService.GetResults("2")

	if result1.Results["Option 1"] != 1 {
		t.Errorf("Expected 1 vote for Option 1 in poll 1, got %d", result1.Results["Option 1"])
	}

	if result2.Results["Option B"] != 1 {
		t.Errorf("Expected 1 vote for Option B in poll 2, got %d", result2.Results["Option B"])
	}
}

func TestResultsHandler(t *testing.T) {
	mockService := mocks.NewMockPollService()
	handler := NewHTTPHandler(mockService)

	// Create a test poll and add some votes
	poll := domain.Poll{
		ID:       "1",
		Question: "Test question?",
		Options:  []string{"Option 1", "Option 2"},
	}
	_ = mockService.CreatePoll(poll)
	_ = mockService.Vote(domain.Vote{PollID: "1", Option: "Option 1"})
	_ = mockService.Vote(domain.Vote{PollID: "1", Option: "Option 2"})

	req := httptest.NewRequest("GET", "/results/1", nil)

	rr := httptest.NewRecorder()
	handler.ResultsHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var result domain.PollResult
	err := json.Unmarshal(rr.Body.Bytes(), &result)
	if err != nil {
		t.Fatal(err)
	}

	if result.Poll.ID != "1" {
		t.Errorf("handler returned unexpected poll ID: got %v want %v",
			result.Poll.ID, "1")
	}

	if result.Results["Option 1"] != 1 || result.Results["Option 2"] != 1 {
		t.Errorf("handler returned unexpected results: %v", result.Results)
	}
}

/*
func TestPollUpdatesHandler(t *testing.T) {
	mockService := mocks.NewMockPollService()
	handler := NewHTTPHandler(mockService)

	// Create a test poll
	poll := domain.Poll{
		ID:       "1",
		Question: "Test question?",
		Options:  []string{"Option 1", "Option 2"},
	}
	mockService.CreatePoll(poll)

	req := httptest.NewRequest("GET", "/poll_updates/1", nil)

	rr := httptest.NewRecorder()

	// Use a channel to signal when the handler has written the first update
	done := make(chan bool)
	go func() {
		handler.PollUpdatesHandler(rr, req)
		done <- true
	}()

	// Add a vote to trigger an update
	mockService.Vote(domain.Vote{PollID: "1", Option: "Option 1"})

	// Wait for the handler to write the update or timeout
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Timeout waiting for handler to write update")
	}

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check that the response contains the expected SSE format
	expectedPrefix := "data: "
	if !bytes.HasPrefix(rr.Body.Bytes(), []byte(expectedPrefix)) {
		t.Errorf("handler response doesn't start with %q", expectedPrefix)
	}
}

*/
