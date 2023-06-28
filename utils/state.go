package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type State struct {
	Token        string
	Applications []Application
}

type Application struct {
	Name      string
	Connected bool
}

func getStatePath() string {
	homePath, _ := os.UserHomeDir()
	statePath := filepath.Join(homePath, ".leaf", "state.txt")

	// Create the directory if it doesn't exist
	if _, err := os.Stat(filepath.Dir(statePath)); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(statePath), 0755)
	}

	// Create the file if it doesn't exist
	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		file, _ := os.Create(statePath)
		file.Close()
	}

	return statePath
}

func GetState() (State, error) {
	tokenFilePath := getStatePath()
	content, err := ioutil.ReadFile(tokenFilePath)
	if err != nil {
		return State{}, err
	}

	if len(content) == 0 {
		// Return default state if content is empty
		return State{
			Token:        "",
			Applications: []Application{},
		}, nil
	}

	var state State
	err = json.Unmarshal(content, &state)
	if err != nil {
		return State{}, err
	}

	return state, nil
}

func SetState(state PartialState) error {
	tokenFilePath := getStatePath()
	currentState, err := GetState()
	if err != nil {
		return err
	}

	updatedState := State{
		Token:        currentState.Token,
		Applications: currentState.Applications,
	}

	if state.Token != "" {
		updatedState.Token = state.Token
	}

	if len(state.Applications) > 0 {
		updatedState.Applications = append(updatedState.Applications, state.Applications...)
	}

	stateBytes, err := json.Marshal(updatedState)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(tokenFilePath, stateBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func AddApplication(name string) error {
	_, err := GetState()
	if err != nil {
		return err
	}

	application := Application{
		Name:      name,
		Connected: false,
	}

	updatedState := PartialState{
		Applications: []Application{application},
	}

	err = SetState(updatedState)
	if err != nil {
		return err
	}

	return nil
}

type PartialState struct {
	Token        string
	Applications []Application
}
