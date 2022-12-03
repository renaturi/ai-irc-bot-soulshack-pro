
package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestChunker_Chunk(t *testing.T) {
	timeout := 1000 * time.Millisecond

	tests := []struct {
		name     string
		input    string
		size     int
		expected []byte
	}{
		{
			name:     "chunk on newline",
			input:    "Hello\nworld",
			size:     350,
			expected: []byte("Hello\n"),
		},
		{
			name:     "chunk on buffer size",
			input:    "Hello",
			size:     5,
			expected: []byte("Hello"),
		},
		{
			name:     "no chunk",
			input:    "Hello",
			size:     10,
			expected: []byte(""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Chunker{
				Size:    tt.size,
				Last:    time.Now(),
				Buffer:  &bytes.Buffer{},
				Timeout: timeout,
			}
			c.Buffer.WriteString(tt.input)

			chunked, chunk := c.Chunk()
			if chunked && string(*chunk) != string(tt.expected) {
				t.Errorf("Chunk() got = %v, want = %v", chunk, tt.expected)
			}
		})
	}