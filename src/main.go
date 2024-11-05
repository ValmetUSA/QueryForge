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

	// Create a MultiLineEntry for output
	output := widget.NewMultiLineEntry()
	output.SetPlaceHolder("Response will appear here.")

	// Create a scroll container for the output to enable vertical scrolling
	scrollOutput := container.NewVScroll(output)
	scrollOutput.SetMinSize(fyne.NewSize(0, 100)) // Set a minimum size for the scroll area

	progress := widget.NewProgressBar()
	progress.Hide()

	askButton := widget.NewButton("Ask", func() {
		question := input.Text
		if question == "" {
			output.SetText("Please enter a question.")
			return
		}

		progress.Show()

		msg := Message{Role: "user", Content: question}
		req := Request{Model: "llama3.2", Stream: false, Messages: []Message{msg}}

		go func() {
			resp, err := talkToOllama(defaultOllamaURL, req)
			if err != nil {
				output.SetText(fmt.Sprintf("Error: %v", err))
			} else {
				output.SetText(resp.Message.Content)
			}
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
