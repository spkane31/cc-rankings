package main

import (
	"rankings" 
	
	"fmt"
	"log"
	"encoding/json"
	// "encoding/csv"
	"os"
	// "bufio"
	// "io"
	"io/ioutil"
	"time"
)

var _ = os.Exit

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func KeyExists(decoded map[string]interface{}, key string) bool {
	val, ok := decoded[key]
	return ok && val != nil
}

func main() {
	fmt.Printf("%v: Establishing DB connection...\n", time.Now().Format("01-02-2006, 15:04:05"))
	db, err := rankings.ConnectToPSQL()
	check(err)

	results_dir := "/home/sean/github/cc-rankings/scraper/RaceResults/"
	race_sum := "raceSummary.json"
	directories, err := ioutil.ReadDir(results_dir)
	check(err)

	// count := 0
	// start := time.Now()

	for _, dir := range directories {
		files, err := ioutil.ReadDir(results_dir + dir.Name() + "/")
		check(err)

		for _, f := range files {
			var index int = 1
			for {
				plan, _ := ioutil.ReadFile(results_dir + dir.Name() + "/" + f.Name() + "/" + race_sum)
				var data map[string]interface{}
				err = json.Unmarshal(plan, &data)
				f_name := fmt.Sprintf("file%v", index)
				index++

				if !KeyExists(data, f_name) {break}

				var m map[string]interface{}
				m = data[f_name].(map[string]interface{})
				added := m["added_to_db"].(bool)
				if !added {
					file_name := fmt.Sprintf("%v", m["file"])
					_, err := os.Open(results_dir + dir.Name() + "/" + f.Name() + fmt.Sprintf("/%v", file_name))
					check(err)

					// reader := csv.NewReader(bufio.NewReader(csvFile))
					// distance := fmt.Sprintf("%v", m["distance"])
					// gender := fmt.Sprintf("%v", m["gender"])
					// course := fmt.Sprintf("%v", data["course"])
					// date := fmt.Sprintf("%v", data["date"])
					// race_name := fmt.Sprintf("%v", data["name"])

					fmt.Println(data)
					fmt.Println(m)

				}

				os.Exit(1)
			}
		}
	}

	os.Exit(1)

	ComputeAverages(db)
	UpdateCorrectionAvg(db)
	os.Exit(1)

	g := make(map[Pair]*Edge)

	FindAllConnections(db, &g)

	PlotAllRaces(db)

	fmt.Printf("%v: Finished!\n", time.Now().Format("01-02-2006, 15:04:05"))
}