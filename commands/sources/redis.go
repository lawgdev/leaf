package sources

import (
	"fmt"
	"leaf/utils"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

func Redis(feed utils.Feed, project utils.Project) error {
	cmd := exec.Command("vector", "generate", "redis")
	vectorGeneratedConfig, err := cmd.Output()
	if err != nil {
		return utils.ParsedError(err, "Failed to generate config", true)
	}

	// ask for redis URL
	var redisUrl = ""
	survey.AskOne(&survey.Input{
		Message: "Redis URL:",
		Suggest: func(toComplete string) []string {
			return []string{"redis://localhost:6379/0", "redis://redis:6379/0", "redis://0.0.0.0:6379/0"}
		},
	}, &redisUrl, survey.WithValidator(survey.Required))

	var parsedVectorConfig = utils.ParseVectorConfig(vectorGeneratedConfig)
	parsedVectorConfig = strings.Replace(parsedVectorConfig, "redis://127.0.0.1:6379/0", redisUrl, 1)

	if err := utils.GenerateConfig(parsedVectorConfig, fmt.Sprintf("%s-redis", feed.Name), project.Namespace, feed.Name); err != nil {
		return utils.ParsedError(err, "Failed to generate config", true)
	}

	return nil
}
