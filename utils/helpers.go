package utils

import "strings"

func ParseVectorConfig(config []byte) string {
	stringedVectorGeneratedConfig := string(config[:])

	lines := strings.Split(stringedVectorGeneratedConfig, "\n")
	lines = lines[2:] // remove first 2 lines of generatedConfig (irrelevant to config)

	return strings.Trim(strings.Join(lines, "\n"), "\n")
}
