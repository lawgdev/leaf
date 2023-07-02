package commands

import (
	"fmt"
	"io"
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

	go printOutput(stdout)
	go printOutput(stderr)

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		return utils.ParsedError(err, "Failed to wait for command", true)
	}

	return nil
}

func printOutput(pipe io.Reader) {
	buf := make([]byte, 1024)
	for {
		n, err := pipe.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}

			return
		}

		fmt.Print(string(buf[:n]))
	}
}
