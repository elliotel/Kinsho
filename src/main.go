package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	archivePath = "dictionary/JMdict_e.gz"
	dictPath    = "dictionary/JMdict_e"
	dictName    = "JMdict_e"
)

func main() {
	f := app.New()
	icon, err := fyne.LoadResourceFromPath("img/kinsho_icon.png")
	if err != nil {
		log.Fatal()
	}
	f.SetIcon(icon)

	if _, err := os.Stat(dictPath + "_0"); os.IsNotExist(err) {
		if !connected() {
			displayConnectionError(f).ShowAndRun()
		} else {
			downloadJMdict()
			decompressAndDeleteGZ(archivePath)
			splitXML()
		}
	}
	complete := make(chan struct{})
	inputChan := make(chan string)
	outputChan := make(chan entry)
	displayGUI(f, inputChan, outputChan, complete)
}

func connected() bool {
	_, err := http.Get("http://ftp.edrdg.org/pub/Nihongo/JMdict_e.gz")
	if err != nil {
		return false
	}
	return true
}

func contains(array []string, s string) bool {
	for _, val := range array {
		if strings.Contains(strings.ToLower(val), strings.ToLower(s)) {
			return true
		}
	}
	return false
}
