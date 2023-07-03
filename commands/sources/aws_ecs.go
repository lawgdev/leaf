package sources

import (
	"fmt"
	"leaf/utils"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

func AwsECSMetrics(feed utils.Feed, project utils.Project) error {
	cmd := exec.Command("vector", "generate", "aws_ecs_metrics")
	vectorGeneratedConfig, err := cmd.Output()
	if err != nil {
		return utils.ParsedError(err, "Failed to generate config", true)
	}

	state, err := utils.GetState()
	if err != nil {
		return utils.ParsedError(err, "Failed to get state", true)
	}

	// ask for endpoint URL
	var endpoint = ""
	survey.AskOne(&survey.Input{
		Message: "Endpoint URL:",
		Suggest: func(toComplete string) []string {
			return []string{"http://169.254.170.2/v2"}
		},
	}, &endpoint, survey.WithValidator(survey.Required))

	var namespace = ""
	survey.AskOne(&survey.Input{
		Message: "Namespace (or empty for disabled):",
	}, &namespace)

	var parsedVectorConfig = utils.ParseVectorConfig(vectorGeneratedConfig)
	parsedVectorConfig = strings.Replace(parsedVectorConfig, "http://169.254.170.2/v2", endpoint, 1)
	parsedVectorConfig = strings.Replace(parsedVectorConfig, "namespace = \"awsecs\"", fmt.Sprintf("namespace = \"%s\"", namespace), 1)

	config := fmt.Sprintf(`
%s

%s`, parsedVectorConfig, utils.DefaultConfig(project.Namespace, feed.Name, state.Token))

	finalPath, err := utils.WriteToPath(fmt.Sprintf("configs/%s-aws-ecs%s.toml", feed.Name, namespace), config)
	if err != nil {
		return utils.ParsedError(err, "Failed to write config", true)
	}

	fmt.Println("Config generated and saved to", finalPath)
	return nil
}
