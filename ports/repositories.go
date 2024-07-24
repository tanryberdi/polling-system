package ports

import "polling-system/domain"

type PollRepository interface {
	CreatePoll(poll domain.Poll) error
	GetPoll(id string) (domain.Poll, error)
	Vote(vote domain.Vote) error
	GetResults(pollID string) (domain.PollResult, error)
}
