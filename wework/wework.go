package wework

import (
	"net/http"
	"sync"
	"time"
)

type AgentClient struct {
	corpID  string
	secret  string
	agentID int
	httpc   *http.Client
	token   Token
}

type Token struct {
	mu        sync.Mutex
	token     string
	expiresAt time.Time
}

func NewAgentClient(corpID string, secret string, agentID int) *AgentClient {
	return &AgentClient{
		corpID:  corpID,
		secret:  secret,
		agentID: agentID,
	}
}
