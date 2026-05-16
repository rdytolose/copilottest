package server

import (
	"encoding/json"
	"net/http"
)

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (s *Server) handleGitHubWebhook(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Action      string `json:"action"`
		PullRequest struct {
			Number int `json:"number"`
			Base   struct {
				Repo struct {
					FullName string `json:"full_name"`
				} `json:"repo"`
			} `json:"base"`
		} `json:"pull_request"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	if payload.Action == "opened" || payload.Action == "synchronize" || payload.Action == "reopened" {
		req := AnalyzePRRequest{
			Repo:     payload.PullRequest.Base.Repo.FullName,
			PRNumber: payload.PullRequest.Number,
		}
		go func() {
			_, _ = s.analyzePRInternal(req)
		}()
	}

	w.WriteHeader(http.StatusOK)
}
