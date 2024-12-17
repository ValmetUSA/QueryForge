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

// Define system content and options for the query
const systemContent = `You are a helpful assistant by the name of PaperPal.
You were designed to help users with their queries using the information from the documents.
You were created by the Finnish company Valmet. You should be friendly and helpful to the users.
All answers should be based on the information from the documents, unless otherwise specified or inferred.
Documents will appear as previous messages in the conversation - you can refer to them directly if needed.
You should not make up any information. If you don't know the answer, you should say so. Do not hallucinate.
If queried about a topic without the needed to refer to the documents, you should answer based on your training data.
Make these answers as helpful as possible - and try to relate the reply back to Valmet (for paper and automation only).
If a document is not provided, you should refer to your own training data for the answer.
`

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

// Global variable to store the document content
var documentsContent string

// Setter function for the Ollama model name
func setOllamaModelName(ollamaModelName string) {
	config.ollamaModelName = ollamaModelName
}

// Setter function for the embedding model name
func setEmbeddingModelName(embeddingModelName string) {
	config.embeddingModelName = embeddingModelName
}

// Function to talk to Ollama
func talkToOllama(userContent string, tempFileLocation string) (string, error) {
	var err error
	if tempFileLocation == "" {
		fmt.Printf("No file to process. Proceeding with query.\n")
	} else {
		// Create a store to save the embeddings
		// Keep it in this function to let the garbage collector clean it up
		store := embeddings.MemoryVectorStore{
			Records: make(map[string]llm.VectorRecord),
		}

		// Read the file content
		fmt.Printf("This is the file location: %s\n", tempFileLocation)
		rulesContent, err := content.ReadTextFile(tempFileLocation)
		if err != nil {
			log.Printf("❌1: Failed to read document: %v", err)
			return "", fmt.Errorf("failed to read document: %w", err)
		}

		// Chunk the content into manageable pieces
		chunks := content.ChunkText(rulesContent, 500, 200)
		fmt.Printf("Chunking complete. Number of chunks: %d\n", len(chunks))

		// Create embeddings from documents and save them in the store
		for idx, doc := range chunks {
			fmt.Printf("Creating embedding from document %d...\n", idx)
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
		fmt.Println("🔎 Searching for similarity...")

		similarities, _ := store.SearchSimilarities(embeddingFromQuestion, 0.4)

		fmt.Printf("🎉 Similarities found: %d\n", len(similarities))
	}

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature:   0.7,
		option.RepeatLastN:   2,
		option.RepeatPenalty: 2.0,
		option.TopK:          10,
		option.TopP:          0.5,
	})

	// If the file is not provided, use the static content
	if tempFileLocation == "" {
		documentsContent = "This is the static fallback content."
	}

	// fmt.Printf("Document Content:\n%s\n", documentsContent)

	var query llm.Query
	// Declare the query variable outside the if block so it is in scope for the function
	if tempFileLocation == "" {
		query = llm.Query{
			Model: config.ollamaModelName,
			Messages: []llm.Message{
				{Role: "system", Content: systemContent},
				{Role: "system", Content: documentsContent},
				{Role: "user", Content: userContent},
			},
			Options: options,
		}
	} else {
		query = llm.Query{
			Model: config.ollamaModelName,
			Messages: []llm.Message{
				{Role: "system", Content: systemContent},
				{Role: "user", Content: userContent},
			},
			Options: options,
		}
	}

	fmt.Println("🤖 Answering query...")

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

	// Clean up the temporary file after use (if applicable)
	// err = deleteTempFile(tempFileLocation)
	// if err != nil {
	// 	dialog.ShowError(err, w)
	// } else {
	// 	fmt.Println("Temporary file deleted successfully.")
	// }

	return responseBuilder.String(), nil
}
