# The Scenario: "Decision by the Minute"

## Problem Statement
In the bustling town of Codeville, a lively debate club named "The Decisives" meets every week. Their discussions range from the trivial, like "Pineapple on pizza: Yes or No?", to the more serious, like "The future of remote work." However, they've encountered a challenge: tallying votes quickly and efficiently during their spirited debates has become cumbersome, leading to delays and diminished engagement.
Seeing an opportunity, they've turned to you. Known for your knack for solving practical problems with technology, you're tasked with creating a solution that allows members to cast their votes on various topics in real-time and see the results unfold as the debate rages on. Your system will not only decide the fate of pineapple on pizza but also serve as a model for real-time decision-making in communities far and wide.

## Your Mission: Build a Simple Real-time Polling System
### Objective
Develop a RESTful service in Go that enables "The Decisives" to create polls, vote in real-time, and view poll results as they come in, keeping the debate lively and engaging.

### Task Overview:
1. Create a new poll: acilitate the creation of new polls with a question and multiple choice answers.
2. Vote on a poll: Allow members to vote on active polls, with the capability to handle a flurry of votes simultaneously.
3. Live results: Implement a real-time update mechanism, so members can watch the voting results change moment by moment.

### Questions
1. What were some of the key assumptions you made, and what trade-offs did you encounter?
2. Imagining this project were to evolve into a full-scale real-world application, what enhancements or next steps would you prioritize to elevate its functionality, user experience, and technical robustness?

## For running and testing the code locally
### Installation
1. Clone the repository
```bash
git clone https://github.com/tanryberdi/polling-system.git
```

2. Path to the project directory
```bash
cd polling-system
```

### Running the code
The code is written in Go. To run the code, please use the following command:
```bash
make run
```

### Running the tests
```bash
make test
```

### Running the linter
```bash
make lint
```

### API Endpoint - For creating poll
```curl
curl -X POST http://localhost:8080/create_poll -d '{"id":"1", "question":"Pineapple on pizza?", "options":["Yes", "No"]}'
```

### API Endpoint - For vote
```curl
curl -X POST http://localhost:8080/vote -d '{"poll_id":"1", "option":"No"}'
```

### API Endpoint - For a multiple vote
```curl
curl -X POST http://localhost:8080/vote_multiple -H "Content-Type: application/json" -d '[
  {"poll_id":"1", "option":"Yes"},
  {"poll_id":"2", "option":"No"}
]'
```

### API Endpoint - For getting results
```curl
curl http://localhost:8080/results/1
```

### API Endpoint - For a live results
```curl
curl http://localhost:8080/poll_updates/1
```
