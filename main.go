package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func retrieveFile(fullPath string) string {
	return filepath.Base(fullPath)
}

func scanAndClassify(root string, extensionMap map[string]string) (map[string]string, error) {
	fileMap := make(map[string]string)

	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !entry.IsDir() {
			ext := strings.ToLower(filepath.Ext(entry.Name()))

			if category, exists := extensionMap[ext]; exists {
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

		srcFile, err := os.Open(fullPath)
		if err != nil {
			return fmt.Errorf("failed to open source file %s: %w", fullPath, err)
		}
		defer srcFile.Close()

		destFile, err := os.Create(destPath)
		if err != nil {
			return fmt.Errorf("failed to create destination file %s: %w", destPath, err)
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			return fmt.Errorf("failed to copy from %s to %s: %w", fullPath, destPath, err)
		}

		fmt.Printf("Copied %s to %s\n", fullPath, destPath)
	}
	return nil
}

func printSummary(fileMap map[string]string) {
	fmt.Println("\n📦 Final file map:")
	for fullPath, category := range fileMap {
		fmt.Printf("%s (%s) => %s\n", fullPath, retrieveFile(fullPath), category)
	}
}
func loadExtensionMap() (map[string]string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	// Get the folder where main.exe is located
	exeDir := filepath.Dir(exePath)

	// Build full path to config.json next to the exe
	path := filepath.Join(exeDir, "config.json")

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var extMap map[string]string
	err = json.Unmarshal(data, &extMap)
	if err != nil {
		return nil, err
	}

	return extMap, nil
}

func main() {

	configmap, err := loadExtensionMap()
	if err != nil {
		fmt.Println("Error loading extension map:", err)
		return
	}
	fmt.Println("Loaded extension map:", configmap)

	start := time.Now()

	root := "./"
	fileMap, err := scanAndClassify(root, configmap)
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
