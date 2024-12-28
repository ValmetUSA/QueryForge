package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/ollama/ollama/api"
)

// Define constants for the query options
var (
	FALSE = false
	TRUE  = true
)

var ollamaModelName = "qwen2.5:0.5b"

// var tempFileLocation string // Global variable to store the temporary file location

func setOllamaModelName(modelName string) {
	ollamaModelName = modelName
}

func getOllamaModelName() string {
	return ollamaModelName
}

// Define system content and options for the query
const systemInstructions = `You are a helpful assistant by the name of PaperPal.
Your purpose is to assist users with questions, mostly related to paper and automation.
You were created by the Finnish company Valmet, a lead developer and supplier of process 
technologies, automation systems and services for the pulp, paper, energy industries.

You should be friendly and helpful to the users. All answers should be based on the information from the documents, 
unless otherwise specified or inferred. Documents will appear as previous messages in the 
conversation - you can refer to them directly if needed. You should not make up any information. If you don't 
know the answer, you should say so. Do not hallucinate.

If queried about a topic without the needed to refer to the documents, you should answer based on your training data.
Make these answers as helpful as possible - and try to relate the reply back to Valmet (for paper and automation only).
`

func talkToOllama(userQuestion string) (string, error) {
	ctx := context.Background()

	// Set the Ollama host
	ollamaRawUrl := os.Getenv("OLLAMA_HOST")
	if ollamaRawUrl == "" {
		ollamaRawUrl = "http://localhost:11434"
	}

	parsedUrl, _ := url.Parse(ollamaRawUrl)
	client := api.NewClient(parsedUrl, http.DefaultClient)

	var messages []api.Message

	// Read context from the temporary file if provided
	var context []byte
	if tempFileLocation != "" {
		fmt.Println("DEBUG: Attempting to read file at:", tempFileLocation)
		var err error
		context, err = os.ReadFile(tempFileLocation)
		if err != nil {
			fmt.Printf("DEBUG: Error reading temp file (%s): %v\n", tempFileLocation, err)
			context = []byte("Error reading context file, defaulting to empty content.")
		}
	} else {
		fmt.Println("DEBUG: No temp file location provided, using default content.")
		context = []byte("No context file provided.")
	}

	// Log the content read for debugging
	fmt.Printf("DEBUG: Content from temp file:\n%s\n", string(context))

	// Prepare the messages for the API request
	messages = []api.Message{
		{Role: "system", Content: systemInstructions},
		{Role: "system", Content: "CONTENT:\n" + string(context)},
		{Role: "user", Content: userQuestion},
	}

	// Configure the chat request
	req := &api.ChatRequest{
		Model:    getOllamaModelName(),
		Messages: messages,
		Options: map[string]interface{}{
			"temperature":    0.4,
			"repeat_last_n":  2,
			"repeat_penalty": 1.8,
			"top_k":          10,
			"top_p":          0.5,
		},
		Stream: &TRUE,
	}

	// Capture response
	responseBuilder := &strings.Builder{}

	err := client.Chat(ctx, req, func(resp api.ChatResponse) error {
		fmt.Print(resp.Message.Content)
		responseBuilder.WriteString(resp.Message.Content)
		return nil
	})

	// Handle errors gracefully
	if err != nil {
		log.Printf("Error during chat request: %v\n", err)
		return "", err
	}

	return responseBuilder.String(), nil
}
