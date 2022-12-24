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
