package commands

import (
	"encoding/json"
	"io"
	"leaf/twig"
	"leaf/utils"
	"log"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
)

func Listen(ctx *cli.Context) error {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return utils.ParsedError(err, "Failed to get home directory", true)
	}

	_, err = utils.GetOriginsFromConfigs()
	if err != nil {
		return utils.ParsedError(err, "Failed to get listenable configs", true)
	}

	state, err := utils.GetState()
	if err != nil {
		return utils.ParsedError(err, "Failed to get state", true)
	}

	// Connect to websocket
	var twigClient = twig.Connect(state.Token)

	var cmd = exec.Command("vector", "-c", homePath+"/.leaf/configs/*.toml")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return utils.ParsedError(err, "Failed to get stdout pipe", true)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return utils.ParsedError(err, "Failed to get stderr pipe", true)
	}

	if err := cmd.Start(); err != nil {
		return utils.ParsedError(err, "Failed to start command", true)
	}

	go printOutput(stdout, twigClient)
	go printOutput(stderr, twigClient)

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		return utils.ParsedError(err, "Failed to wait for command", true)
	}

	return nil
}

func printOutput(pipe io.Reader, twigClient *twig.Twig) {
	buf := make([]byte, 1024)
	for {
		n, err := pipe.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}

			return
		}

		var message = string(buf[:n])

		// check if message is parseable into a json object
		var obj map[string]interface{}
		if err := json.Unmarshal([]byte(message), &obj); err != nil {
			println("not parsable")
			continue
		}

		twigClient.WS.WriteJSON(twig.CreateLogMessage{
			Op: 5,
			Data: twig.CreateLogData{
				Message:          message,
				Level:            "info",
				ProjectNamespace: "test",
				FeedName:         "test",
			},
		})
	}
}
