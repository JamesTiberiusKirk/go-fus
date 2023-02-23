package renderer

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// defaultFileHandler new default file handler
// This is just supposed to read the file and return a string as content.
func defaultFileHandler() FileHandler {
	return func(config Config, tplFile string) (string, error) {
		// Get the absolute path of the root template
		path, err := filepath.Abs(config.Root + string(os.PathSeparator) + tplFile)
		if err != nil {
			return "", fmt.Errorf("ViewEngine path:%v error: %w", path, err)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("ViewEngine render read name:%v, path:%v, error: %w", tplFile, path, err)
		}
		return string(data), nil
	}
}

// InFolderFileHandler new default file handler
// This is just supposed to read the file and return a string as content.
func InFolderFileHandler() FileHandler {
	return func(config Config, tplFile string) (string, error) {
		// Get the absolute path of the root template

		log.Println("tplFile:", tplFile)

		path, err := filepath.Abs(config.Root + string(os.PathSeparator) + tplFile)
		if err != nil {
			return "", fmt.Errorf("ViewEngine path:%v error: %w", path, err)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("ViewEngine render read name:%v, path:%v, error: %w", tplFile, path, err)
		}
		return string(data), nil
	}
}
