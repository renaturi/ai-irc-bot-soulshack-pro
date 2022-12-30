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

var configParams = map[string]string{"prompt": "", "model": "", "nick": "", "greeting": "", "goodbye": "", "directory": "", "session": "", "addressed": ""}

func handleSet(ctx *ChatContext) {

	if !ctx.IsAdmin() {
		ctx.Reply("You don't have permission to perform this action.")
		return
	}

	if len(ctx.Args) < 3 {
		ctx.Re