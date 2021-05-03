package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
)

type Result struct {
}

func parseDoc() {
	xmlFile, err := os.Open("JMdict_e")
	if err != nil {
		log.Fatal(err)
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	decoder.Strict = false
	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}

		switch startElement := token.(type) {
		case xml.StartElement:
			if startElement.Name.Local == "gloss" {
				token, _ := decoder.Token()
				switch token.(type) {
				case xml.CharData:

					fmt.Println(string([]byte(token.(xml.CharData))))
				}
			}
		}
	}
}
