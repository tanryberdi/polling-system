package repositories

import (
	"fmt"
	"sync"

	"polling-system/domain"
)

type MemoryRepository struct {
	polls     map[string]*domain.Poll
	votes     map[string]map[string]int
	pollMutex sync.RWMutex
	voteMutex sync.RWMutex
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		polls: make(map[string]*domain.Poll),
		votes: make(map[string]map[string]int),
	}
}

func (r *MemoryRepository) CreatePoll(poll domain.Poll) error {
	r.pollMutex.Lock()
	defer r.pollMutex.Unlock()
	r.polls[poll.ID] = &poll
	r.votes[poll.ID] = make(map[string]int)
	return nil
}

func (r *MemoryRepository) GetPoll(id string) (domain.Poll, error) {
	r.pollMutex.RLock()
	defer r.pollMutex.RUnlock()
	poll, ok := r.polls[id]
	if !ok {
		return domain.Poll{}, fmt.Errorf("poll not found")
	}
	return *poll, nil
}

func (r *MemoryRepository) Vote(vote domain.Vote) error {
	r.voteMutex.Lock()
	defer r.voteMutex.Unlock()

	r.pollMutex.RLock()
	_, ok := r.polls[vote.PollID]
	r.pollMutex.RUnlock()

	if !ok {
		return fmt.Errorf("poll not found")
	}

	if _, ok := r.votes[vote.PollID]; !ok {
		r.votes[vote.PollID] = make(map[string]int)
	}
	r.votes[vote.PollID][vote.Option]++
	return nil
}

func (r *MemoryRepository) GetResults(pollID string) (domain.PollResult, error) {
	r.pollMutex.RLock()
	defer r.pollMutex.RUnlock()
	r.voteMutex.RLock()
	defer r.voteMutex.RUnlock()

	poll, ok := r.polls[pollID]
	if !ok {
		return domain.PollResult{}, fmt.Errorf("poll not found")
	}

	return domain.PollResult{
		Poll:    *poll,
		Results: r.votes[pollID],
	}, nil
}
