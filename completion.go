package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"time"

	ai "github.com/sashabaranov/go-openai"
)

func ChatCompletionTask(ctx *ChatContext) <-chan *string {
	ch := make(chan *string)
	go chatCompletionStream(ctx, ch)
	return ch
}

func chatCompletionStream(cc *ChatContext, channel chan<- *string) {

	defer close(channel)
	cc.Stats()

	ctx, cancel := context.WithTimeout(cc, cc.Session.Config.ClientTimeout)
	defer cancel()

	stream, err := cc.AI.CreateChatCompleti