package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ledongthuc/pdf"
)

var tempFileLocation = "" // Global variable to store the temporary file location

// mergeFilesToTemp reads all files from the selected directory and outputs their
// combined content into a single temporary file, and returns its location.
func mergeFilesToTemp(dir string) (string, error) {
	tempDir := os.TempDir()
	tempFile, err := os.CreateTemp(tempDir, "merged_*.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}

	fmt.Printf("Temporary file created at: %s\n", tempFile.Name())

	fileCount := 0
	// Traverse the directory and process each file
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and specific unwanted files like .DS_Store
		if info.IsDir() || strings.EqualFold(info.Name(), ".ds_store") {
			return nil
		}

		fileCount++
		if fileCount > 50 {
			return fmt.Errorf("too many files in directory: %d", fileCount)
		}

		// Append valid files to the temp file
		if err := appendFileContents(tempFile, path); err != nil {
			return fmt.Errorf("failed to process file %s: %w", path, err)
		}
		return nil
	})
	if err != nil {
		_ = tempFile.Close()
		_ = os.Remove(tempFile.Name()) // Clean up temp file on error
		return "", err
	}

	if err := tempFile.Close(); err != nil {
		return "", fmt.Errorf("failed to close temporary file: %w", err)
	}

	return tempFile.Name(), nil
}

func appendFileContents(tempFile *os.File, filePath string) error {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".txt":
		return appendTextFileContents(tempFile, filePath)
	case ".pdf":
		return appendPdfFileContents(tempFile, filePath)
	default:
		return fmt.Errorf("unsupported file format: %s", ext)
	}
}

// appendTextFileContents appends the contents of a text file.
func appendTextFileContents(tempFile *os.File, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open text file %s: %w", filePath, err)
	}
	defer file.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		return fmt.Errorf("failed to copy text file contents: %w", err)
	}

	if _, err := tempFile.WriteString("\n"); err != nil {
		return fmt.Errorf("failed to write newline: %w", err)
	}

	return nil
}

func appendPdfFileContents(tempFile *os.File, filePath string) error {
	// Open the PDF file
	f, r, err := pdf.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open PDF file %s: %w", filePath, err)
	}
	defer f.Close()

	// Loop through all pages to extract text
	totalPages := r.NumPage()
	for i := 1; i <= totalPages; i++ {
		page := r.Page(i)
		if page.V.IsNull() {
			return fmt.Errorf("failed to read page %d in PDF %s", i, filePath)
		}

		// Extract text from the page
		text, err := page.GetPlainText(nil)
		if err != nil {
			return fmt.Errorf("failed to extract text from page %d in PDF %s: %w", i, filePath, err)
		}

		// Write the extracted text to the temporary file
		if _, err := tempFile.WriteString(text + "\n"); err != nil {
			return fmt.Errorf("failed to write PDF text: %w", err)
		}
	}

	return nil
}

// deleteTempFile deletes the specified temporary file.
func deleteTempFile() error {
	if err := os.Remove(tempFileLocation); err != nil {
		return fmt.Errorf("can not delete temporary file: %w", err)
	}
	return nil
}

func getTempFileLocation() string {
	return tempFileLocation
}

func setTempFileLocation(location string) {
	tempFileLocation = location
}

// NOTE: Uncomment the main function to run the file merger
// func main() {
// 	// Change "your_directory_path" to the directory you want to process
// 	directory := "your_directory_path"

// 	if err := mergeFilesToTemp(directory); err != nil {
// 		fmt.Printf("Error: %v\n", err)
// 	} else {
// 		fmt.Println("Temporary file processed and deleted successfully.")
// 	}
// }
