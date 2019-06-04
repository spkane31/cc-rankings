package main

import (
	// "fmt"
	"net/http"
	"log"
	"net/url"
	"strconv"
	"strings"
	// "io/ioutil"

	"github.com/PuerkitoBio/goquery"
)

// var count int
var links []string

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func processElement(index int, element *goquery.Selection) {
	href, exists := element.Attr("href")
	if exists {
		// fmt.Println(href)
		ParseUrl(href)
	}
}

func ParseUrl(s string) {
	parsedUrl, err := url.Parse(s)
	check(err)
	path := parsedUrl.Path

	if strings.Contains(path, "/results/") {
		links = append(links, parsedUrl.Path)	
	}
}

func GetUrlMonthYear(month, year int) {
	y := strconv.Itoa(year)
	m := strconv.Itoa(month)

	response, err := http.PostForm(
		"https://www.tfrrs.org/results_search.html",
		url.Values{
			"sport": {"xc"},
			"state": {""},
			"month": {m},
			"year": {y},
		},
	)

	check(err)
	defer response.Body.Close()
	document, err := goquery.NewDocumentFromReader(response.Body)
	check(err)
	document.Find("a").Each(processElement)
}

func ScrapePage(link string) {
	// This function will scrape a race page and create the CSV, and the 
}

func main() {
	
	log.Println("Scraping TFRRS!")
	GetUrlMonthYear(11, 2018)
	log.Printf("Found %d Links!", len(links))
	// log.Println(links)
	response, err := http.Get("https://www.tfrrs.org/results_search.html")
	check(err)
	defer response.Body.Close()

	// fmt.Println(response)

	// document, err := goquery.NewDocumentFromReader(response.Body)
	// check(err)

	// document.Find("a").Each(processElement)
}