package commands

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
)

func Listen(ctx *cli.Context) error {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return cli.Exit("Failed to get home directory", 1)
	}

	var cmd = exec.Command("vector", "-c", homePath+"/.leaf/configs/*.toml")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return cli.Exit("Failed to get stdout pipe", 1)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return cli.Exit("Failed to get stderr pipe", 1)
	}

	if err := cmd.Start(); err != nil {
		return cli.Exit("Failed to start command", 1)
	}

	go printOutput(stdout)
	go printOutput(stderr)

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		return cli.Exit("Failed to wait for command", 1)
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
