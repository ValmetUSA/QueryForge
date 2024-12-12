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
	w.Resize(fyne.NewSize(200, 500))

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

	askButton := widget.NewButton("Query the AI", func() {
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

	// Settings button with menu containing checkboxes - Not yet functional
	settingsButton := widget.NewButtonWithIcon("Settings", theme.SettingsIcon(), func() {

		// Select base conversational model for the AI - select 3b by default
		pickBaseModel := widget.NewLabel("Base AI Model:")
		selectModel := widget.NewSelect([]string{"llama3.2:1b", "llama3.2:3b"}, func(selected string) {
			//fmt.Println("Selected model:", selected)
		})
		selectModel.Selected = "llama3.2:3b" // Set the default selection to "llama3.2:3b"

		// Select embedding model for the AI - select 33m by default
		pickEmbeddingModel := widget.NewLabel("Embedding Model:")
		selectEmbeddingModel := widget.NewSelect([]string{"all-minilm:33m", "all-minilm:125m"}, func(selected string) {
			//fmt.Println("Selected embedding model:", selected)
		})
		selectEmbeddingModel.Selected = "all-minilm:33m" // Set the default selection to "all-minilm:33m"

		settingsMenu := container.NewVBox(
			pickBaseModel,
			selectModel,
			pickEmbeddingModel,
			selectEmbeddingModel,
		)

		setAIfromUserPrefs(selectModel.Selected, selectEmbeddingModel.Selected)
		dialog.ShowCustom("Settings", "Close", settingsMenu, w)
	})

	// Folder picker for selecting a directory to run the RAG search within
	folderPicker := widget.NewButton("Select Folder", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				fmt.Println("Error opening folder:", err)
			}
			if uri == nil {
				fmt.Println("No folder selected")
				return
			}

			// Start the chunking process for the RAG search - TODO: Implement chunking calls
			input.SetText(uri.String()) // Remove this line when the chunking process is implemented
		}, w)
	})

	resetButton := widget.NewButton("Clear All", func() {
		input.SetText("")
		output.SetText("")
	})

	// Toolbar with copy, and paste actions
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentCopyIcon(), func() {
			clipboard := w.Clipboard()
			clipboard.SetContent(output.Text) // Save text to clipboard
		}),
		widget.NewToolbarAction(theme.ContentPasteIcon(), func() {
			clipboard := w.Clipboard()
			input.SetText(input.Text + clipboard.Content()) // Append clipboard text to entry
		}),
	)

	content := container.NewVBox(
		image,
		input,
		container.NewCenter(
			container.NewHBox(
				folderPicker,
				askButton,
			),
		),
		progress,
		container.NewCenter(
			container.NewHBox(
				toolbar,
				resetButton,
			),
		),
		scrollOutput,
		container.NewCenter(
			container.NewHBox(
				aboutButton,
				settingsButton,
			),
		),
	)
	w.SetContent(content)
	w.ShowAndRun()
}
