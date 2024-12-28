// Writen by VII @ Valmet, Inc.
// A lightweight app for edge device RAG document searches.
// Developed by Kenneth Alexander Jenkins, Valmet USA - Atlanta, Georgia: 2024.
// This file contains the main function for the Valmet QueryForge GUI application.
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
	// Create a new Fyne application
	a := app.New()
	w := a.NewWindow("Valmet QueryForge")
	w.Resize(fyne.NewSize(200, 500)) // Edit this line to change the window size: width x height (pixels)

	// Load the Valmet logo image from a static resource
	image := canvas.NewImageFromResource(resourceValmetlogosmallPng)
	image.FillMode = canvas.ImageFillOriginal
	image.Resize(fyne.NewSize(75, 25)) // Edit this line to change the logo image size (default is 75x25 pixels)

	// Create an Entry for user input
	input := widget.NewEntry()
	input.SetPlaceHolder("Type your question here.")

	// Create a MultiLineEntry for output with text wrapping enabled
	output := widget.NewMultiLineEntry()
	output.SetPlaceHolder("Response will appear here.")
	output.Wrapping = fyne.TextWrapWord // Allow text wrapping to wrap at word boundaries

	// Create a vertical scroll container for the output
	scrollOutput := container.NewVScroll(output)
	scrollOutput.SetMinSize(fyne.NewSize(380, 200)) // Sinimum size for the scroll area

	// Progress bar to show the AI query progress
	progress := widget.NewProgressBar()
	progress.Hide()

	// Ask button to query the AI
	askButton := widget.NewButton("Query the AI", func() {
		question := input.Text
		if question == "" {
			output.SetText("Please enter a question.")
			return
		}

		// Show the progress bar and set the value to 0
		progress.Show()
		progress.SetValue(0) // Start the progress bar at 0

		// Start a goroutine to query the AI and update the progress bar
		// Note: This allows the UI to remain responsive while the AI is processing the question,
		// thanks to multi-threading built into Go.
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
		dialog.ShowInformation("About", "QueryForge \n by VII @ Valmet, Inc.\n\nA lightweight app for edge device RAG document searches.\n\nBuilt with ❤️ by Valmet USA - Atlanta, Georgia.", w)
	})

	// Settings button with menu containing checkboxes - Not yet functional
	settingsButton := widget.NewButtonWithIcon("Settings", theme.SettingsIcon(), func() {

		// Select base conversational model for the AI - selected 1b by default
		pickBaseModel := widget.NewLabel("Base AI Model:")

		// Select the model from the dropdown
		selectModel := widget.NewSelect([]string{"qwen2.5:0.5b", "llama3.2:3b"}, func(selected string) {
			fmt.Println("Selected model:", selected)
			setOllamaModelName(selected)
		})

		// // Select the embedding model for the AI - selected 33m by default
		// pickEmbeddingModel := widget.NewLabel("Embedding Model:")
		// selectEmbeddingModel := widget.NewSelect([]string{"all-minilm:33m", "all-minilm:22m"}, func(selected string) {
		// 	fmt.Println("Selected embedding model:", selected)
		// 	setEmbeddingModelName(selected)
		// })

		// Function to set the AI model from user preferences
		settingsMenu := container.NewVBox(
			pickBaseModel,
			selectModel,
			// pickEmbeddingModel,
			// selectEmbeddingModel,
		)

		// Show the settings menu with the selected AI models
		//config.SetAIFromUserPrefs(selectModel.Selected, selectEmbeddingModel.Selected)
		dialog.ShowCustom("Settings", "Close", settingsMenu, w)
	})

	// Folder picker for selecting a directory to run the RAG search within
	folderPicker := widget.NewButton("Select Folder \n (PDF, TXT Formats Only)", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				fmt.Println("Error opening folder:", err)
			}
			if uri == nil {
				fmt.Println("No folder selected, RAG search will not be used.")
				return
			}

			// Start the chunking process for the RAG search
			fmt.Println("Selected folder:", uri.String())

			// Start a goroutine to scan the directory and merge the files
			go func() {
				// Show a dialog to inform the user that the files are being processed
				dialog.ShowInformation("Processing Files", "This may take a while - please wait...", w)

				// Call mergeFilesToTemp and handle the result
				tempFileLocation, err := mergeFilesToTemp(uri.Path())
				if err != nil {
					dialog.ShowError(err, w)
					return
				}

				// Set the temporary file location for the AI query
				setTempFileLocation(tempFileLocation)

				// Notify the user of success and provide the location of the temporary file
				dialog.ShowInformation("Files Processed", fmt.Sprintf("Files processed successfully."), w)

				// // Clean up the temporary file after use
				// err = deleteTempFile(tempFileLocation)
				// if err != nil {
				// 	dialog.ShowError(err, w)
				// } else {
				// 	fmt.Println("Temporary file deleted successfully.")
				// }
			}()
		}, w)
	})

	// Reset button to clear all text fields
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

	// This is the main content of the window
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
