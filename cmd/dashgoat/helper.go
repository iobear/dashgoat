package main

import (
	"strings"
	"time"
)

//indexOf - does the value exist, and where
func indexOf(slice []string, item string) int {

	for i := range slice {

		if slice[i] == item {
			return i
		}

	}

	return -1
}

//add2url add url path root url
func add2url(path string, route string) string {
	var result strings.Builder

	if path == "/" {
		path = ""
	}

	result.WriteString(path)
	result.WriteString(route)

	return result.String()
}

//validate and enrich input from POST
func validateUpdate() bool {

	if fromPost.UpdateKey == updatekey {
		fromPost.UpdateKey = "valid"
	} else {
		return false
	}

	fromPost.Seen = time.Now().Unix()

	msglength := len(fromPost.Message)
	if msglength > 254 {
		fromPost.Message = string(fromPost.Message[0:254])
	}

	severitylen := len(fromPost.Severity)
	if severitylen > 10 {
		fromPost.Severity = string(fromPost.Severity[0:10])
	}
	fromPost.Severity = strings.ToLower(fromPost.Severity)

	statuslen := len(fromPost.Status)
	if statuslen > 10 {
		fromPost.Status = string(fromPost.Status[0:10])
	}
	fromPost.Status = strings.ToLower(fromPost.Status)

	if fromPost.Severity == "" {

		if fromPost.Status == "ok" {
			fromPost.Severity = "info"

		} else {
			fromPost.Severity = "error"

		}

	}

	return true
}
