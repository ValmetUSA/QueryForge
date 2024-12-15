package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// mergeFilesToTemp reads all files from the selected directory and outputs their
// combined content into a single temporary file, prints it, and deletes it.
func mergeFilesToTemp(dir string) error {
	tempDir := os.TempDir()
	tempFile, err := os.CreateTemp(tempDir, "merged_*.txt")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer tempFile.Close()

	// Traverse the directory and process each file
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if err := appendFileContents(tempFile, path); err != nil {
				return fmt.Errorf("failed to process file %s: %w", path, err)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Print the content of the temporary file to the terminal
	if err := printFileContents(tempFile.Name()); err != nil {
		return fmt.Errorf("failed to print temporary file contents: %w", err)
	}

	// Delete the temporary file
	if err := os.Remove(tempFile.Name()); err != nil {
		return fmt.Errorf("failed to delete temporary file: %w", err)
	}

	return nil
}

// appendFileContents appends the contents of a given file to the temp file.
func appendFileContents(tempFile *os.File, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	// Add a newline separator between files (optional)
	if _, err := tempFile.WriteString("\n"); err != nil {
		return fmt.Errorf("failed to write newline: %w", err)
	}

	return nil
}

// printFileContents reads and prints the contents of a file to the terminal.
func printFileContents(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file for printing: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	return nil
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
