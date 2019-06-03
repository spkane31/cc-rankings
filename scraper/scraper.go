package main

import (
	"fmt"
	"net/http"
	"log"
	"strings"
	"io/ioutil"
)

func main() {
	fmt.Println("Scraping TFRRS!")
	response, err := http.Get("https://www.tfrrs.org/")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response)

	dataInBytes, err := ioutil.ReadAll(response.Body)
	pageContent := string(dataInBytes)

	titleStartIndex := strings.Index(pageContent, "<title>")
	fmt.Println(titleStartIndex)
}