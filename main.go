package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func retrieveFile(fullPath string) string {
	return filepath.Base(fullPath)
}

func scanAndClassify(root string) (map[string]string, error) {
	fileMap := make(map[string]string)

	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !entry.IsDir() {
			ext := strings.ToLower(filepath.Ext(entry.Name()))
			category := ""

			switch ext {
			case ".txt", ".pdf", ".docx":
				category = "document"
			case ".jpg", ".jpeg", ".png", ".gif":
				category = "image"
			}

			if category != "" {
				fileMap[path] = category
			}
		}
		return nil
	})

	return fileMap, err
}

func copyFiles(fileMap map[string]string) error {
	for fullPath, category := range fileMap {
		destDir := filepath.Join("organized", category)
		if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", destDir, err)
		}

		destPath := filepath.Join(destDir, retrieveFile(fullPath))
		if err := os.Rename(fullPath, destPath); err != nil {
			return fmt.Errorf("failed to move file %s to %s: %w", fullPath, destPath, err)
		}
		fmt.Printf("Moved %s to %s\n", fullPath, destPath)
	}
	return nil
}

func printSummary(fileMap map[string]string) {
	fmt.Println("\n📦 Final file map:")
	for fullPath, category := range fileMap {
		fmt.Printf("%s (%s) => %s\n", fullPath, retrieveFile(fullPath), category)
	}
}

func main() {
	start := time.Now()

	root := "./"
	fileMap, err := scanAndClassify(root)
	if err != nil {
		fmt.Println("Error scanning files:", err)
		return
	}

	copyFiles(fileMap)

	printSummary(fileMap)

	duration := time.Since(start)
	fmt.Printf("\n⏱ Time taken: %s\n", duration)
	fmt.Println("\nPress Enter to exit...")
	fmt.Scanln()
}
