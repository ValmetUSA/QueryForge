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

// Sauced from here: https://github.com/parakeet-nest/parakeet/tree/main/examples/23-rag-with-chunker

const defaultOllamaURL = "http://localhost:11434"
const ollamaModelName = "llama3.2:3b"
const embeddingModelName = "all-minilm:33m"

func talkToOllama(userContent string) (string, error) {

	store := embeddings.MemoryVectorStore{
		Records: make(map[string]llm.VectorRecord),
	}

	rulesContent, err := content.ReadTextFile("./test_doc.txt")
	if err != nil {
		log.Fatalln("😡:", err)
	}
	chunks := content.ChunkText(rulesContent, 500, 200)

	// Create embeddings from documents and save them in the store
	for idx, doc := range chunks {
		fmt.Println("Creating embedding from document ", idx)
		embedding, err := embeddings.CreateEmbedding(
			defaultOllamaURL,
			llm.Query4Embedding{
				Model:  embeddingModelName,
				Prompt: doc,
			},
			strconv.Itoa(idx),
		)
		if err != nil {
			fmt.Println("😡:", err)
		} else {
			store.Save(embedding)
		}
	}

	systemContent := `You are a helpful assistant by the name of PaperPal.
	You were designed to help users with their queries using the information from the documents.
	You were created by the Finnish company Valmet Oyj. You should be friendly and helpful to the users.
	All answers should be based on the information from the documents, unless otherwise specified or inferred.
	Documents will appear as previous messages in the conversation - you can refer to them directly if needed.
	You should not make up any information. If you don't know the answer, you should say so.
	`

	// Create an embedding from the question
	embeddingFromQuestion, err := embeddings.CreateEmbedding(
		defaultOllamaURL,
		llm.Query4Embedding{
			Model:  embeddingModelName,
			Prompt: userContent,
		},
		"question",
	)
	if err != nil {
		log.Fatalln("😡:", err)
	}
	fmt.Println("🔎 searching for similarity...")

	similarities, _ := store.SearchSimilarities(embeddingFromQuestion, 0.3)

	fmt.Println("🎉 similarities:", len(similarities))

	documentsContent := embeddings.GenerateContentFromSimilarities(similarities)

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature:   0.0,
		option.RepeatLastN:   2,
		option.RepeatPenalty: 2.0,
		option.TopK:          10,
		option.TopP:          0.5,
	})

	query := llm.Query{
		Model: ollamaModelName,
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
		log.Fatal("😡:", err)
	} else {
		return responseBuilder.String(), nil
	}

	return "", nil
}
