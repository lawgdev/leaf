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

	// ask for syslog address
	var addressUrl = ""
	if err := survey.AskOne(&survey.Input{
		Message: "Address URL:",
		Suggest: func(toComplete string) []string {
			return []string{"127.0.0.1:514", "0.0.0.0:514"}
		},
	}, &addressUrl, survey.WithValidator(survey.Required)); err != nil {
		return utils.ParsedError(err, "Failed to get address URL", true)
	}

	var mode = ""
	if err := survey.AskOne(&survey.Select{
		Message: "Mode:",
		Options: []string{"tcp", "udp"},
	}, &mode, survey.WithValidator(survey.Required)); err != nil {
		return utils.ParsedError(err, "Failed to get mode", true)
	}

	var parsedVectorConfig = utils.ParseVectorConfig(vectorGeneratedConfig)
	parsedVectorConfig = strings.Replace(parsedVectorConfig, "0.0.0.0:514", addressUrl, 1)
	parsedVectorConfig = strings.Replace(parsedVectorConfig, "tcp", mode, 1)

	if err := utils.GenerateConfig(parsedVectorConfig, fmt.Sprintf("%s-syslog", feed.Name), project.Namespace, feed.Name); err != nil {
		return utils.ParsedError(err, "Failed to generate config", true)
	}

	return nil
}
