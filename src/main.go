package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Valmet QueryForge")
	w.Resize(fyne.NewSize(550, 500))

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

		go func() {
			progress.SetValue(0.5)
			Response, err := talkToOllama(question) // Pass question directly without using a pointer
			if err != nil {
				output.SetText(fmt.Sprintf("Error: %v", err))
			} else {
				output.SetText(Response) // Set text directly from the returned string
			}
			progress.SetValue(1.0)
			progress.Hide()
		}()
	})

	// About button to span the top of the window
	aboutButton := widget.NewButtonWithIcon("About", theme.InfoIcon(), func() {
		dialog.ShowInformation("About", "QueryForge by VII @ Valmet, Inc.\n\nA lightweight app for edge device RAG document searches.\n\nBuilt with ❤️ by Valmet USA - Atlanta, Georgia.", w)
	})

	// Folder picker for selecting a directory to run the RAG search within
	folderPicker := widget.NewButton("Select Folder (PDF FILES ONLY)", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			input.SetText(uri.String())
		}, w)
	})

	resetButton := widget.NewButton("Clear All", func() {
		input.SetText("")
		output.SetText("")
	})

	// TODO: Implement a toolbar with cut, copy, and paste actions
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentCutIcon(), func() {}),
		widget.NewToolbarAction(theme.ContentCopyIcon(), func() {}),
		widget.NewToolbarAction(theme.ContentPasteIcon(), func() {}),
	)

	content := container.NewVBox(
		image,
		input,
		container.NewCenter(
			container.NewHBox(
				folderPicker,
				askButton,
				resetButton,
			),
		),
		progress,
		container.NewCenter(
			container.NewHBox(
				toolbar,
			),
		),
		scrollOutput,
		aboutButton,
	)
	w.SetContent(content)
	w.ShowAndRun()
}
