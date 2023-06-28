package sources

import (
	"context"
	"fmt"
	"leaf/utils"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func generateConfig(source, includeContainer, project, feed, token string) string {
	return fmt.Sprintf(`
%s
include_containers = [ "%s" ]

[sinks.lawg_sink]
type = "http"
encoding.codec = "json"
inputs = ["source0"]
uri = "http://100.127.114.55:8080/v1/projects/%s/feeds/%s/logs"
auth.strategy = "bearer"
auth.token = "%s"
`, strings.TrimSpace(source), includeContainer, project, feed, token)
}

func DockerLogs(feed utils.Feed, project utils.Project) {
	client, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := client.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	var containerNames = []string{}
	for _, container := range containers {
		containerNames = append(containerNames, fmt.Sprintf("%s (%s)", container.Names[0][1:], container.ID[:12]))
	}

	var selectedContainer string
	survey.AskOne(&survey.Select{
		Message: "Select a container:",
		Options: containerNames,
	}, &selectedContainer, survey.WithValidator(survey.Required))

	selectedContainer = strings.Split(selectedContainer, " ")[0]

	cmd := exec.Command("vector", "generate", "docker_logs")
	vectorGeneratedConfig, err := cmd.Output()
	stringedVectorGeneratedConfig := string(vectorGeneratedConfig[:])

	if err != nil {
		println("Vector does not exist, please install at https://vector.dev/docs/setup/installation/")
		return
	}

	lines := strings.Split(stringedVectorGeneratedConfig, "\n")
	lines = lines[2:] // remove first 2 lines of generatedConfig (irrelevant to config)

	stringedVectorGeneratedConfig = strings.Join(lines, "\n")

	state, err := utils.GetState()
	if err != nil {
		println("Error: " + err.Error())
		return
	}

	config := generateConfig(stringedVectorGeneratedConfig, selectedContainer, project.Namespace, feed.Name, state.Token)

	finalPath, err := utils.WriteToPath(fmt.Sprintf("configs/%s-%s.toml", feed.Name, selectedContainer), config)
	if err != nil {
		fmt.Println("Failed to write config file:", err)
		return
	}

	fmt.Println("Config generated and saved to", finalPath)
}
