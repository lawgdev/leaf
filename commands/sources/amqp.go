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

	state, err := utils.GetState()
	if err != nil {
		return utils.ParsedError(err, "Failed to get state", true)
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

	config := fmt.Sprintf(`
%s

%s`, parsedVectorConfig, utils.DefaultConfig(project.Namespace, feed.Name, state.Token))

	finalPath, err := utils.WriteToPath(fmt.Sprintf("configs/%s-amqp%s.toml", feed.Name, queueName), config)
	if err != nil {
		return utils.ParsedError(err, "Failed to write config", true)
	}

	fmt.Println("Config generated and saved to", finalPath)
	return nil
}
