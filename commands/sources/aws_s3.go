package sources

import (
	"fmt"
	"leaf/utils"
	"os/exec"

	"github.com/AlecAivazis/survey/v2"
)

func AwsS3(feed utils.Feed, project utils.Project) error {
	// ask for queue URL for sqs (its required)
	var queueUrl = ""
	survey.AskOne(&survey.Input{
		Message: "SQS Queue URL:",
	}, &queueUrl, survey.WithValidator(survey.Required))

	var region = ""
	survey.AskOne(&survey.Input{
		Message: "Region:",
		Suggest: func(toComplete string) []string {
			return []string{"us-east-1", "us-east-2", "us-west-1", "us-west-2", "ap-east-1", "ap-south-1", "ap-northeast-2", "ap-southeast-1", "ap-southeast-2", "ap-northeast-1", "ca-central-1", "cn-north-1", "cn-northwest-1", "eu-central-1", "eu-west-1", "eu-west-2", "eu-west-3", "eu-north-1", "me-south-1", "sa-east-1"}
		},
	}, &region, survey.WithValidator(survey.Required))

	var accessKeyId = ""
	survey.AskOne(&survey.Input{
		Message: "Access Key ID:",
	}, &accessKeyId, survey.WithValidator(survey.Required))

	var secretAccessKey = ""
	survey.AskOne(&survey.Input{
		Message: "Secret Access Key:",
	}, &secretAccessKey, survey.WithValidator(survey.Required))

	cmd := exec.Command("vector", "generate", "aws_s3")
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

[sources.source0.auth]
access_key_id = "%s"
secret_access_key = "%s"
region = "%s"

%s`, parsedVectorConfig, accessKeyId, secretAccessKey, region, utils.DefaultConfig(project.Namespace, feed.Name, state.Token))

	finalPath, err := utils.WriteToPath(fmt.Sprintf("configs/%s-aws-s3.toml", feed.Name), config)
	if err != nil {
		return utils.ParsedError(err, "Failed to write config", true)
	}

	fmt.Println("Config generated and saved to", finalPath)
	return nil
}
