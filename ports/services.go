package ports

import "polling-system/domain"

type PollService interface {
	CreatePoll(poll domain.Poll) error
	Vote(vote domain.Vote) error
	GetResults(pollID string) (domain.PollResult, error)
	VoteMultiple(votes domain.MultiVote) error
}
