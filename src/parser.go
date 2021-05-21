package main

import (
	"encoding/xml"
	"github.com/gojp/kana"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
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

	var wg sync.WaitGroup

	sectionChan := make(chan entry)

	input := <-inputChan
	if strings.TrimSpace(input) == "" {
		complete <- struct{}{}
		return
	}
	inputHiragana := kana.RomajiToHiragana(input)
	inputKatakana := kana.RomajiToKatakana(input)

	entries := make([]entry, 0)

	files, err := ioutil.ReadDir("dictionary/")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.Contains(file.Name(), dictName) {
			wg.Add(1)
			xmlFile, err := os.Open("dictionary/" + file.Name())
			if err != nil {
				log.Fatal(err)
			}
			go parseSection(input, inputHiragana, inputKatakana, xmlFile, sectionChan, &wg)
		}
	}

	go func() {
		wg.Wait()
		close(sectionChan)
	}()

	for entry := range sectionChan {
		entries = append(entries, entry)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].priority > entries[j].priority
	})
	if len(entries) > 0 {
		for i := 0; i < len(entries) && i < entryAmount; i++ {
			outputChan <- entries[i]
		}
	}
	complete <- struct{}{}
}

func parseSection(input string, inputHiragana string, inputKatakana string, xmlFile *os.File, sectionChan chan entry, wg *sync.WaitGroup) {
	decoder := xml.NewDecoder(xmlFile)
	decoder.Strict = false

	var current entry
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
				current.kanji = append(current.kanji, kanji)
				if input == kanji {
					current.priority += 100
				}

			case "reb":
				ganaKana := string(token.(xml.CharData))
				current.kana = append(current.kana, ganaKana)
				if input == ganaKana {
					current.priority += 100
				}
			case "gloss":
				def := string(token.(xml.CharData))
				current.def = append(current.def, def)
				if input == strings.ToLower(def) {
					current.priority += 100
				}
			case "ke_pri":
				priority := string(token.(xml.CharData))
				switch priority {
				case "ichi1":
					current.priority += 10
				case "ichi2":
					current.priority += 5
				case "spec1":
					current.priority += 30
				case "spec2":
					current.priority += 10
				}
				runes := []rune(priority)
				if string(runes[0])+string(runes[1]) == "nf" {
					number, err := strconv.Atoi(string(runes[2]) + string(runes[3]))
					if err != nil {
						log.Fatal(err)
					}
					current.priority += 49 - number
				}

			}

		case xml.EndElement:
			if element.Name.Local == "entry" {
				if contains(current.def, input) ||
					contains(current.kana, input) ||
					contains(current.kanji, input) ||
					contains(current.kana, inputHiragana) ||
					contains(current.kana, inputKatakana) {
					sectionChan <- current
				}
				current = entry{}
			}
		}
	}
	wg.Done()
	err := xmlFile.Close()
	if err != nil {
		log.Fatal(err)
	}
}
