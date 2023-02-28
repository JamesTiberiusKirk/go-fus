package fus

import (
	"fmt"
	"os"
	"path/filepath"
)

type FileHandlerType string

const (
	SingleFolder FileHandlerType = "single-file"
	MultiFolder  FileHandlerType = "multi-folder" // TODO: Not inplemented
)

// FileHandler file handler interface.
type FileHandler func(config viewEngineConfig, tplFile string) (content string, err error)

func getFileHandler(t FileHandlerType) FileHandler {
	switch t {
	case MultiFolder:
		// TODO: Need to actually implement this
		return func(config viewEngineConfig, tplFile string) (string, error) {
			// Get the absolute path of the root template
			path, err := filepath.Abs(config.Root + string(os.PathSeparator) + tplFile)
			if err != nil {
				return "", fmt.Errorf("ViewEngine handler:%s path:%v error: %w", MultiFolder, path, err)
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return "", fmt.Errorf("ViewEngine render read handler:%s name:%v, path:%v, error: %w",
					MultiFolder, tplFile, path, err)
			}
			return string(data), nil
		}
	case SingleFolder:
		return func(config viewEngineConfig, tplFile string) (string, error) {
			// Get the absolute path of the root template
			path, err := filepath.Abs(config.Root + string(os.PathSeparator) + tplFile)
			if err != nil {
				return "", fmt.Errorf("ViewEngine handler:%s path:%v error: %w", SingleFolder, path, err)
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return "", fmt.Errorf("ViewEngine render read handler:%s name:%v, path:%v, error: %w",
					SingleFolder, tplFile, path, err)
			}
			return string(data), nil
		}
	default:
		return getFileHandler(SingleFolder)
	}
}
