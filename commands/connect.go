package commands

import (
	"leaf/commands/sources"
	"leaf/utils"
	"os/exec"

	"github.com/AlecAivazis/survey/v2"
	"github.com/urfave/cli/v2"
)

type Source struct {
	Name    string
	Handler func(feed utils.Feed, project utils.Project) error
}

var supportedSources = map[string]Source{
	"docker_logs": {
		Name:    "Docker Logs",
		Handler: sources.DockerLogs,
	},
	"kubernetes_logs": {
		Name:    "Kubernetes Logs",
		Handler: sources.KubernetesLogs,
	},
	"amqp": {
		Name:    "AMQP",
		Handler: sources.AMQP,
	},
	"redis": {
		Name:    "Redis",
		Handler: sources.Redis,
	},
	"aws_ecs_metrics": {
		Name:    "AWS ECS Metrics",
		Handler: sources.AwsECSMetrics,
	},
	"aws_s3": {
		Name:    "AWS S3",
		Handler: sources.AwsS3,
	},
	"file": {
		Name:    "File",
		Handler: sources.File,
	},
	"syslog": {
		Name:    "Syslog",
		Handler: sources.Syslog,
	},
}

var supportSourcesOrder = []string{
	"docker_logs",
	"kubernetes_logs",
	"amqp",
	"redis",
	"aws_ecs_metrics",
	"aws_s3",
	"file",
	"syslog",
}

func Connect(ctx *cli.Context) error {
	var cmd = exec.Command("vector", "--help")

	_, err := cmd.Output()
	if err != nil {
		return utils.ParsedError(err, "Unable to find Vector! Please install it at: https://vector.dev/docs/setup/installation")
	}

	state, err := utils.GetState()
	if err != nil {
		return utils.ParsedError(err, "Unable to get state", true)
	}

	projectSpinner := utils.Spinner.AddSpinner("Fetching Projects")
	utils.Spinner.Start()

	me, err := utils.GetMe(state.Token)

	if err != nil {
		projectSpinner.Error()
		projectSpinner.UpdateMessage(err.Error())
		utils.Spinner.Stop()

		return cli.Exit("", 1)
	}

	projectSpinner.UpdateMessage("Fetched Projects")
	projectSpinner.Complete()
	utils.Spinner.Stop()

	var projects = []string{}
	for _, project := range me.Data.Projects {
		projects = append(projects, project.Namespace)
	}

	var selectedProjectNamespace = ""
	survey.AskOne(&survey.Select{
		Message: "Select Project",
		Options: projects,
	}, &selectedProjectNamespace, survey.WithValidator(survey.Required))

	// get the selected project by namespace
	var selectedProject utils.Project
	for _, project := range me.Data.Projects {
		if project.Namespace == selectedProjectNamespace {
			selectedProject = project
		}
	}

	var feeds = []string{}
	for _, feed := range selectedProject.Feeds {
		// if feed is an event feed and not an application feed
		if feed.Type == 0 {
			continue
		}
		feeds = append(feeds, feed.Name)
	}

	if len(feeds) == 0 {
		println("No application feeds found")
		return nil
	}

	survey.AskOne(&survey.Select{
		Message: "Select An Application Feed",
		Options: feeds,
	}, &selectedProjectNamespace, survey.WithValidator(survey.Required))

	var selectedFeed utils.Feed
	for _, feed := range selectedProject.Feeds {
		if feed.Name == selectedProjectNamespace {
			selectedFeed = feed
		}
	}

	var sourceNames = []string{}
	for _, k := range supportSourcesOrder {
		sourceNames = append(sourceNames, supportedSources[k].Name)
	}

	var selectedSourceName = ""
	survey.AskOne(&survey.Select{
		Message: "Select A Source",
		Options: sourceNames,
	}, &selectedSourceName, survey.WithValidator(survey.Required))

	var selectedSource Source
	for _, source := range supportedSources {
		if source.Name == selectedSourceName {
			selectedSource = source
		}
	}

	return selectedSource.Handler(selectedFeed, selectedProject)
}
