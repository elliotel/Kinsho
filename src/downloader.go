package main

import (
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	entriesPerFile = 1000000
)

type xmlEntry struct {
	XMLName   xml.Name `xml:"entry"`
	Kanjis    Kanjis
	GanaKanas GanaKanas
	Defs      Defs
	Prioritys Prioritys
}

type Kanjis struct {
	Kanji []Kanji
}

type GanaKanas struct {
	GanaKana []GanaKana
}

type Defs struct {
	Def []Def
}

type Prioritys struct {
	Priority []Priority
}

type Kanji struct {
	XMLName xml.Name `xml:"keb"`
	Kanji   string   `xml:",innerxml"`
}

type GanaKana struct {
	XMLName  xml.Name `xml:"reb"`
	GanaKana string   `xml:",innerxml"`
}

type Def struct {
	XMLName xml.Name `xml:"gloss"`
	Def     string   `xml:",innerxml"`
}

type Priority struct {
	XMLName  xml.Name `xml:"ke_pri"`
	Priority string   `xml:",innerxml"`
}

func downloadJMdict() {

	if _, err := os.Stat("dictionary/"); os.IsNotExist(err) {
		err := os.Mkdir("dictionary/", os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	files, err := ioutil.ReadDir("dictionary/")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.Contains(file.Name(), dictName) {
			if err := os.Remove("dictionary/" + file.Name()); err != nil {
				log.Fatal(err)
			}
		}
	}

	fullURL := "http://ftp.edrdg.org/pub/Nihongo/JMdict_e.gz"

	file, err := os.Create(archivePath)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	resp, err := client.Get(fullURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)

	defer file.Close()
}

func decompressAndDeleteGZ(fileName string) {

	defer func() {
		err := os.Remove(archivePath)
		if err != nil {
			log.Fatal(err)
		}
	}()

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		log.Fatal(err)
	}
	defer gzReader.Close()

	writer, err := os.Create(dictPath)
	if err != nil {
		log.Fatal(err)
	}
	defer writer.Close()

	_, err = io.Copy(writer, gzReader)
	if err != nil {
		log.Fatal(err)
	}
}

func splitXML() {

	var wg sync.WaitGroup

	xmlFile, err := os.Open(dictPath)
	if err != nil {
		log.Fatal(err)
	}

	decoder := xml.NewDecoder(xmlFile)
	decoder.Strict = false
	xmlEntries := make([]xmlEntry, 1)

	var current xmlEntry
	fileNumber := 0
	runs := 0
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
				kanjiString := string(token.(xml.CharData))
				current.Kanjis.Kanji = append(current.Kanjis.Kanji, Kanji{Kanji: kanjiString})

			case "reb":
				ganaKanaString := string(token.(xml.CharData))
				current.GanaKanas.GanaKana = append(current.GanaKanas.GanaKana, GanaKana{GanaKana: ganaKanaString})
			case "gloss":
				defString := html.EscapeString(string(token.(xml.CharData)))
				current.Defs.Def = append(current.Defs.Def, Def{Def: defString})
			case "ke_pri":
				priorityString := string(token.(xml.CharData))
				current.Prioritys.Priority = append(current.Prioritys.Priority, Priority{Priority: priorityString})
			}

		case xml.EndElement:
			if element.Name.Local == "entry" {
				xmlEntries = append(xmlEntries, current)
				current = xmlEntry{}
			}

		}
		runs++
		if runs%entriesPerFile == 0 {

			wg.Add(1)
			go createXml(xmlEntries, fileNumber, &wg)
			fileNumber++
			xmlEntries = make([]xmlEntry, 0)
		}
	}

	wg.Add(1)
	go createXml(xmlEntries, fileNumber, &wg)

	wg.Wait()

	xmlFile.Close()
	if err := os.Remove(dictPath); err != nil {
		log.Fatal(err)
	}
}

func createXml(xmlEntries []xmlEntry, fileNumber int, wg *sync.WaitGroup) {

	fileName := "dictionary/JMdict_e_" + strconv.FormatUint(uint64(fileNumber), 10)
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, entry := range xmlEntries {

		out, err := xml.MarshalIndent(entry, "", "   ")

		if err != nil {
			panic(err)
		}

		if _, err := file.Write(out); err != nil {
			log.Fatal(err)
		}
		if _, err := file.WriteString("\n"); err != nil {
			log.Fatal(err)
		}
	}
	file.Close()
	wg.Done()
}
