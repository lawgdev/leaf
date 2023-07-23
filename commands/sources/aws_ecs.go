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

	// ask for endpoint URL
	var endpoint = ""
	if err := survey.AskOne(&survey.Input{
		Message: "Endpoint URL:",
		Suggest: func(toComplete string) []string {
			return []string{"http://169.254.170.2/v2"}
		},
	}, &endpoint, survey.WithValidator(survey.Required)); err != nil {
		return utils.ParsedError(err, "Failed to get endpoint URL", true)
	}

	var namespace = ""
	if err := survey.AskOne(&survey.Input{
		Message: "Namespace (or empty for disabled):",
	}, &namespace); err != nil {
		return utils.ParsedError(err, "Failed to get namespace", true)
	}

	var parsedVectorConfig = utils.ParseVectorConfig(vectorGeneratedConfig)
	parsedVectorConfig = strings.Replace(parsedVectorConfig, "http://169.254.170.2/v2", endpoint, 1)
	parsedVectorConfig = strings.Replace(parsedVectorConfig, "namespace = \"awsecs\"", fmt.Sprintf("namespace = \"%s\"", namespace), 1)

	if err := utils.GenerateConfig(parsedVectorConfig, fmt.Sprintf("%s-aws-ecs%s", feed.Name, namespace), project.Namespace, feed.Name); err != nil {
		return utils.ParsedError(err, "Failed to generate config", true)
	}

	return nil
}
