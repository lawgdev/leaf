package sources

import (
	"fmt"
	"leaf/utils"
	"os/exec"
)

func KubernetesLogs(feed utils.Feed, project utils.Project) error {
	cmd := exec.Command("vector", "generate", "kubernetes_logs")
	vectorGeneratedConfig, err := cmd.Output()
	if err != nil {
		return utils.ParsedError(err, "Failed to generate config", true)
	}

	var parsedVectorConfig = utils.ParseVectorConfig(vectorGeneratedConfig)

	if err := utils.GenerateConfig(parsedVectorConfig, fmt.Sprintf("%s-kubernetes", feed.Name), project.Namespace, feed.Name); err != nil {
		return utils.ParsedError(err, "Failed to generate config", true)
	}

	return nil
}
