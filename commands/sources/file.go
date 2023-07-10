package sources

import (
	"fmt"
	"leaf/utils"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

func File(feed utils.Feed, project utils.Project) error {
	cmd := exec.Command("vector", "generate", "file")
	vectorGeneratedConfig, err := cmd.Output()
	if err != nil {
		return utils.ParsedError(err, "Failed to generate config", true)
	}

	var includePaths = []string{}
	var askForPaths = true

	for askForPaths {
		var includePath = ""
		survey.AskOne(&survey.Input{
			Message: "Include Path:",
			Suggest: func(toComplete string) []string {
				return []string{"/var/log/**/*.log"}
			},
		}, &includePath)

		if includePath == "" && len(includePaths) > 0 {
			askForPaths = false
			break
		}

		if len(includePaths) == 0 {
			println("You can add more or leave blank to continue.")
		}

		includePaths = append(includePaths, includePath)
	}

	var excludePaths = []string{}
	var askForExcludePaths = true

	for askForExcludePaths {
		var excludePath = ""
		survey.AskOne(&survey.Input{
			Message: "Exclude Path (optional):",
			Suggest: func(toComplete string) []string {
				return []string{"/var/log/binary-file.log"}
			},
		}, &excludePath)

		if excludePath == "" {
			askForPaths = false
			break
		}

		if len(excludePath) == 0 {
			println("You can add more or leave blank to continue.")
		}

		excludePaths = append(excludePaths, excludePath)
	}

	var parsedVectorConfig = utils.ParseVectorConfig(vectorGeneratedConfig)
	parsedVectorConfig = strings.Replace(parsedVectorConfig, "include = [\"/var/log/**/*.log\"]", "include = [\""+strings.Join(includePaths, "\", \"")+"\"]", 1)
	parsedVectorConfig = strings.Replace(parsedVectorConfig, "exclude = []", "exclude = [\""+strings.Join(excludePaths, "\", \"")+"\"]", 1)

	if err := utils.GenerateConfig(parsedVectorConfig, fmt.Sprintf("%s-file", feed.Name), project.Namespace, feed.Name); err != nil {
		return utils.ParsedError(err, "Failed to generate config", true)
	}

	return nil
}
