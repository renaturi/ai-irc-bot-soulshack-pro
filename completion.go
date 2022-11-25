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

	stream, err := cc.AI.CreateChatCompletionStream(ctx, ai.ChatCompletionRequest{
		MaxTokens: cc.Session.Config.MaxTokens,
		Model:     cc.Personality.Model,
		Messages:  cc.Session.GetHistory(),
		Stream:    true,
	})

	if err != nil {
		senderror(err, channel)
		return
	}

	defer stream.Close()
	chunker := &Chunker{
		Size:    cc.Session.Config.Chunkmax,
		Last:    time.Now(),
		Timeout: cc.Session.Config.Chunkdelay,
		Buffer:  &bytes.Buffer{},
	}

	for {
		response, err := stream.Recv()
		if err != nil {
			send(chunker.Buffer.String(), channel)
			if !errors.Is(err, io.EOF) {
				senderror(err, channel)
			}
			return
		}
		if len(response.Choices) != 0 {
			chunker.Buffer.WriteString(response.Choices[0].Delta.Content)
		}
		for {
			if ready, chunk := chunker.Chunk(); ready {
				send(string(*chunk), channel)
			} else {
				break
			}
		}
	}
}

func senderror(err error, channel chan<- *string) {
	e := err.Error()
	channel <- &e
}

func send(chunk string, channel chan<- *string) {
	channel <- &chunk
}

type Chunker struct {
	Size    int
	Last    time.Time
	Buffer  *bytes.Buffer
	Timeout time.Duration
}

func (c *Chunker) Chunk() (bool, *[]byte) {

	end := c.Size
	if c.Buffer.Len() < 