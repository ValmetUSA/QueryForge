package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas" // Import canvas for image handling
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Response struct {
	Model              string    `json:"model"`
	CreatedAt          time.Time `json:"created_at"`
	Message            Message   `json:"message"`
	Done               bool      `json:"done"`
	TotalDuration      int64     `json:"total_duration"`
	LoadDuration       int       `json:"load_duration"`
	PromptEvalCount    int       `json:"prompt_eval_count"`
	PromptEvalDuration int       `json:"prompt_eval_duration"`
	EvalCount          int       `json:"eval_count"`
	EvalDuration       int64     `json:"eval_duration"`
}

const defaultOllamaURL = "http://localhost:11434/api/chat"

func main() {
	// Initialize the GUI app
	a := app.New()
	w := a.NewWindow("Valmet QueryForge")
	w.Resize(fyne.NewSize(400, 300))

	// Load the PNG image and set its size
	image := canvas.NewImageFromFile("valmet_logo_small.png") // Update with your image path
	image.FillMode = canvas.ImageFillOriginal                 // Adjust fill mode as needed
	image.Resize(fyne.NewSize(75, 25))                        // Lock the image size to 75 x 25 pixels

	// Input and output widgets
	input := widget.NewEntry()
	input.SetPlaceHolder("Type your question here...")
	output := widget.NewLabel("Response will appear here.")
	output.Wrapping = fyne.TextWrapWord // Enable word wrapping
	progress := widget.NewProgressBar() // Create a new progress bar
	progress.Hide()                     // Initially hide the progress bar

	// Ask button with functionality
	askButton := widget.NewButton("Ask", func() {
		question := input.Text
		if question == "" {
			output.SetText("Please enter a question.")
			return
		}

		// Show the progress bar
		progress.Show()

		// Create request and send it
		msg := Message{Role: "user", Content: question}
		req := Request{Model: "llama3.2", Stream: false, Messages: []Message{msg}}

		go func() {
			resp, err := talkToOllama(defaultOllamaURL, req)
			if err != nil {
				output.SetText(fmt.Sprintf("Error: %v", err))
			} else {
				// Display response in the output label
				output.SetText(resp.Message.Content)
			}

			// Hide the progress bar after response
			progress.Hide()
		}()
	})

	// Layout for the window
	content := container.NewVBox(
		image, // Add the image widget to the layout
		input,
		askButton,
		progress,
		output,
	)
	w.SetContent(content)
	w.ShowAndRun()
}

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

	// If no content was gathered, return an error
	if fullContent == "" {
		return nil, fmt.Errorf("no response content received from server")
	}

	// Prepare the final response
	return &Response{
		Model:   ollamaReq.Model,
		Message: Message{Role: "assistant", Content: fullContent},
		Done:    true,
	}, nil
}
