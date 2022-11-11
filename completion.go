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

func chatCompletionStr