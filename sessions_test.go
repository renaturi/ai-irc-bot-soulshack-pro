
package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"

	ai "github.com/sashabaranov/go-openai"
	vip "github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestChatSession(t *testing.T) {
	vip.Set("session", 1*time.Hour)
	vip.Set("history", 10)

	//log.SetOutput(io.Discard)

	ctx := &ChatContext{
		Personality: &Personality{
			Prompt: "You are a helpful assistant.",
		},
	}

	t.Run("Test interactions and message history", func(t *testing.T) {
		session1 := sessions.Get("session1")
		session1.Message(ctx, ai.ChatMessageRoleUser, "Hello!")
		session1.Message(ctx, ai.ChatMessageRoleAssistant, "Hi there!")

		assert.Len(t, session1.History, 3)
		assert.Equal(t, session1.History[1].Content, "Hello!")
		assert.Equal(t, session1.History[2].Content, "Hi there!")
	})
}
func TestExpiry(t *testing.T) {
	//log.SetOutput(io.Discard)
	ctx := &ChatContext{
		Personality: &Personality{
			Prompt: "You are a helpful assistant.",
		},
	}
	t.Run("Test session expiration and trimming", func(t *testing.T) {
		vip.Set("session", 500*time.Millisecond)
		vip.Set("history", 20)

		session2 := sessions.Get("session2")
		session2.Message(ctx, ai.ChatMessageRoleUser, "How are you?")
		session2.Message(ctx, ai.ChatMessageRoleAssistant, "I'm doing great, thanks!")
		session2.Message(ctx, ai.ChatMessageRoleUser, "What's your name?")

		time.Sleep(2 * time.Second)
		session3 := sessions.Get("session2")

		assert.NotEqual(t, session2, session3, "Expired session should not be reused")
		assert.Len(t, session3.History, 0, "New session history should be empty")

		session3.Message(ctx, ai.ChatMessageRoleUser, "Hello again!")
		session3.Message(ctx, ai.ChatMessageRoleAssistant, "Hi! Nice to see you again!")

		assert.Len(t, session3.History, 3, "History should include the latest 2 messages plus the initial system message")
		assert.Equal(t, session3.History[1].Content, "Hello again!")
		assert.Equal(t, session3.History[2].Content, "Hi! Nice to see you again!")
	})
}

func TestSessionConcurrency(t *testing.T) {
	vip.Set("session", 1*time.Hour)
	vip.Set("history", 10)

	log.SetOutput(io.Discard)

	t.Run("Test session concurrency", func(t *testing.T) {
		vip.Set("session", 1*time.Hour)
		vip.Set("history", 500*2000)

		ctx := &ChatContext{
			Personality: &Personality{
				Prompt: "You are a helpful assistant.",
			},
		}

		const concurrentUsers = 1000
		const messagesPerUser = 500

		startTime := time.Now()

		var wg sync.WaitGroup
		wg.Add(concurrentUsers)

		for i := 0; i < concurrentUsers; i++ {
			go func(userIndex int) {
				defer wg.Done()
				sessionID := fmt.Sprintf("usersession%d", userIndex)
				session := sessions.Get(sessionID)

				for j := 0; j < messagesPerUser; j++ {
					session.Message(ctx, ai.ChatMessageRoleUser, fmt.Sprintf("User %d message %d", userIndex, j))
					session.Message(ctx, ai.ChatMessageRoleAssistant, fmt.Sprintf("Assistant response to user %d message %d", userIndex, j))
				}
			}(i)
		}

		wg.Wait()

		for i := 0; i < concurrentUsers; i++ {
			sessionID := fmt.Sprintf("usersession%d", i)
			session := sessions.Get(sessionID)
			assert.Len(t, session.History, messagesPerUser*2+1, "Each session should have the correct number of messages")
		}
		elapsedTime := time.Since(startTime)
		totalMessages := concurrentUsers * messagesPerUser * 2
		messagesPerSecond := float64(totalMessages) / elapsedTime.Seconds()
		t.Logf("Processed %d messages in %v, which is %.2f messages per second\n", totalMessages, elapsedTime, messagesPerSecond)
	})
}

func TestSingleSessionConcurrency(t *testing.T) {
	log.SetOutput(io.Discard)

	t.Run("Test single session concurrency", func(t *testing.T) {
		vip.Set("session", 1*time.Hour)
		vip.Set("history", 500*200)

		ctx := &ChatContext{
			Personality: &Personality{
				Prompt: "You are a helpful assistant.",
			},
		}

		const concurrentUsers = 500
		const messagesPerUser = 100

		startTime := time.Now()

		session := sessions.Get("concurrentSession")

		var wg sync.WaitGroup
		wg.Add(concurrentUsers)

		for i := 0; i < concurrentUsers; i++ {
			go func(userIndex int) {
				defer wg.Done()
				for j := 0; j < messagesPerUser; j++ {
					session.Message(ctx, ai.ChatMessageRoleUser, fmt.Sprintf("User %d message %d", userIndex, j))
					session.Message(ctx, ai.ChatMessageRoleAssistant, fmt.Sprintf("Assistant response to user %d message %d", userIndex, j))
				}
			}(i)
		}

		wg.Wait()

		elapsedTime := time.Since(startTime)
		totalMessages := concurrentUsers * messagesPerUser * 2
		messagesPerSecond := float64(totalMessages) / elapsedTime.Seconds()

		assert.Len(t, session.History, totalMessages+1, "The session should have the correct number of messages")
		t.Logf("Processed %d messages in %v, which is %.2f messages per second\n", totalMessages, elapsedTime, messagesPerSecond)