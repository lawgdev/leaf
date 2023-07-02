package commands

import (
	"os"
	"runtime"

	"github.com/urfave/cli/v2"
)

func Upstart(c *cli.Context) error {
	switch runtime.GOOS {
	case "darwin":
		return macOS(c)
	case "linux":
		return ubuntu(c)
	default:
		return cli.Exit("This command is only available on macOS and Ubuntu.", 1)
	}
}

func macOS(c *cli.Context) error {
	var launchdConfig = `
<?xml version="1.0" encoding="UTF-8"?>
	<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
	<plist version="1.0">
	  <dict>
	    <key>Label</key>
	    <string>dev.lawg.leaf</string>
	    <key>ProgramArguments</key>
	    <array>
	      <string>/usr/local/bin/leaf</string>
	      <string>listen</string>
	    </array>
	    <key>RunAtLoad</key>
	    <true/>
	  </dict>
	</plist>
`

	f, err := os.Create("/Library/LaunchDaemons/dev.lawg.leaf.plist")
	if err != nil {
		return cli.Exit("Failed to create launchd config", 1)
	}

	defer f.Close()

	_, err = f.WriteString(launchdConfig)
	if err != nil {
		return cli.Exit("Failed to write launchd config", 1)
	}

	return nil
}

func ubuntu(c *cli.Context) error {
	if _, err := os.Stat("/lib/systemd"); os.IsNotExist(err) {
		return cli.Exit("This command is only available on Ubuntu 6.10 and higher.", 1)
	}

	var systemCtlConfig = `
[Unit]
Description=Leaf
After=network.target

[Service]
Type=simple
Restart=always
RestartSec=1
ExecStart=/usr/bin/leaf listen
RemainAfterExit=true

[Install]
WantedBy=multi-user.target
`

	f, err := os.Create("/lib/systemd/system/leaf.service")
	if err != nil {
		return cli.Exit("Failed to create upstart config", 1)
	}

	defer f.Close()

	_, err = f.WriteString(systemCtlConfig)
	if err != nil {
		return cli.Exit("Failed to write systemctl config", 1)
	}
	return nil
}
