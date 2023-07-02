package commands

import (
	"fmt"
	"leaf/utils"

	"github.com/AlecAivazis/survey/v2"
	"github.com/urfave/cli/v2"
)

func Login(ctx *cli.Context) error {
	token := ""
	survey.AskOne(&survey.Input{Message: "Enter your lawg.dev token (https://app.lawg.dev/user/settings):"}, &token, survey.WithValidator(survey.Required))

	spinner := utils.Spinner.AddSpinner("Authenticating with lawg")
	utils.Spinner.Start()
	me, err := utils.GetMe(token)

	if err != nil {
		if err.Error() == "You must be authenticated to access this route" {
			spinner.Error()
			spinner.UpdateMessage("Invalid token, please use a valid API or session token.")
			utils.Spinner.Stop()

			return cli.Exit("", 1)
		}

		spinner.Error()
		spinner.UpdateMessage(err.Error())
		utils.Spinner.Stop()

		return cli.Exit("", 1)
	}

	err = utils.SetState(utils.PartialState{
		Token: token,
	})

	if err != nil {
		return cli.Exit("Failed to save state", 1)
	}

	var user = me.Data.User

	spinner.Complete()
	spinner.UpdateMessage(fmt.Sprintf("Logged in as %s (%s)", user.Username, user.Email))
	utils.Spinner.Stop()

	return nil
}
