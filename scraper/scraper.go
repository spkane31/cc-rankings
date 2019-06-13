package main

import (
	"fmt"
	// "net/http"
	"log"
	// "net/url"
	// "strconv"
	// "strings"
	// "unicode"
	"encoding/json"
	"encoding/csv"
	"os"
	"bufio"
	"io"
	"io/ioutil"

	_ "github.com/PuerkitoBio/goquery"
)

// var count int
var links []string

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
var HomePath string

func main() {
	fmt.Println("Testing Connections for now!")
	db := ConnectToPSQL()

	csvFile, err := os.Open("/home/sean/github/cc-rankings/scraper/RaceResults/NCAADICrossCountryChampionships/mens.csv")

	plan, _ := ioutil.ReadFile("RaceResults/NCAADICrossCountryChampionships/raceSummary.json")
	var data map[string]string
	err = json.Unmarshal(plan, &data)
	fmt.Println(data["name"])

	reader := csv.NewReader(bufio.NewReader(csvFile))

	for i := 0; i < 10; i++ {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else {
			check(err)
			fmt.Println(line)

			CreateResult(db, "MALE", line, data)
			// os.Exit(1)
		}
	}

	// CreateRunner(db, "sean", "kane")//, "UC", "SR")
	// InsertRunner(db, "nick", "dehaven", "UC", "JR")
	// InsertRunner(db, "evan", "sargent", "OSU", "FR")

	// QueryRunnerFromID(db, 1)
	// QueryRunnerFromID(db, 2)
	// QueryRunnerFromID(db, 3)

	// r := GetRunners(db, 2)
	// fmt.Println(r)


	os.Exit(1)
	log.Println("Scraping TFRRS!")
	os.MkdirAll("RaceResults", os.ModePerm)
	HomePath, err := os.Getwd()
	HomePath = HomePath + "/RaceResults/"
	check(err)

	GetUrlMonthYear(11, 2018)
	log.Println(links)
	log.Printf("Found %d Links!", len(links))

	// Invoke as goroutines to maximize efficiency
	for i := 0; i < len(links); i += 4 {
		if i < len(links) {ScrapePage(links[i])}
		if i+1 < len(links) {
			go ScrapePage(links[i+1])
		}
		if i+2 < len(links) {
			go ScrapePage(links[i+2])
		}
		if i+3 < len(links) {
			go ScrapePage(links[i+3])
		}
	}

	log.Println("Finished!")
}