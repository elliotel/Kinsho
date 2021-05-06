package main

import (
	"encoding/xml"
	"fmt"
	"github.com/gojp/kana"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	entryAmount = 150
)

type entry struct {
	kanji string
	kana string
	def string
	priority int
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



	entries := make([]entry, 1)

	test := 0
	index := 0
	for {
		//fmt.Println(entries[test])
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
				entries[index].kanji = kanji
				if input == kanji {
					entries[index].priority += 100
				}

			case "reb":
				ganaKana := string(token.(xml.CharData))
				entries[index].kana = ganaKana
				if input == ganaKana {
					entries[index].priority += 100
				}
			case "gloss":
				def := string(token.(xml.CharData))
				entries[index].def = def
				if input == def {
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
				if string(runes[0]) + string(runes[1]) == "nf" {
					number, err := strconv.Atoi(string(runes[2]) + string(runes[3]))
					if err != nil {
						log.Fatal(err)
					}
					entries[index].priority += 49 - number
				}

			}

		case xml.EndElement:
			if element.Name.Local == "entry" {
				if strings.Contains(strings.ToLower(entries[index].def), input) ||
					strings.Contains(strings.ToLower(entries[index].kana), input) ||
					strings.Contains(strings.ToLower(entries[index].kanji), input) ||
					strings.Contains(strings.ToLower(entries[index].kana), inputHiragana) ||
					strings.Contains(strings.ToLower(entries[index].kana), inputKatakana) {
					if index == test && test != 0 {
						fmt.Println(entries[index].kanji + " added " + " with " + strconv.Itoa(entries[index].priority) + " priority")
					}
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
	if !(entries[0].kanji == "" && entries[0].kana == "" && entries[0].priority == 0) {
		for i := 0; i < len(entries) && i < entryAmount; i++ {
			fmt.Printf("Word: %s, Freq: %d\n", entries[i].kanji, entries[i].priority)
			outputChan <- entries[i]
		}
	}
	complete <- struct{}{}
}
