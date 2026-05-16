package main

import (
	"log"
	"net/http"

	"devin-router/internal/config"
	"devin-router/internal/devin"
	"devin-router/internal/github"
	"devin-router/internal/server"
	"devin-router/internal/tokenpool"
)

func main() {
	cfg := config.Load()

	pool := tokenpool.NewPool(cfg.DevinAPIKeys)
	devinClient := devin.NewClient(cfg.DevinAPIBase, pool)
	ghClient := github.NewClient(cfg.GitHubToken)

	srv := server.New(cfg, devinClient, ghClient, pool)

	log.Printf("listening on %s", cfg.ListenAddr)
	if err := http.ListenAndServe(cfg.ListenAddr, srv.Router()); err != nil {
		log.Fatal(err)
	}
}
