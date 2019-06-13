package main

import (
	"fmt"
	"log"
	"encoding/json"
	"encoding/csv"
	"os"
	"bufio"
	"io"
	"io/ioutil"
	"time"

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
var debug bool

func main() {
	debug = false
	fmt.Println("Testing Connections for now!")
	db := ConnectToPSQL()
	directories := []string{
		"NCAADivisionIWestRegionCrossCountryChampionships",
		"NCAADivisionISouthRegionCrossCountryChampionships",
		"NCAADivisionISoutheastRegionCrossCountryChampionships",
		"NCAADivisionISouthCentralRegionCrossCountryChampionships",
		"NCAADivisionINortheastRegionCrossCountryChampionships",
		"NCAADivisionIMountainRegionCrossCountryChampionships",
		"NCAADivisionIMidwestRegionCrossCountryChampionships",
		"NCAADivisionIMidAtlanticRegionCrossCountryChampionships",
		"NCAADICrossCountryChampionships",
	}

	// dir := "RaceResults/NCAADivisionIWestRegionCrossCountryChampionships/"
	count := 0
	start := time.Now()
	for _, dir := range directories {
		fmt.Println(dir)

		csvFile, err := os.Open("RaceResults/" + dir + "/mens.csv")
		check(err)
	
		plan, _ := ioutil.ReadFile("RaceResults/" + dir + "/raceSummary.json")
		var data map[string]string
		err = json.Unmarshal(plan, &data)
		fmt.Println(data["name"])
	
		reader := csv.NewReader(bufio.NewReader(csvFile))
		
		for {
			line, err := reader.Read()
			if err == io.EOF {
				break
			} else {
				check(err)
				// fmt.Println(line)
	
				CreateResult(db, "MALE", line, data)
				// os.Exit(1)
			}
			count++
		}
	}
	elapsed := time.Since(start)
	seconds := elapsed.Seconds()
	entries_second := float64(count) / seconds
	fmt.Printf("Took %v nanoseconds to insert %v entries. Average %v per second.\n", seconds, count, int64(entries_second))
	os.Exit(1)

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