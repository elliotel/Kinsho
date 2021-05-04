package main

import (
	"os"
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
