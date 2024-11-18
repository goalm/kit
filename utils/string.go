package utils

import (
	"fmt"
	"regexp"
)

func FilePathToName(str string) (string, error) {
	re := regexp.MustCompile(`[\/|\\]*(\w*)\.\w*$`)
	match := re.FindStringSubmatch(str)
	if len(match) > 0 {
		return match[1], nil
	}
	return "", fmt.Errorf("No match found")
}
