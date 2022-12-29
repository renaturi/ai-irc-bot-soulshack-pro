package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	ai "github.com/sashabaranov/go-openai"
	vip "github.com/spf13/viper"
)

func sendGreeting(ctx *ChatContext) {
	log.Println("sending greeting...")
	ctx.Session.Message(ctx, ai.ChatMessageRoleAssistant, ctx.Personality.Greeting)
	rch := ChatCompletionTask(ctx)
	_ = spoolFromChannel(ctx, rch)
	ctx.Session.Reset()
}

func spoolFromChannel(ctx *ChatContext, msgch <-chan *string) *string {
	all := strings.Builder{}
	for reply := range msgch {
		all.WriteString(*reply)
		sendMessage(ctx, reply)
	}
	s := all.String()
	return &s
}

func sendMessage(ctx *ChatContext, message *string) {
	log.Println("<<", ctx.Personality.Nick, *message)
	ctx.Reply(*message)
}

var configParams = map[string]string{"prompt": "", "model": "", "nick": "", "greeting": "", "goodbye