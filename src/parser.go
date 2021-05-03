package main

import (
	"encoding/xml"
	"log"
	"os"
	"strings"
)

type Result struct {
}

func parseDoc(inOut chan string, complete chan struct{}) {

	input :=<- inOut
	xmlFile, err := os.Open(resultName)
	if err != nil {
		log.Fatal(err)
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	decoder.Strict = false

	var hiragana string
	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}
		switch startElement := token.(type) {
		case xml.StartElement:
			if startElement.Name.Local == "reb"{
				token, _ := decoder.Token()
				switch token.(type) {
				case xml.CharData:
					hiragana = string([]byte(token.(xml.CharData)))
				}
			}
			if startElement.Name.Local == "gloss" {
				token, _ := decoder.Token()
				switch token.(type) {
				case xml.CharData:
					def := string([]byte(token.(xml.CharData)))
					if strings.Contains(strings.ToLower(def), input){
						select {
						case <-complete:
							return
						case inOut <- hiragana + "\n" + def:
						}
					}
				}
			}
		}
	}
}
