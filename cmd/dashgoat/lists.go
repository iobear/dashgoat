package main

import "fmt"

//listServiceIDs - list all reported services
func listServiceIDs() []string {

	var resultList = []string{}

	for key := range serviceStateList {
		resultList = append(resultList, key)
	}

	return resultList
}

//uniqList - search struct item, return a list of unique values
func uniqList(item string) []string {

	var resultList = []string{}
	var resultStr = ""
	var listOfServices = listServiceIDs()

	for _, id := range listOfServices {

		if item == "service" {
			resultStr = serviceStateList[id].Service

		} else if item == "host" {
			resultStr = serviceStateList[id].Host

		} else if item == "status" {
			resultStr = serviceStateList[id].Status

		} else if item == "message" {
			resultStr = serviceStateList[id].Message

		} else if item == "severity" {
			resultStr = serviceStateList[id].Severity

		} else if item == "nextupdatesec" {
			resultInt := serviceStateList[id].NextUpdateSec
			resultStr = fmt.Sprintf("%d", resultInt)

		} else if item == "seen" {
			int64Unix := serviceStateList[id].Seen
			resultStr = fmt.Sprintf("%d", int64Unix)

		} else if item == "change" {
			int64Unix := serviceStateList[id].Change
			resultStr = fmt.Sprintf("%d", int64Unix)
		}

		if item == "tags" {
			for _, value := range serviceStateList[id].Tags {
				if indexOf(resultList, value) == -1 {
					resultList = append(resultList, value)
				}
			}

		} else {
			if indexOf(resultList, resultStr) == -1 {
				resultList = append(resultList, resultStr)
			}
		}
	}

	return resultList
}
