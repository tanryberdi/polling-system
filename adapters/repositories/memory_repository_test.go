package repositories

import (
	"sync"
	"testing"

	"polling-system/domain"
)

func TestCreatePoll(t *testing.T) {
	repo := NewMemoryRepository()

	poll := domain.Poll{
		ID:       "1",
		Question: "Test question?",
		Options:  []string{"Option 1", "Option 2"},
	}

	err := repo.CreatePoll(poll)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify the poll was created
	createdPoll, err := repo.GetPoll("1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if createdPoll.Question != poll.Question {
		t.Errorf("Expected question %s, got %s", poll.Question, createdPoll.Question)
	}
}

func TestVote(t *testing.T) {
	repo := NewMemoryRepository()

	// Create a poll first
	poll := domain.Poll{
		ID:       "1",
		Question: "Test question?",
		Options:  []string{"Option 1", "Option 2"},
	}
	repo.CreatePoll(poll)

	vote := domain.Vote{
		PollID: "1",
		Option: "Option 1",
	}

	err := repo.Vote(vote)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify the vote was recorded
	results, err := repo.GetResults("1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if results.Results["Option 1"] != 1 {
		t.Errorf("Expected 1 vote for Option 1, got %d", results.Results["Option 1"])
	}
}

func TestGetResults(t *testing.T) {
	repo := NewMemoryRepository()

	// Create a poll and add some votes
	poll := domain.Poll{
		ID:       "1",
		Question: "Test question?",
		Options:  []string{"Option 1", "Option 2"},
	}
	repo.CreatePoll(poll)
	repo.Vote(domain.Vote{PollID: "1", Option: "Option 1"})
	repo.Vote(domain.Vote{PollID: "1", Option: "Option 2"})
	repo.Vote(domain.Vote{PollID: "1", Option: "Option 1"})

	results, err := repo.GetResults("1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if results.Results["Option 1"] != 2 || results.Results["Option 2"] != 1 {
		t.Errorf("Unexpected results: %v", results.Results)
	}
}

func TestVoteNonExistentPoll(t *testing.T) {
	repo := NewMemoryRepository()

	vote := domain.Vote{
		PollID: "non_existent",
		Option: "Option 1",
	}

	err := repo.Vote(vote)
	if err == nil {
		t.Error("Expected an error when voting on a non-existent poll, got nil")
	}
}

func TestGetResultsNonExistentPoll(t *testing.T) {
	repo := NewMemoryRepository()

	_, err := repo.GetResults("non_existent")
	if err == nil {
		t.Error("Expected an error when getting results for a non-existent poll, got nil")
	}
}

//

func TestConcurrentVoting(t *testing.T) {
	repo := NewMemoryRepository()

	poll := domain.Poll{
		ID:       "1",
		Question: "Test question?",
		Options:  []string{"Option 1", "Option 2"},
	}
	err := repo.CreatePoll(poll)
	if err != nil {
		t.Fatalf("Failed to create poll: %v", err)
	}

	votesCount := 1000
	var wg sync.WaitGroup
	wg.Add(votesCount)

	// Create a channel to collect errors
	errorChan := make(chan error, votesCount)

	for i := 0; i < votesCount; i++ {
		go func() {
			defer wg.Done()
			vote := domain.Vote{
				PollID: "1",
				Option: "Option 1",
			}
			if err := repo.Vote(vote); err != nil {
				errorChan <- err
			}
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(errorChan)

	// Check for any errors
	for err := range errorChan {
		t.Errorf("Error during concurrent voting: %v", err)
	}

	// Verify the final vote count
	results, err := repo.GetResults("1")
	if err != nil {
		t.Fatalf("Failed to get results: %v", err)
	}

	if results.Results["Option 1"] != votesCount {
		t.Errorf("Expected %d votes, got %d", votesCount, results.Results["Option 1"])
	}
}

/*
func TestCreateExistingPoll(t *testing.T) {
	repo := NewMemoryRepository()

	poll := domain.Poll{
		ID:       "1",
		Question: "Test question?",
		Options:  []string{"Option 1", "Option 2"},
	}

	err := repo.CreatePoll(poll)
	if err != nil {
		t.Errorf("Unexpected error creating poll: %v", err)
	}

	// Attempt to create the same poll again
	err = repo.CreatePoll(poll)
	if err == nil {
		t.Error("Expected an error when creating a poll with an existing ID, got nil")
	}
}

*/
