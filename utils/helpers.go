package utils

import (
	"fmt"
	"strings"
)

func ParseVectorConfig(config []byte) string {
	stringedVectorGeneratedConfig := string(config[:])

	lines := strings.Split(stringedVectorGeneratedConfig, "\n")
	lines = lines[2:] // remove first 2 lines of generatedConfig (irrelevant to config)

	return strings.Trim(strings.Join(lines, "\n"), "\n")
}

func GenerateConfig(config, sourceName, projectNamespace, feedName string) error {
	_, err := WriteIfNotExists("./configs/sinks.toml", `[transforms.lawg_transform]
type = "remap"
inputs = []
source = '''
message = .message
level = "info"
. = {}
.message = message
.level = level
'''

[sinks.lawg_sink]
type = "console"
encoding.codec = "json"
inputs = ["lawg_transform"]
`)

	if err != nil {
		return err
	}

	finalPath, err := WriteToPath(fmt.Sprintf("./configs/%s.toml", sourceName), fmt.Sprintf(`%s

[transforms.%s-transform]
type = "remap"
inputs = ["%s"]
source = '''
.project_namespace = "%s"
.feed_name = "%s"
'''
`, strings.Replace(config, "source0", sourceName, 1), sourceName, sourceName, projectNamespace, feedName))

	if err != nil {
		return err
	}

	if err := AddSourceToSink(sourceName); err != nil {
		return err
	}

	fmt.Println("Config generated and saved to " + finalPath)

	return nil
}

func AddSourceToSink(sourceName string) error {
	content, err := GetFileContents("./configs/sinks.toml")
	if err != nil {
		return err
	}

	lines := strings.Split(content, "\n")
	inputsLine := lines[2]
	inputsLine = inputsLine[:len(inputsLine)-1]

	println(inputsLine)
	if inputsLine == "inputs = [" {
		lines[2] = fmt.Sprintf("%s\"%s-transform\"]", inputsLine, sourceName)
	} else {
		println("meow")
		lines[2] = fmt.Sprintf("%s, \"%s-transform\"]", inputsLine, sourceName)
	}

	_, err = WriteToPath("/configs/sinks.toml", strings.Join(lines, "\n"))

	if err != nil {
		return err
	}

	return nil
}
