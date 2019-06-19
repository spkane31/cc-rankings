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
	if false {
		debug = false
		fmt.Println("Testing Connections for now!")
		db := ConnectToPSQL()
		directories := []string{
			"NCAADivisionIWestRegionCrossCountryChampionship",
			"NCAADivisionISouthRegionCrossCountryChampionship",
			"NCAADivisionISoutheastRegionCrossCountryChampionship",
			"NCAADivisionISouthCentralRegionCrossCountryChampionship",
			"NCAADivisionINortheastRegionCrossCountryChampionship",
			"NCAADivisionIMountainRegionCrossCountryChampionship",
			"NCAADivisionIMidwestRegionCrossCountryChampionship",
			"NCAADivisionIMidAtlanticRegionCrossCountryChampionship",
			"NCAADivisionIGreatLakesRegionCrossCountryChampionship",
			"NCAADICrossCountryChampionship",
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
	}


	log.Println("Scraping TFRRS!")
	os.MkdirAll("RaceResults", os.ModePerm)
	HomePath, err := os.Getwd()
	HomePath = HomePath + "/RaceResults/"
	check(err)

	// links = []string{
	// 	"/results/xc/14624/Monmouth_University_XC_Tune-Up",
	// }

	GetUrlMonthYear(9, 2018)
	// log.Println(links)
	log.Printf("Found %d Links!", len(links))

	// ScrapePage(links[0])
	// os.Exit(1)

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