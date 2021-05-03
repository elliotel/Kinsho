package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func parseDoc() {
	file, err := os.Open("edict2.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	searchTerm := "stockholm"

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(strings.ToLower(scanner.Text()), strings.ToLower(searchTerm)) {
			fmt.Println(scanner.Text())
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	}
