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
		// "NCAADivisionIWestRegionCrossCountryChampionships",
		// "NCAADivisionISouthRegionCrossCountryChampionships",
		// "NCAADivisionISoutheastRegionCrossCountryChampionships",
		// "NCAADivisionISouthCentralRegionCrossCountryChampionships",
		// "NCAADivisionINortheastRegionCrossCountryChampionships",
		// "NCAADivisionIMountainRegionCrossCountryChampionships",
		// "NCAADivisionIMidwestRegionCrossCountryChampionships",
		// "NCAADivisionIMid-AtlanticRegionCrossCountryChampionships",
		// "NCAADivisionIGreatLakesRegionCrossCountryChampionships",
		"NCAADIVISIONICROSSCOUNTRYCHAMPIONSHIPS",
	}

	// For analysis of speed. Interested in how many values can be loaded/second
	count := 0
	start := time.Now()

	// Iterate through each directory and then see how many years are in each directory
	for _, dir := range directories {
		fmt.Println(dir)
		for i := 1; ; i++ {
			files, err := ioutil.ReadDir("RaceResults2/" + dir + "/")
			check(err)
			for _, f := range files {
				var index int = 1
				for {
					fmt.Printf("\n%v\n", f.Name())
					plan, _ := ioutil.ReadFile("RaceResults2/" + dir + "/" + f.Name() + "/raceSummary.json")
					var data map[string]interface{}
					err = json.Unmarshal(plan, &data)
					f_name := fmt.Sprintf("file%v", index)
					index++
					
					// The files are simply labeled as "file1", "file2", etc. Once we can't find one, we break 
					// and go to a new directory.
					if !KeyExists(data, f_name) {break}

					var m map[string]interface{}
					m = data[f_name].(map[string]interface{})

					file_name := fmt.Sprintf("%v", m["file"])
					csvFile, err := os.Open("RaceResults2/" + dir + "/" + f.Name() + fmt.Sprintf("/%v", file_name))
					check(err)

					reader := csv.NewReader(bufio.NewReader(csvFile))
					distance := fmt.Sprintf("%v", m["distance"])
					gender := fmt.Sprintf("%v", m["gender"])
					course := fmt.Sprintf("%v", m["course"])
					date := fmt.Sprintf("%v", data["date"])
					race_name := fmt.Sprintf("%v", data["name"])
					for {
						line, err := reader.Read()
						if err == io.EOF {
							break
						} else {
							check(err)
							result_id := CreateResult(db, line, distance, gender, course, date, race_name)
							fmt.Println(result_id, line)
							os.Exit(1)
							count++
						}
					}
					fmt.Println(count)
				}
			}

			fmt.Println(count)

			os.Exit(1)
	
			plan, _ := ioutil.ReadFile("RaceResults2/" + dir + "/raceSummary.json")
			var data map[string]interface{}
			err = json.Unmarshal(plan, &data)
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