package main

import (
	"encoding/xml"
	"github.com/gojp/kana"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	entryAmount = 20
)

type entry struct {
	kanji    []string
	kana     []string
	def      []string
	priority int
}

func parseDoc(inputChan chan string, outputChan chan entry, complete chan struct{}) {

	input := <-inputChan
	if strings.TrimSpace(input) == "" {
		complete <- struct{}{}
		return
	}
	inputHiragana := kana.RomajiToHiragana(input)
	inputKatakana := kana.RomajiToKatakana(input)
	xmlFile, err := os.Open(resultName)
	if err != nil {
		log.Fatal(err)
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	decoder.Strict = false

	entries := make([]entry, 1)

	index := 0
	for {
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
				kanji := string(token.(xml.CharData))
				entries[index].kanji = append(entries[index].kanji, kanji)
				if input == kanji {
					entries[index].priority += 100
				}

			case "reb":
				ganaKana := string(token.(xml.CharData))
				entries[index].kana = append(entries[index].kana, ganaKana)
				if input == ganaKana {
					entries[index].priority += 100
				}
			case "gloss":
				def := string(token.(xml.CharData))
				entries[index].def = append(entries[index].def, def)
				if input == strings.ToLower(def) {
					entries[index].priority += 100
				}
			case "ke_pri":
				priority := string(token.(xml.CharData))
				switch priority {
				case "ichi1":
					entries[index].priority += 10
				case "ichi2":
					entries[index].priority += 5
				case "spec1":
					entries[index].priority += 30
				case "spec2":
					entries[index].priority += 10
				}
				runes := []rune(priority)
				if string(runes[0])+string(runes[1]) == "nf" {
					number, err := strconv.Atoi(string(runes[2]) + string(runes[3]))
					if err != nil {
						log.Fatal(err)
					}
					entries[index].priority += 49 - number
				}

			}

		case xml.EndElement:
			if element.Name.Local == "entry" {
				if contains(entries[index].def, input) ||
					contains(entries[index].kana, input) ||
					contains(entries[index].kanji, input) ||
					contains(entries[index].kana, inputHiragana) ||
					contains(entries[index].kana, inputKatakana) {
					entries = append(entries, entry{})
					index++
				}
				entries[index] = entry{}
			}
		}

	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].priority > entries[j].priority
	})
	if !(entries[0].kanji == nil && entries[0].kana == nil && entries[0].priority == 0) {
		for i := 0; i < len(entries) && i < entryAmount; i++ {
			outputChan <- entries[i]
		}
	}
	complete <- struct{}{}
}
