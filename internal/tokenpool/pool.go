package tokenpool

import (
	"errors"
	"sync"
)

type Token struct {
	Key       string
	Exhausted bool
}

type Pool struct {
	mu     sync.Mutex
	tokens []Token
	next   int
}

func NewPool(keys []string) *Pool {
	tokens := make([]Token, 0, len(keys))
	for _, k := range keys {
		tokens = append(tokens, Token{Key: k})
	}
	return &Pool{tokens: tokens}
}

func (p *Pool) Acquire() (Token, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.tokens) == 0 {
		return Token{}, errors.New("no tokens configured")
	}

	for i := 0; i < len(p.tokens); i++ {
		idx := (p.next + i) % len(p.tokens)
		if !p.tokens[idx].Exhausted {
			p.next = idx + 1
			return p.tokens[idx], nil
		}
	}
	return Token{}, errors.New("all tokens exhausted")
}

func (p *Pool) MarkExhausted(key string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for i := range p.tokens {
		if p.tokens[i].Key == key {
			p.tokens[i].Exhausted = true
			return
		}
	}
}
