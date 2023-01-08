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
		ctx.Reply(fmt.Sprintf("Usage: /set %s <value>", keysAsString(configParams)))
		return
	}

	param, v := ctx.Args[1], ctx.Args[2:]
	value := strings.Join(v, " ")
	if _, ok := configParams[param]; !ok {
		ctx.Reply(fmt.Sprintf("Unknown parameter. Supported parameters: %v", keysAsString(configParams)))
		return
	}

	// set on global config
	vip.Set(param, value)
	ctx.Reply(fmt.Sprintf("%s set to: %s", param, vip.GetString(param)))

	if param == "nick" {
		ctx.Client.Cmd.Nick(value)
	}

	ctx.Session.Reset()
}

func handleGet(ctx *ChatContext) {

	tokens := ctx.Args
	if len(tokens) < 2 {
		ctx.Reply(fmt.Sprintf("Usage: /get %s", keysAsString(configParams)))
		return
	}

	param := tokens[1]
	if _, ok := configParams[param]; !ok {
		ctx.Reply(fmt.Sprintf("Unknown parameter. Supported parameters: %v", keysAsString(configParams)))
		return
	}

	value := vip.GetString(param)
	ctx.Reply(fmt.Sprintf("%s: %s", param, value))
}

func handleSave(ctx *ChatContext) {

	tokens := ctx.Args
	if !ctx.IsAdmin() {
		ctx.Reply("You don't have permission to perform this action.")
		return
	}

	if len(tokens) < 2 {
		ctx.Reply("Usage: /save <name>")
		return
	}

	filename := tokens[1]

	v := vip.New()

	v.Set("nick", ctx.Personality.Nick)
	v.Set("prompt", ctx.Personality.Prompt)
	v.Set("model", ctx.Personality.Model)
	v.Set("greeting"