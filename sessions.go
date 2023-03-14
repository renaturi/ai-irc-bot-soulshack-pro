
package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	ai "github.com/sashabaranov/go-openai"
	vip "github.com/spf13/viper"
)

var sessions = Chats{
	sessionMap: make(map[string]*ChatSession),
	mu:         sync.RWMutex{},
}

type Chats struct {
	sessionMap map[string]*ChatSession
	mu         sync.RWMutex
}

type SessionConfig struct {
	MaxTokens      int
	SessionTimeout time.Duration
	MaxHistory     int
	ClientTimeout  time.Duration
	Chunkdelay     time.Duration
	Chunkmax       int
}

type ChatSession struct {
	Config     SessionConfig
	Name       string
	History    []ai.ChatCompletionMessage
	mu         sync.RWMutex
	Last       time.Time
	Totalchars int
}

func (s *ChatSession) GetHistory() []ai.ChatCompletionMessage {
	s.mu.RLock()
	defer s.mu.RUnlock()

	historyCopy := make([]ai.ChatCompletionMessage, len(s.History))
	copy(historyCopy, s.History)

	return historyCopy
}

func (s *ChatSession) Message(ctx *ChatContext, role string, message string) *ChatSession {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.History) == 0 {
		s.History = append(s.History, ai.ChatCompletionMessage{Role: ai.ChatMessageRoleSystem, Content: ctx.Personality.Prompt})
		s.Totalchars += len(ctx.Personality.Prompt)
	}

	s.History = append(s.History, ai.ChatCompletionMessage{Role: role, Content: message})
	s.Totalchars += len(message)
	s.Last = time.Now()

	s.trim()

	return s
}