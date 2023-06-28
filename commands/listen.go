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
		return fmt.Errorf("failed to get user's home directory: %v", err)
	}

	var cmd = exec.Command("vector", "-c", homePath+"/.leaf/configs/*.toml")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	go printOutput(stdout)
	go printOutput(stderr)

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
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
