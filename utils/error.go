package utils

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

func ParsedError(err error, genericError string, printErr ...bool) error {
	if strings.Contains(err.Error(), "permission denied") {
		return cli.Exit("This command requires root privileges.", 1)
	}

	if strings.Contains(err.Error(), "You must be authenticated to access this route") {
		return cli.Exit("Please use `leaf init` to login to your lawg account.", 1)
	}

	if len(printErr) > 0 && printErr[0] {
		return cli.Exit(fmt.Sprintf("%s: %s", genericError, err.Error()), 1)
	}

	return cli.Exit(fmt.Sprintf("%s", genericError), 1)
}
