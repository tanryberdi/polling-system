package domain

type Poll struct {
	ID       string
	Question string
	Options  []string
}

type Vote struct {
	PollID string `json:"poll_id"`
	Option string `json:"option"`
}

type PollResult struct {
	Poll    Poll
	Results map[string]int
}
