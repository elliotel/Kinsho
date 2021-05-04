package main

import (
	"encoding/xml"
	"github.com/gojp/kana"
	"log"
	"os"
	"strings"
)

const (
	entryAmount = 20
)

type entry struct {
	kanji string
	kana string
	def string
}

func parseDoc(inputChan chan string, outputChan chan entry, complete chan struct{}) {

	input := <-inputChan
	inputHiragana := kana.RomajiToHiragana(input)
	inputKatakana := kana.RomajiToKatakana(input)
	xmlFile, err := os.Open(resultName)
	if err != nil {
		log.Fatal(err)
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	decoder.Strict = false



	entries := make([]entry, entryAmount)
	i := 0

	for i < entryAmount {
		token, _ := decoder.Token()
		if token == nil {
			break
		}
		switch element := token.(type) {
		case xml.StartElement:
			startElement := element.Name.Local
			token, _ := decoder.Token()
			switch startElement {
			case "keb":
				entries[i].kanji = string(token.(xml.CharData))
			case "reb":
				entries[i].kana = string(token.(xml.CharData))
			case "gloss":
				entries[i].def = string(token.(xml.CharData))
			}
		case xml.EndElement:
			if element.Name.Local == "entry" {
				if strings.Contains(strings.ToLower(entries[i].def), input) ||
					strings.Contains(strings.ToLower(entries[i].kana), input) ||
					strings.Contains(strings.ToLower(entries[i].kanji), input) ||
					strings.Contains(strings.ToLower(entries[i].kana), inputHiragana) ||
					strings.Contains(strings.ToLower(entries[i].def), inputKatakana) {
					outputChan <- entries[i]
					i++
				}
			}
		}

	}
	complete <- struct{}{}
}
