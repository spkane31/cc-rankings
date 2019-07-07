package main

import (
	"rankings" 
	
	"fmt"
	"log"
	"encoding/json"
	"encoding/csv"
	"os"
	"bufio"
	"io"
	"io/ioutil"
	"time"
	"math"
)

var _, _, _ = os.Exit, math.Sqrt, log.Fatal

func check(e error) {
	if e != nil {
		panic(e)
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
	// completed_dir := "/home/sean/github/cc-rankings/scraper/Completed/"
	race_sum := "raceSummary.json"
	directories, err := ioutil.ReadDir(results_dir)
	check(err)

	count := 0
	start := time.Now()
	var no_hiccups bool
	// var entire_dir_added bool

	f, err := os.OpenFile("app.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	check(err)

	log.SetOutput(f)
	log.Println("This is a test log entry")

	for _, dir := range directories {
		files, err := ioutil.ReadDir(results_dir + dir.Name() + "/")
		check(err)
		no_hiccups = true
		// entire_dir_added = true
		for _, f := range files {
			var index int = 1
			json_file := results_dir + dir.Name() + "/" + f.Name() + "/" + race_sum
			plan, _ := ioutil.ReadFile(json_file)

			var data map[string]interface{}
			err = json.Unmarshal(plan, &data)
			check(err)
			
			for {
				f_name := fmt.Sprintf("file%v", index)
				index++

				if !KeyExists(data, f_name) {break}

				var m map[string]interface{}
				m = data[f_name].(map[string]interface{})
				added := m["added_to_db"].(bool)
				// if added {added = false}
				if !added {
					file_name := fmt.Sprintf("%v", m["file"])
					csvFile, err := os.Open(results_dir + dir.Name() + "/" + f.Name() + fmt.Sprintf("/%v", file_name))
					check(err)

					reader := csv.NewReader(bufio.NewReader(csvFile))
					distance := fmt.Sprintf("%v", m["distance"])
					gender := fmt.Sprintf("%v", m["gender"])
					course := fmt.Sprintf("%v", data["course"])
					date := fmt.Sprintf("%v", data["date"])
					race_name := fmt.Sprintf("%v", data["name"])
					place := 1
					// _, _, _, _, _, _, _ = db, distance, gender, course, date, race_name, place
					n, mean, variance := 0, 0.0, 0.0
					
					if distance == "N/A" || gender == "N/A" {
						log.Printf("Skipping Race. Distance = %v. Race: %v. Gender = %v\n", distance, race_name, gender)
						no_hiccups = false
						// entire_dir_added = false
						break
					} 

					for {
						line, err := reader.Read()
						if err == io.EOF {
							break
						} else {
							check(err)
							if len(line) <= 4 {
								fmt.Println("ERROR: Not correct line length: ", line)
								no_hiccups = false
								break
							} else {
								n, mean, variance = UpdateStats(n, mean, variance, rankings.GetTime(line[4]))
																								
								// line is of the format: last, first, year, team, time
								runner, result, race_id := rankings.CreateResult(db, line, distance, gender, course, date, race_name, place)
								all_results := rankings.FindResultsForRunner(db, runner)
								
								rankings.UpdateAverage(db, race_id, mean)
								rankings.UpdateStdDev(db, race_id, variance)
								if len(*all_results) > 1 {
									// This runner has multiple results, go through these and add to the graph
									rankings.AddToGraph(db, all_results, result)
								}
								place++
								count++
							} 
						}
					}

					fmt.Printf("%v results in %s seconds.\t", count, time.Now().Sub(start))
					fmt.Printf("%v results per second.\n", math.Round(float64(count) / time.Now().Sub(start).Seconds()))

					if no_hiccups {
						m["added_to_db"] = true
						data[f_name] = m
					}
					
				}
			}

			output, err := json.MarshalIndent(&data, "", "\t")
			check(err)

			err = ioutil.WriteFile(json_file, output, 0644)
			check(err)
		}

		fmt.Printf("Finished %v\n", dir.Name())

	}

	g := rankings.BuildGraph(db)
	// fmt.Println(g.Length())

	rankings.FindCorrections(g)
	os.Exit(1)

	// ComputeAverages(db)
	// UpdateCorrectionAvg(db)
	// os.Exit(1)

	// g := make(map[Pair]*Edge)

	// FindAllConnections(db, &g)

	// PlotAllRaces(db)

	fmt.Printf("%v: Finished!\n", time.Now().Format("01-02-2006, 15:04:05"))
}

func UpdateStats(n int, mean, S, new float64) (int, float64, float64) {
	// fmt.Println(new)

	prev_mean := mean
	mean = mean + (new - mean) / (float64(n+1))
	S = S + (new - mean) * (new - prev_mean)

	return n+1, mean, S

}