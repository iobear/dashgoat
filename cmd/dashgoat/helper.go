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
func (ss *ServiceState) validateUpdate() bool {

	if ss.UpdateKey == updatekey {
		ss.UpdateKey = "valid"
	} else {
		return false
	}

	if ss.Probe == 0 {
		ss.Probe = time.Now().Unix()
	}

	msglength := len(ss.Message)
	if msglength > 254 {
		ss.Message = string(ss.Message[0:254])
	}

	severitylen := len(ss.Severity)
	if severitylen > 10 {
		ss.Severity = string(ss.Severity[0:10])
	}
	ss.Severity = strings.ToLower(ss.Severity)

	statuslen := len(ss.Status)
	if statuslen > 10 {
		ss.Status = string(ss.Status[0:10])
	}
	ss.Status = strings.ToLower(ss.Status)

	if ss.Severity == "" {

		if ss.Status == "ok" {
			ss.Severity = "info"

		} else {
			ss.Severity = "error"

		}

	}

	return true
}
