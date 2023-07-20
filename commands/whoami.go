package commands

import (
	"fmt"
	"leaf/utils"

	"github.com/urfave/cli/v2"
)

func WhoAmI(c *cli.Context) error {
	state, err := utils.GetState()

	if err != nil {
		return utils.ParsedError(err, "Failed to get state", true)
	}

	spinner := utils.Spinner.AddSpinner("Authenticating with lawg")
	utils.Spinner.Start()
	me, err := utils.GetMe(state.Token)

	if err != nil {
		if err.Error() == "You must be authenticated to access this route" {
			spinner.Error()
			spinner.UpdateMessage("Invalid token or not logged in, please use a valid API or session token and login with 'leaf login'")
			utils.Spinner.Stop()

			return cli.Exit("", 1)
		}

		spinner.Error()
		spinner.UpdateMessage(err.Error())
		utils.Spinner.Stop()

		return cli.Exit("", 1)
	}

	return cli.Exit(fmt.Sprintf("Logged in as %s (%s)", me.Data.User.Username, me.Data.User.Email), 0)
}
