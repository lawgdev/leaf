package main

import (
	"log"
	"os"

	"leaf/commands"

	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "login",
				Aliases: []string{"init"},
				Usage:   "Login to lawg",
				Action:  commands.Login,
			},
			{
				Name:    "connect",
				Aliases: []string{"c", "con"},
				Usage:   "Connect your apps to an application feed on lawg",
				Action:  commands.Connect,
			},
			{
				Name:    "listen",
				Aliases: []string{"l", "li"},
				Usage:   "Listen to your application feeds on lawg",
				Action:  commands.Listen,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
