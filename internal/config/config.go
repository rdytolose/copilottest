package config

import (
	"flag"
	"os"
	"strings"
)

type Config struct {
	ListenAddr     string
	TargetRepo     string
	GitHubToken    string
	DevinAPIKeys   []string
	DevinAPIBase   string
	MaxPromptChars int
}

func Load() Config {
	var listen = flag.String("listen", ":8080", "listen address")
	var repo = flag.String("repo", "", "target repo (owner/name)")
	var maxChars = flag.Int("max-prompt-chars", 120000, "max chars to include from diff")

	flag.Parse()

	keysEnv := os.Getenv("DEVIN_API_KEYS")
	keys := splitAndTrim(keysEnv)

	cfg := Config{
		ListenAddr:     *listen,
		TargetRepo:     *repo,
		GitHubToken:    os.Getenv("GITHUB_TOKEN"),
		DevinAPIKeys:   keys,
		DevinAPIBase:   getEnvDefault("DEVIN_API_BASE", "https://api.devin.ai/v1"),
		MaxPromptChars: *maxChars,
	}

	return cfg
}

func splitAndTrim(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func getEnvDefault(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}
