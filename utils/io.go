package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func WriteToPath(path, text string) (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user's home directory: %v", err)
	}

	finalPath := filepath.Join(homePath, ".leaf", path)
	dir := filepath.Dir(finalPath)

	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}

	err = ioutil.WriteFile(finalPath, []byte(text), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write file: %v", err)
	}

	return finalPath, nil
}
