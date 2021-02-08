package utils

import (
	"bufio"
	"strings"
)

func Int32Ptr(i int32) *int32 { return &i }

func StringToLines(s string) (lines []string, err error) {
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err = scanner.Err()
	return lines, err
}

func Contains(source []string, value string) bool {
	for _, item := range source {
		if item == value {
			return true
		}
	}
	return false
}
