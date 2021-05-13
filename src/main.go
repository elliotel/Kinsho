package main

import (
	"os"
	"strings"
)

const (
	archiveName = "JMdict_e.gz"
	resultName  = "JMdict_e"
)

func main() {

	if _, err := os.Stat(resultName); os.IsNotExist(err) {
		downloadJMdict()
		decompressAndDeleteGZ(archiveName)
	}
	complete := make(chan struct{})
	inputChan := make(chan string)
	outputChan := make(chan entry)
	displayGUI(inputChan, outputChan, complete)
}

func contains(array []string, s string) bool {
	for _, val := range array {
		if strings.Contains(strings.ToLower(val), strings.ToLower(s)) {
			return true
		}
	}
	return false
}

func containsExact(array []string, s string) bool {
	for _, val := range array {
		if strings.ToLower(val) == strings.ToLower(s) {
			return true
		}
	}
	return false
}
