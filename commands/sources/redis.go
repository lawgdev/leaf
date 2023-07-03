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

	state, err := utils.GetState()
	if err != nil {
		return utils.ParsedError(err, "Failed to get state", true)
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

	config := fmt.Sprintf(`
%s

%s`, parsedVectorConfig, utils.DefaultConfig(project.Namespace, feed.Name, state.Token))

	finalPath, err := utils.WriteToPath(fmt.Sprintf("configs/%s-redis.toml", feed.Name), config)
	if err != nil {
		return utils.ParsedError(err, "Failed to write config", true)
	}

	fmt.Println("Config generated and saved to", finalPath)
	return nil
}
