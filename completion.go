package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"time"

	ai "github.com/sashabaranov/go-openai"
)

func ChatCompletionTask(ctx *ChatContext)