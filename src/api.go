package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const defaultOllamaURL = "http://localhost:11434/api/chat"

func talkToOllama(url string, ollamaReq Request) (*Response, error) {
	js, err := json.Marshal(&ollamaReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	client := http.Client{}
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(js))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer httpResp.Body.Close()

	var fullContent string
	decoder := json.NewDecoder(httpResp.Body)
	for {
		var part Response
		if err := decoder.Decode(&part); err != nil {
			break
		}
		fullContent += part.Message.Content

		if part.Done {
			break
		}
	}

	if fullContent == "" {
		return nil, fmt.Errorf("no response content received from server")
	}

	return &Response{
		Model:   ollamaReq.Model,
		Message: Message{Role: "assistant", Content: fullContent},
		Done:    true,
	}, nil
}
