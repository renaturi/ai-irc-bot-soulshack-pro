
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
