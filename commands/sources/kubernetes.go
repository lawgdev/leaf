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

	state, err := utils.GetState()
	if err != nil {
		return utils.ParsedError(err, "Failed to get state", true)
	}

	var parsedVectorConfig = utils.ParseVectorConfig(vectorGeneratedConfig)

	config := fmt.Sprintf(`
%s

%s`, parsedVectorConfig, utils.DefaultConfig(project.Namespace, feed.Name, state.Token))

	finalPath, err := utils.WriteToPath(fmt.Sprintf("configs/%s-kubernetes.toml", feed.Name), config)
	if err != nil {
		return utils.ParsedError(err, "Failed to write config", true)
	}

	fmt.Println("Config generated and saved to", finalPath)
	return nil
}
