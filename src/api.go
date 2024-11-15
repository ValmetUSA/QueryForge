package main

import (
	"context"
	"strings"

	"github.com/henomis/lingoose/llm/ollama"
	"github.com/henomis/lingoose/thread"
)

const defaultOllamaURL = "http://localhost:11434/api"
const ollamaModelName = "llama3.2"

func talkToOllama(request string) (string, error) {
	ollamaThread := thread.New()

	ollamaThread.AddMessage(thread.NewUserMessage().AddContent(
		thread.NewTextContent(request),
	))

	var responseBuilder strings.Builder
	err := ollama.New().WithEndpoint(defaultOllamaURL).WithModel(ollamaModelName).
		WithStream(func(s string) {
			//fmt.Print(s)
			responseBuilder.WriteString(s) // Append each streamed response chunk to the builder
		}).Generate(context.Background(), ollamaThread)

	if err != nil {
		return "", err
	}

	return responseBuilder.String(), nil // Return the full response as a string
}
