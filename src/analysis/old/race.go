package main

import (
	"database/sql"
	"fmt"
	"os"
	"math"
	"math/rand"
	"time"
	"encoding/csv"
	// "bufio"
	// "io"

	_ "github.com/lib/pq"
)

var _, _, _ = fmt.Println, os.Exit, csv.NewReader

type Race struct {
	id int
	name string
	course string
	distance int
	gender string
	is_base bool
	average float64
	std_dev float64
	correction_avg float64
	correction_graph float64
}

type Instance struct {
	id int
	date string
	race_id int 
	average float64
	std_dev float64
}

func UpdateAverage(db *sql.DB, id int, average float64) {
	update := `UPDATE races SET average=$2 WHERE id=$1;`
	_, err := db.Exec(update, id, average)
	check(err)
}
func UpdateStdDev(db *sql.DB, id int, std_dev float64) {
	update := `UPDATE races SET std_dev=$2 WHERE id=$1;`
	_, err := db.Exec(update, id, std_dev)
	check(err)
}

func UpdateRace(db *sql.DB, id int, correction float64) {
	update := `UPDATE races SET correction=$2 WHERE id=$1;`
	_, err := db.Exec(update, id, correction)
	check(err)
}

func ComputeAverages(db *sql.DB) {
	fmt.Printf("%v: Updating Race Averages\n", time.Now().Format("01-02-2006, 15:04:05"))
	races := GetAllRacesByGender(db, "MENS")
	for i := range *races {
		insts := FindAllInstances(db, (*races)[i].id)
		var results []Result
		for i := range *insts {
			r := GetInstanceResults(db, (*insts)[i].id)
			for _, each := range *r {
				results = append(results, each)
			}
		}

		results = FilterDNFs(&results)
		mean, std_dev := CalculateStatistics(&results)

		// fmt.Printf("Race ID: %v\tMean: %v\t St. Dev: %v\n", (*races)[i].id, mean, std_dev)
		UpdateAverage(db, (*races)[i].id, mean)
		UpdateStdDev(db, (*races)[i].id, std_dev)
	}
}

func UpdateCorrectionAvg(db * sql.DB) {
	fmt.Printf("%v: Updating Race correction_avg\n", time.Now().Format("01-02-2006, 15:04:05"))

	races := GetAllRacesByGender(db, "MENS")
	base := rand.Intn(len(*races))
	SetAsBase(db, base)
	fmt.Printf("\tRace %d is the base race.\n", base)

	race := FindRace(db, base)

	var temp Race
	for i := range *races {
		if i != base && (*races)[i].gender == "MENS" {
			temp = FindRace(db, (*races)[i].id)
			dif := race.average - temp.average
			UpdateRaceCorrectionAvg(db, (*races)[i].id, dif)
			fmt.Printf("Race: %v\tCorrection: %v\n", (*races)[i].name, dif)
		}
	}

}

func UpdateRaceCorrectionAvg(db *sql.DB, id int, dif float64) {
	update := `UPDATE races SET correction_avg=$2 WHERE id=$1;`
	_, err := db.Exec(update, id, dif)
	check(err)
}

func SetAsBase(db *sql.DB, id int) {
	update := `UPDATE races SET is_base=$2 WHERE id=$1;`
	_, err := db.Exec(update, id, true)
	check(err)
}

func CalculateStatistics(results *[]Result) (float64, float64) {
	size := len(*results)
	sum_x, sum_y, sum_xx, sum_xy := 0.0, 0.0, 0.0, 0.0

	mean := 0.0
	S := 0.0
	
	var t float64

	for i := range *results {
		sum_x += float64(i+1)
		sum_xx += float64((i+1) * (i+1))

		t = GetTime((*results)[i].time)
		
		sum_y += t
		sum_xy += (float64(i+1) * t)
		prev_mean := mean
		mean = mean + (t - mean) / float64(i+1)

		S = S + (t - mean) * (t - prev_mean)
	}

	return mean, math.Sqrt(S / float64(size))
}

func FindRace(db *sql.DB, id int) Race {
	var ret Race
	query := `SELECT * FROM races WHERE id=$1;`
	row := db.QueryRow(query, id)
	err := row.Scan(&ret.id, &ret.name, &ret.course, &ret.distance, &ret.gender, &ret.is_base, &ret.average, &ret.std_dev, &ret.correction_avg, &ret.correction_graph)
	check(err)
	return ret
}

func FindRaceByCourseName(db *sql.DB, course, name, gender string) Race {
	var ret Race
	query := `SELECT id, name, course, distance, gender FROM races WHERE (course=$1 AND name=$2 AND gender=$3);`
	row := db.QueryRow(query, course, name, gender)
	err := row.Scan(&ret.id, &ret.name, &ret.course, &ret.distance, &ret.gender)
	check(err)
	return ret
}

func FindAllInstances(db *sql.DB, id int) *[]Instance {
	var ret []Instance

	queryStatement := `SELECT * from race_instances WHERE race_id=$1;`
	rows, err := db.Query(queryStatement, id)
	check(err)
	defer rows.Close()

	for rows.Next() {
		var inst Instance
		err = rows.Scan(&inst.id, &inst.date, &inst.race_id, &inst.average, &inst.std_dev)
		ret = append(ret, inst)
	}

	return &ret
}


func GetInstanceResults(db *sql.DB, id int) *[]Result {
	var ret []Result

	queryStatement := `SELECT * FROM results WHERE race_instance_id=$1;`
	rows, err := db.Query(queryStatement, id)
	check(err)
	defer rows.Close()

	for rows.Next() {
		var result Result
		err = rows.Scan(&result.id, &result.distance, &result.unit, &result.rating, &result.time, &result.race_instance_id, &result.runner_id, &result.scaled_time, &result.time_float)
		check(err)
		ret = append(ret, result)
	}

	return &ret
}

func GetAllRacesByGender(db *sql.DB, gender string) *[]Race {
	var ret []Race

	queryStatement := `SELECT id, name, course, distance, gender FROM races WHERE gender=$1;`
	rows, err := db.Query(queryStatement, gender)
	check(err)
	defer rows.Close()

	for rows.Next() {
		var race Race
		err = rows.Scan(&race.id, &race.name, &race.course, &race.distance, &race.gender)
		ret = append(ret, race)
	}

	return &ret
}