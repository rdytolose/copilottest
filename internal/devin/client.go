package devin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"devin-router/internal/tokenpool"
)

type Client struct {
	baseURL string
	pool    *tokenpool.Pool
	http    *http.Client
}

func NewClient(baseURL string, pool *tokenpool.Pool) *Client {
	return &Client{
		baseURL: baseURL,
		pool:    pool,
		http:    &http.Client{Timeout: 60 * time.Second},
	}
}

type CreateSessionRequest struct {
	Prompt string `json:"prompt"`
	Title  string `json:"title,omitempty"`
}

type CreateSessionResponse struct {
	SessionID string `json:"session_id"`
	URL       string `json:"url"`
	IsNew     bool   `json:"is_new_session"`
}

func (c *Client) CreateSession(prompt, title string) (CreateSessionResponse, error) {
	token, err := c.pool.Acquire()
	if err != nil {
		return CreateSessionResponse{}, err
	}

	body := CreateSessionRequest{Prompt: prompt, Title: title}
	b, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/sessions", c.baseURL), bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer "+token.Key)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return CreateSessionResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 429 || resp.StatusCode == 403 {
		c.pool.MarkExhausted(token.Key)
		return c.CreateSession(prompt, title)
	}
	if resp.StatusCode >= 300 {
		return CreateSessionResponse{}, errors.New("devin create session failed: " + resp.Status)
	}

	var out CreateSessionResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return CreateSessionResponse{}, err
	}
	return out, nil
}
