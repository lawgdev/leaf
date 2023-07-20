package main

import (
	"leaf/commands"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Name:  "lawg",
		Usage: "your logs, one place.",
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
			{
				Name:   "upstart",
				Usage:  "Configure lawg to start on systemctl",
				Action: commands.Upstart,
			},
			{
				Name:   "whoami",
				Usage:  "Get your user information",
				Action: commands.WhoAmI,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
