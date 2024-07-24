package domain

type Poll struct {
	ID       string
	Question string
	Options  []string
}

type Vote struct {
	PollID string
	Option string
}

type PollResult struct {
	Poll    Poll
	Results map[string]int
}
