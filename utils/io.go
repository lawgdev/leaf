package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func WriteIfNotExists(path, text string) (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user's home directory: %v", err)
	}

	finalPath := filepath.Join(homePath, ".leaf", path)
	if _, err := os.Stat(finalPath); os.IsNotExist(err) {
		return WriteToPath(path, text)
	}

	return path, nil
}

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

	if err := ioutil.WriteFile(finalPath, []byte(text), 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %v", err)
	}

	return finalPath, nil
}

func GetFileContents(filename string) (string, error) {
	var path string

	if !strings.HasPrefix(filename, "/") {
		homePath, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user's home directory: %v", err)
		}

		path = filepath.Join(homePath, ".leaf", filename)
	} else {
		path = filename
	}

	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	return string(fileContents), nil
}

func WalkMatch(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func GetOriginsFromConfigs() ([]string, error) {
	var configs []string = []string{}
	entries, err := WalkMatch(os.Getenv("HOME")+"/.leaf/configs/", "*.toml")

	for _, entry := range entries {
		if strings.Contains(entry, "sinks.toml") {
			continue
		}

		contents, err := GetFileContents(entry)
		if err != nil {
			return nil, err
		}

		firstLine := strings.Split(contents, "\n")[0]
		configs = append(configs, firstLine[1:])
	}

	if err != nil {
		return nil, err
	}

	return configs, nil
}
