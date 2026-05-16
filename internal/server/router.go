package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) Router() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/health", s.handleHealth).Methods("GET")
	r.HandleFunc("/api/analyze/pr", s.handleAnalyzePR).Methods("POST")
	r.HandleFunc("/webhook/github", s.handleGitHubWebhook).Methods("POST")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web")))

	return r
}
