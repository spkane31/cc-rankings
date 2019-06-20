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

var _, _ = csv.NewReader, bufio.NewReader

// var count int
var links []string

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func KeyExists(decoded map[string]interface{}, key string) bool {
	val, ok := decoded[key]
	return ok && val != nil
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
		"NCAADivisionIMid-AtlanticRegionCrossCountryChampionships",
		"NCAADivisionIGreatLakesRegionCrossCountryChampionships",
		"NCAADICrossCountryChampionships",
		// "BigSkyConferenceChampionship",
		// "2018SunBeltConferenceCrossCountryChampionship",
		// "AmericanAthleticConferenceChampionship",
		// "Big12CrossCountryChampionship",
		// "Big8ConferenceChampionshi",
		// "BIGEASTCrossCountryChampionship",

	}

	// dir := "RaceResults/NCAADivisionIWestRegionCrossCountryChampionships/"
	count := 0
	start := time.Now()
	for _, dir := range directories {
		fmt.Println(dir)
		for i := 1; ; i++ {
	
			plan, _ := ioutil.ReadFile("RaceResults2/" + dir + "/raceSummary.json")
			var data map[string]interface{}
			err := json.Unmarshal(plan, &data)
			// fmt.Println(data)

			f_name := fmt.Sprintf("file%v", i)

			if !KeyExists(data, f_name) {break}

			var m map[string]interface{}
			m = data[f_name].(map[string]interface{})

			file_name := fmt.Sprintf("%v",m["file"])
			csvFile, err := os.Open("RaceResults2/" + dir + fmt.Sprintf("/%v", file_name))
			check(err)
			reader := csv.NewReader(bufio.NewReader(csvFile))
			distance := fmt.Sprintf("%v", m["distance"])
			gender := fmt.Sprintf("%v", m["gender"])
			// if gender == "WOMENS" {break}
			course := fmt.Sprintf("%v", data["course"])
			date := fmt.Sprintf("%v", data["date"])
			race_name := fmt.Sprintf("%v", data["name"])
			for {
				line, err := reader.Read()
				if err == io.EOF {
					break
				} else {
					check(err)
					// fmt.Println(line)
					// fmt.Println(line)
					// break
					CreateResult(db, line, distance, gender, course, date, race_name)
					count++
					// os.Exit(1)
				}
			}
			// os.Exit(1)
		}

	}
	elapsed := time.Since(start)
	seconds := elapsed.Seconds()
	entries_second := float64(count) / seconds
	fmt.Printf("Took %v nanoseconds to insert %v entries. Average %v per second.\n", seconds, count, int64(entries_second))
	os.Exit(1)
}