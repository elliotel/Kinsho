package main

import (
	"os"
	"strings"
)

const (
	archivePath = "dictionary/JMdict_e.gz"
	dictPath    = "dictionary/JMdict_e"
	dictName    = "JMdict_e"
)

func main() {

	if _, err := os.Stat(dictPath + "_0"); os.IsNotExist(err) {
		downloadJMdict()
		decompressAndDeleteGZ(archivePath)
		splitXML()
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
