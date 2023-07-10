package sources

import (
	"fmt"
	"leaf/utils"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

func AMQP(feed utils.Feed, project utils.Project) error {
	cmd := exec.Command("vector", "generate", "amqp")
	vectorGeneratedConfig, err := cmd.Output()
	if err != nil {
		return utils.ParsedError(err, "Failed to generate config", true)
	}

	// ask for amqp URL
	var amqpUrl = ""
	survey.AskOne(&survey.Input{
		Message: "AMQP URL:",
		Suggest: func(toComplete string) []string {
			return []string{"amqp://default:pass@127.0.0.1:5672/%2f", "amqp://default:pass@0.0.0.0:5672/%2f"}
		},
	}, &amqpUrl, survey.WithValidator(survey.Required))

	var queueName = ""
	survey.AskOne(&survey.Input{
		Message: "Queue name (or leave blank for all):",
	}, &queueName)

	var parsedVectorConfig = utils.ParseVectorConfig(vectorGeneratedConfig)
	parsedVectorConfig = strings.Replace(parsedVectorConfig, "amqp://user:password@127.0.0.1:5672/%2f?timeout=10", amqpUrl, 1)
	parsedVectorConfig = strings.ReplaceAll(parsedVectorConfig, "queue = \"\"", fmt.Sprintf("queue = \"%s\"", queueName))

	if err := utils.GenerateConfig(parsedVectorConfig, fmt.Sprintf("%s-amqp%s", feed.Name, queueName), project.Namespace, feed.Name); err != nil {
		return utils.ParsedError(err, "Failed to generate config", true)
	}

	return nil
}
