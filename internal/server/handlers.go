package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"devin-router/internal/config"
	"devin-router/internal/devin"
	"devin-router/internal/github"
	"devin-router/internal/tokenpool"
)

type Server struct {
	cfg   config.Config
	devin *devin.Client
	gh    *github.Client
	pool  *tokenpool.Pool
}

func New(cfg config.Config, devinClient *devin.Client, ghClient *github.Client, pool *tokenpool.Pool) *Server {
	return &Server{cfg: cfg, devin: devinClient, gh: ghClient, pool: pool}
}

type AnalyzePRRequest struct {
	Repo        string `json:"repo"`
	PRNumber    int    `json:"pr_number"`
	PromptExtra string `json:"prompt_extra,omitempty"`
}

type AnalyzePRResponse struct {
	SessionID string `json:"session_id"`
	URL       string `json:"url"`
}

func (s *Server) handleAnalyzePR(w http.ResponseWriter, r *http.Request) {
	var req AnalyzePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	resp, err := s.analyzePRInternal(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) analyzePRInternal(req AnalyzePRRequest) (AnalyzePRResponse, error) {
	repo := req.Repo
	if repo == "" {
		repo = s.cfg.TargetRepo
	}
	if repo == "" {
		return AnalyzePRResponse{}, fmt.Errorf("repo is required")
	}
	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return AnalyzePRResponse{}, fmt.Errorf("repo must be owner/name")
	}

	diff, err := s.gh.FetchPRDiff(parts[0], parts[1], req.PRNumber)
	if err != nil {
		return AnalyzePRResponse{}, err
	}

	if len(diff) > s.cfg.MaxPromptChars {
		diff = diff[:s.cfg.MaxPromptChars] + "\n\n[truncated]"
	}

	prompt := fmt.Sprintf(
		"Review this GitHub PR diff. Provide: summary, risks, suggested tests, and any issues.\n\nPR: %s#%d\n\nDIFF:\n%s",
		repo, req.PRNumber, diff,
	)
	if req.PromptExtra != "" {
		prompt += "\n\nExtra instructions:\n" + req.PromptExtra
	}

	session, err := s.devin.CreateSession(prompt, fmt.Sprintf("PR Review %s#%d", repo, req.PRNumber))
	if err != nil {
		return AnalyzePRResponse{}, err
	}

	return AnalyzePRResponse{SessionID: session.SessionID, URL: session.URL}, nil
}
