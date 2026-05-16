package server

import (
	"net/http"
)

func (s *Server) Router() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/api/analyze/pr", s.handleAnalyzePR)
	mux.HandleFunc("/webhook/github", s.handleGitHubWebhook)

	mux.Handle("/", http.FileServer(http.Dir("./web")))

	return mux
}
