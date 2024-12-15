package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/content"
	"github.com/parakeet-nest/parakeet/embeddings"
	"github.com/parakeet-nest/parakeet/enums/option"
	"github.com/parakeet-nest/parakeet/llm"
)

// Some of this was sauced from here: https://github.com/parakeet-nest/parakeet/tree/main/examples/23-rag-with-chunker

// Constant declarations
const defaultOllamaURL = "http://localhost:11434"

// Config holds the AI model configurations
type Config struct {
	ollamaModelName    string
	embeddingModelName string
}

// Default configuration
var config = Config{
	ollamaModelName:    "llama3.2:1b",
	embeddingModelName: "all-minilm:33m",
}

// Setter function for the Ollama model name
func setOllamaModelName(ollamaModelName string) {
	config.ollamaModelName = ollamaModelName
}

// Setter function for the embedding model name
func setEmbeddingModelName(embeddingModelName string) {
	config.embeddingModelName = embeddingModelName
}

// Function to talk to Ollama
func talkToOllama(userContent string) (string, error) {

	// Create a store to save the embeddings
	// Keep it in this function to let the garbage collector clean it up
	// our users will likely only use this function once to get their answer
	store := embeddings.MemoryVectorStore{
		Records: make(map[string]llm.VectorRecord),
	}

	// Read from here how to do that: https://parakeet-nest.github.io/parakeet/embeddings/
	rulesContent, err := content.ReadTextFile("./test_doc.txt")
	if err != nil {
		log.Printf("❌1: Failed to read document: %v", err)
		return "", fmt.Errorf("failed to read document: %w", err)
	}
	chunks := content.ChunkText(rulesContent, 500, 200)

	// Create embeddings from documents and save them in the store
	for idx, doc := range chunks {
		fmt.Println("Creating embedding from document ", idx)
		embedding, err := embeddings.CreateEmbedding(
			defaultOllamaURL,
			llm.Query4Embedding{
				Model:  config.embeddingModelName,
				Prompt: doc,
			},
			strconv.Itoa(idx),
		)
		if err != nil {
			log.Printf("❌2: Failed to create embedding for document %d: %v", idx, err)
			continue
		}
		store.Save(embedding)
	}

	systemContent := `You are a helpful assistant by the name of PaperPal.
	You were designed to help users with their queries using the information from the documents.
	You were created by the Finnish company Valmet Oyj. You should be friendly and helpful to the users.
	All answers should be based on the information from the documents, unless otherwise specified or inferred.
	Documents will appear as previous messages in the conversation - you can refer to them directly if needed.
	You should not make up any information. If you don't know the answer, you should say so. Do not halucinate.
	If queried about a topic without the needed to refer to the documents, you should answer based on your training data.
	Make these answers as helpful as possible - and try to relate the reply back to valmet (for paper and automation only).
	`

	// Create an embedding from the question
	embeddingFromQuestion, err := embeddings.CreateEmbedding(
		defaultOllamaURL,
		llm.Query4Embedding{
			Model:  config.embeddingModelName,
			Prompt: userContent,
		},
		"question",
	)
	if err != nil {
		log.Printf("❌3: Failed to create embedding from question, is Ollama on? - %v", err)
		return "", fmt.Errorf("failed to create embedding from question, is Ollama on? 🦙 \n\nError Data: %w", err)
	}
	fmt.Println("🔎 searching for similarity...")

	similarities, _ := store.SearchSimilarities(embeddingFromQuestion, 0.4)

	fmt.Println("🎉 similarities:", len(similarities))

	documentsContent := embeddings.GenerateContentFromSimilarities(similarities)

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature:   0.7,
		option.RepeatLastN:   2,
		option.RepeatPenalty: 2.0,
		option.TopK:          10,
		option.TopP:          0.5,
	})

	query := llm.Query{
		Model: config.ollamaModelName,
		Messages: []llm.Message{
			{Role: "system", Content: systemContent},
			{Role: "system", Content: documentsContent},
			{Role: "user", Content: userContent},
		},
		Options: options,
	}

	fmt.Println("🤖 answer:")

	// Answer the question
	var responseBuilder strings.Builder
	_, err = completion.ChatStream(defaultOllamaURL, query,
		func(answer llm.Answer) error {
			fmt.Print(answer.Message.Content)
			responseBuilder.WriteString(answer.Message.Content)
			return nil
		})

	if err != nil {
		log.Printf("❌4: Failed to get response: %v", err)
		return responseBuilder.String(), fmt.Errorf("failed to get response: %w", err)
	}

	return responseBuilder.String(), nil
}
