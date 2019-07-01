package rankings

import (
	"database/sql"
	"fmt"
	"os"
	"math"
	// "math/rand"
	"time"
	// "encoding/csv"
	// "bufio"
	// "io"

	_ "github.com/lib/pq"
)

type Race struct {
	id int
	name string
	course string
	distance int
	gender string
	// inserted_at
	// updated_at 
	is_base bool
	average float64
	std_dev float64
	correction_avg float64
	correction_graph float64
}

func AddRace(db *sql.DB, name, course, gender, d string) (int) {
	id, err := GetRaceByCourse(db, name, course, gender, d)
	if err == sql.ErrNoRows {
		distance := GetDistance(d)
		sqlStatement := `INSERT INTO races (name, course, distance, gender, inserted_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6);`
		_, err := db.Exec(sqlStatement, name, course, distance, gender, time.Now(), time.Now())
		check(err)

		id, err = GetRaceByCourse(db, name, course, gender, d)
	}
	return id
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

func ComputeAverage(db *sql.DB, id int) {
	race := GetRaceByID(db, id)
	insts := FindAllInstances(db, (*race).id)
	var results []Result
	for i := range *insts {
		r := GetInstanceResults(db, (*insts)[i].id)
		for _, each := range *r {
			results = append(results, each)
		}
	}

	trimmed_results := FilterDNFs(&results)
	mean, std_dev := CalculateStatistics(trimmed_results)

	// fmt.Printf("Race ID: %v\tMean: %v\t St. Dev: %v\n", (*races)[i].id, mean, std_dev)
	UpdateAverage(db, (*race).id, mean)
	UpdateStdDev(db, (*race).id, std_dev)
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

func GetDistance(d string) int {
	if d == "10000" {return 10000}
	if d == "8000" {return 8000}
	if d == "6000" {return 6000}
	if d == "5000" {return 5000}

	return -1
}

func GetRaceByCourse(db *sql.DB, name, course, gender, d string) (int, error) {
	var id int
	query := `SELECT id FROM races WHERE (name=$1 AND course=$2 AND gender=$3 AND distance=$4);`
	
	distance := GetDistance(d)
	if distance == -1 {
		fmt.Printf("Distance not recognized")
		os.Exit(1)
		// TODO - Error handling this
	}
	row := db.QueryRow(query, name, course, gender, distance)
	err := row.Scan(&id)
	return id, err
}

func GetRaceByID(db *sql.DB, id int) *Race {
	query := `SELECT id FROM races WHERE (id=$1);`
	row := db.QueryRow(query, id)
	var ret Race
	err := row.Scan(&ret.id, &ret.name, &ret.course, &ret.distance, &ret.gender, &ret.is_base, &ret.average, &ret.std_dev, &ret.correction_avg, &ret.correction_graph)
	check(err)
	return &ret
}

func UpdateAverage(db *sql.DB, id int, average float64) {
	update := `UPDATE races SET average=$2 WHERE id=$1;`
	_, err := db.Exec(update, id, average)
	check(err)
}

func UpdateRace(db *sql.DB, id int, correction float64) {
	update := `UPDATE races SET correction=$2 WHERE id=$1;`
	_, err := db.Exec(update, id, correction)
	check(err)
}

func UpdateStdDev(db *sql.DB, id int, std_dev float64) {
	update := `UPDATE races SET std_dev=$2 WHERE id=$1;`
	_, err := db.Exec(update, id, std_dev)
	check(err)
}