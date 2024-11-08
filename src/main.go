package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Valmet QueryForge")
	w.Resize(fyne.NewSize(400, 300))

	image := canvas.NewImageFromFile("valmet_logo_small.png")
	image.FillMode = canvas.ImageFillOriginal
	image.Resize(fyne.NewSize(75, 25))

	input := widget.NewEntry()
	input.SetPlaceHolder("Type your question here...")

	// Create a MultiLineEntry for output with text wrapping enabled
	output := widget.NewMultiLineEntry()
	output.SetPlaceHolder("Response will appear here.")
	output.Wrapping = fyne.TextWrapWord // Set text wrapping to wrap by words

	// Create a vertical scroll container for the output
	scrollOutput := container.NewVScroll(output)
	scrollOutput.SetMinSize(fyne.NewSize(380, 200)) // Set a minimum size for the scroll area

	progress := widget.NewProgressBar()
	progress.Hide()

	askButton := widget.NewButton("Ask", func() {
		question := input.Text
		if question == "" {
			output.SetText("Please enter a question.")
			return
		}

		progress.Show()
		progress.SetValue(0) // Start the progress bar at 0

		msg := Message{Role: "user", Content: question}
		req := Request{Model: "llama3.2", Stream: false, Messages: []Message{msg}}

		go func() {
			progress.SetValue(0.5) // Update halfway through while waiting for a response
			resp, err := talkToOllama(defaultOllamaURL, req)
			if err != nil {
				output.SetText(fmt.Sprintf("Error: %v", err))
			} else {
				output.SetText(resp.Message.Content)
			}
			progress.SetValue(1.0) // Set progress to 100% when done
			progress.Hide()
		}()
	})

	content := container.NewVBox(
		image,
		input,
		askButton,
		progress,
		scrollOutput,
	)
	w.SetContent(content)
	w.ShowAndRun()
}
