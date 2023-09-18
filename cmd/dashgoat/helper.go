/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	"os"
	"strconv"
	"strings"
)

// indexOf - does the value exist, and where
func indexOf(slice []string, item string) int {

	for i := range slice {

		if slice[i] == item {
			return i
		}

	}

	return -1
}

// contains Does the value exist, https://gosamples.dev/slice-contains/
func contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

// add2url add url path root url
func add2url(path string, route string) string {
	var result strings.Builder

	if path == "/" {
		path = ""
	}

	result.WriteString(path)
	result.WriteString(route)

	return result.String()
}

// isExists Does the given directory of filepath exist?
func isExists(path string, task string) bool {
	fileStat, err := os.Stat(path)

	if err != nil {
		return false
	}

	if task == "path" || task == "directory" {
		return fileStat.IsDir()
	}

	return true
}

// Convert string value to an integer
func str2int(str_to_convert string) int {

	result, err := strconv.Atoi(str_to_convert)
	if err != nil {
		return 0
	}

	return result
}
