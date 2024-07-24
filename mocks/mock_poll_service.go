package mocks

import (
	"fmt"

	"polling-system/domain"
)

type MockPollService struct {
	polls map[string]*domain.Poll
	votes map[string]map[string]int
}

func NewMockPollService() *MockPollService {
	return &MockPollService{
		polls: make(map[string]*domain.Poll),
		votes: make(map[string]map[string]int),
	}
}

func (m *MockPollService) CreatePoll(poll domain.Poll) error {
	m.polls[poll.ID] = &poll
	m.votes[poll.ID] = make(map[string]int)
	return nil
}

func (m *MockPollService) Vote(vote domain.Vote) error {
	if _, ok := m.votes[vote.PollID]; !ok {
		return fmt.Errorf("poll not found")
	}
	m.votes[vote.PollID][vote.Option]++
	return nil
}

func (m *MockPollService) GetResults(pollID string) (domain.PollResult, error) {
	poll, ok := m.polls[pollID]
	if !ok {
		return domain.PollResult{}, fmt.Errorf("poll not found")
	}
	return domain.PollResult{
		Poll:    *poll,
		Results: m.votes[pollID],
	}, nil
}
