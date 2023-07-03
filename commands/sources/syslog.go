package sources

import (
	"fmt"
	"leaf/utils"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

func Syslog(feed utils.Feed, project utils.Project) error {
	cmd := exec.Command("vector", "generate", "syslog")
	vectorGeneratedConfig, err := cmd.Output()
	if err != nil {
		return utils.ParsedError(err, "Failed to generate config", true)
	}

	state, err := utils.GetState()
	if err != nil {
		return utils.ParsedError(err, "Failed to get state", true)
	}

	// ask for syslog address
	var addressUrl = ""
	survey.AskOne(&survey.Input{
		Message: "Address URL:",
		Suggest: func(toComplete string) []string {
			return []string{"127.0.0.1:514", "0.0.0.0:514"}
		},
	}, &addressUrl, survey.WithValidator(survey.Required))

	var mode = ""
	survey.AskOne(&survey.Select{
		Message: "Mode:",
		Options: []string{"tcp", "udp"},
	}, &mode, survey.WithValidator(survey.Required))

	var parsedVectorConfig = utils.ParseVectorConfig(vectorGeneratedConfig)
	parsedVectorConfig = strings.Replace(parsedVectorConfig, "0.0.0.0:514", addressUrl, 1)
	parsedVectorConfig = strings.Replace(parsedVectorConfig, "tcp", mode, 1)

	config := fmt.Sprintf(`
%s

%s`, parsedVectorConfig, utils.DefaultConfig(project.Namespace, feed.Name, state.Token))

	finalPath, err := utils.WriteToPath(fmt.Sprintf("configs/%s-syslog.toml", feed.Name), config)
	if err != nil {
		return utils.ParsedError(err, "Failed to write config", true)
	}

	fmt.Println("Config generated and saved to", finalPath)
	return nil
}
