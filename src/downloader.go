package main

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"os"
)

func downloadJMdict() {
	fullURL := "http://ftp.edrdg.org/pub/Nihongo/JMdict_e.gz"

	file, err := os.Create(archiveName)
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
		err := os.Remove(archiveName)
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

	writer, err := os.Create(resultName)
	if err != nil {
		log.Fatal(err)
	}
	defer writer.Close()

	_, err = io.Copy(writer, gzReader)
	if err != nil {
		log.Fatal(err)
	}
}
