package services

import (
	"fmt"

	"polling-system/domain"
	"polling-system/ports"
)

type PollService struct {
	repo ports.PollRepository
}

func NewPollService(repo ports.PollRepository) *PollService {
	return &PollService{repo: repo}
}

func (s *PollService) CreatePoll(poll domain.Poll) error {
	return s.repo.CreatePoll(poll)
}

func (s *PollService) Vote(vote domain.Vote) error {
	// Check if the poll exists before voting
	_, err := s.repo.GetPoll(vote.PollID)
	if err != nil {
		return fmt.Errorf("poll not found: %w", err)
	}
	return s.repo.Vote(vote)
}

func (s *PollService) GetResults(pollID string) (domain.PollResult, error) {
	return s.repo.GetResults(pollID)
}
