package rankings

import (
	"database/sql"
	"fmt"
	"time"
	"os"
	"strings"
	"strconv"
	// "math"

	_ "github.com/lib/pq"
)

var _, _, _ = fmt.Println, time.Now, os.Exit

type Result struct {
	id int
	distance int
	unit string
	rating float64
	time string
	race_instance_id int
	runner_id int
	scaled_time float64
	time_float float64
}

func CreateResult(db *sql.DB, details []string, distance, gender, course, date, race_name string, place int) (int, int, int, int) {
	// details = [last, first, year, school, time]
	debug := false

	team_id, err := FindTeam(db, details[3])
	if err == sql.ErrNoRows {
		if debug {fmt.Println("No team found, creating a new one")}
		team_id = AddTeam(db, details[3])
	} else {
		check(err)
	}
	if debug {fmt.Printf("Team ID: %d\n", team_id)}

	runner_id, err := GetRunnerID(db, details[1], details[0], details[2], gender, team_id)
	if err == sql.ErrNoRows {
		if debug {fmt.Println("No Runner found, creating a new one")}
		runner_id = AddRunner(db, details[1], details[0], details[2], gender, team_id)
	} else {
		check(err)
	}
	if debug {fmt.Printf("Runner ID: %d\n", runner_id)}

	ConnectRunnerTeam(db, runner_id, team_id)

	race_id, err := GetRaceByCourse(db, race_name, course, gender, distance)
	if err == sql.ErrNoRows {
		if debug {fmt.Println("This Race does not exist, creating a new RACE")}
		race_id = AddRace(db, race_name, course, gender, distance)
	} else {
		check(err)
	}
	if debug {fmt.Printf("Race ID: %d\n", race_id)}

	instance_id, err := GetInstance(db, date, race_id)
	if err == sql.ErrNoRows {
		if debug {fmt.Println("New Race Instance")}
		instance_id = AddInstance(db, date, race_id, gender, distance)
	} else {
		check(err)
	}
	if debug {fmt.Printf("Instance ID: %d\n", instance_id)}

	// Now we create the result and link it to the runner and the race instance
	result_id, err := FindResult(db, details[4], distance, runner_id, instance_id)
	if err == sql.ErrNoRows {
		if debug {fmt.Println("Adding Result")}
		result_id = AddResult(db, details[4], distance, runner_id, instance_id, gender, place, date)
	} else {
		check(err)
	}
	if debug {fmt.Printf("Result ID: %d\n", result_id)}
	
	return runner_id, result_id, race_id, instance_id
}

func FindResult(db *sql.DB, time string, distance string, runner_id, instance_id int) (int, error) {
	var id int
	d := GetDistance(distance)
	queryStatement := `SELECT id FROM results WHERE (time=$1 AND distance=$2 AND runner_id=$3 AND race_instance_id=$4);`
	row := db.QueryRow(queryStatement, time, d, runner_id, instance_id)
	err := row.Scan(&id)
	return id, err
}

func FindResultByID(db *sql.DB, id int) Result {
	query := `SELECT * FROM results WHERE (id=$1);`

	row := db.QueryRow(query, id)
	var ret Result
	err := row.Scan(&ret.id, &ret.distance, &ret.unit, &ret.rating, &ret.time, &ret.race_instance_id, &ret.runner_id, &ret.scaled_time, &ret.time_float)
	check(err)
	return ret
}

func AddResult(db *sql.DB, t string, distance string, runner_id, instance_id int, gender string, place int, date string) int {
	var id int
	var scaled float64
	time_float := GetTime(t)
	if distance == "10000" && gender == "MALE" {
		scaled = time_float / 1.268
	} else if distance == "8000" && gender == "MALE" {
		scaled = time_float
	} else if distance == "5000" && gender == "FEMALE" {
		scaled = time_float
	} else if distance == "6000" && gender == "FEMALE" {
		scaled = time_float / 1.213
	} else {
		scaled = 0.0
	}
	sqlStatement := `INSERT INTO results 
										(distance, time, runner_id, unit, race_instance_id, scaled_time, time_float, date, gender, place, inserted_at, updated_at) 
										VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);`
	var unit string
	d := GetDistance(distance)
	
	_, err := db.Exec(sqlStatement, d, t, runner_id, unit, instance_id, scaled, time_float, date, gender, place, time.Now(), time.Now())
	check(err)

	id, err = FindResult(db, t, distance, runner_id, instance_id)

	return id
}

func GetTime(time string) float64 {
	if time == "DNF" || time == "DNS" {return 0}

	var ret float64
	t := strings.Split(time, ":")
	mult := 1.0
	for i := len(t)-1; i >= 0; mult *= 60.0 {
		f, _ := strconv.ParseFloat(strings.Replace((t[i]), " ", "", -1), 16)
		ret += f * mult
		i--
	}
	return ret
}

func FindResultsForRunner(db * sql.DB, id int) *[]int {
	var result_ids []int

	query := `SELECT id FROM results WHERE runner_id=$1 ORDER BY race_instance_id;`
	rows, err := db.Query(query, id)
	check(err)
	defer rows.Close()

	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		result_ids = append(result_ids, id)
	}
	
	return &result_ids
}

func GetRaceResults(db *sql.DB, id int) *[]Result {
	var results []Result

	query := `SELECT * FROM results WHERE race_id=$1;`
	rows, err := db.Query(query, id)
	check(err)
	defer rows.Close()

	for rows.Next() {
		var result Result
		err = rows.Scan(&result.id, &result.distance, &result.unit, &result.rating, &result.time, &result.race_instance_id, &result.runner_id, &result.scaled_time, &result.time_float)
		results = append(results, result)
	}

	return &results
}

func FilterDNFs(results *[]Result) *[]Result {
	var ret []Result

	for _, each := range *results {
		if GetTime(each.time) != 0 {
			ret = append(ret, each)
		}
	}

	return &ret

}

func ScaleTime(time float64, distance int) float64 {
	if distance == 6000 {
		return time / 1.213
	} else if distance == 10000 {
		return time / 1.268
	} else if distance == 5000 || distance == 8000 {
		return time
	} else {
		return -1
	}
}

func CheckResultsYears(db *sql.DB, a, b int) bool {

	query := `SELECT date FROM results WHERE id=$1;`

	row := db.QueryRow(query, a)
	date_a := time.Now()
	err := row.Scan(&date_a)
	check(err)
	
	row = db.QueryRow(query, b)
	date_b := time.Now()
	err = row.Scan(&date_b)
	check(err)

	return date_b.Year() == date_a.Year()
}

func GetRaceIDFromResult(db *sql.DB, id int) int {
	query := `SELECT race_instance_id FROM results WHERE id=$1;`
	var inst_id int
	row := db.QueryRow(query, id)
	err := row.Scan(&inst_id)
	check(err)

	query = `SELECT race_id FROM race_instances WHERE id=$1;`
	var race_id int
	row = db.QueryRow(query, inst_id)
	err = row.Scan(&race_id)
	check(err)

	return race_id
}

func GetEdgeInformation(db *sql.DB, result_a, result_b int) (int, int, float64) {

	result_query := `SELECT race_instance_id, time, distance FROM results WHERE id=$1;`
	var inst_id_a int
	var time_a string
	var dist_a int

	row := db.QueryRow(result_query, result_a)
	err := row.Scan(&inst_id_a, &time_a, &dist_a)
	check(err)
	scaled_a := ScaleTime(GetTime(time_a), dist_a)

	race_query := `SELECT race_id FROM race_instances WHERE id=$1;`
	var race_id_a int
	row = db.QueryRow(race_query, inst_id_a)
	err = row.Scan(&race_id_a)
	check(err)

	var time_b string
	var inst_id_b int
	var race_id_b int
	var dist_b int

	row = db.QueryRow(result_query, result_b)
	err = row.Scan(&inst_id_b, &time_b, &dist_b)
	check(err)
	scaled_b := ScaleTime(GetTime(time_b), dist_b)

	row = db.QueryRow(race_query, inst_id_b)
	err = row.Scan(&race_id_b)

	return race_id_a, race_id_b, scaled_b - scaled_a
}

func UpdateRatings(db *sql.DB) {
	query := `SELECT id, correction_graph FROM races WHERE correction_graph != 0;`
	rows, err := db.Query(query)
	check(err)
	defer rows.Close()

	// TODO - Do I want to keep race_instances or move to something else?
	resultsQuery := `SELECT id, scaled_time, gender FROM results WHERE race_instance_id =$1;`
	updateQuery := `UPDATE results SET rating=$2 WHERE id=$1;`
	for rows.Next() {
		// Get information for the race
		var id int
		var correction_graph float64
		err = rows.Scan(&id, &correction_graph)
		check(err)

		// Get information for all of the result
		results, err := db.Query(resultsQuery, id) 
		check(err)
		defer results.Close()

		for results.Next() {
			var result_id int
			var rating float64
			var scaled_time float64
			var gender string

			err = results.Scan(&result_id, &scaled_time, &gender)

			if scaled_time != 0 {
				check(err)

				if gender == "MALE" {
					rating = (1900 - scaled_time - correction_graph) / (8000.0 / 1609.0)
				} else if gender == "FEMALE" {
					rating = (1350 - scaled_time - correction_graph) / (5000.0 / 1609.0)
				}

				_, err := db.Exec(updateQuery, result_id, rating)
				check(err)
			} else {
				_, err := db.Exec(updateQuery, result_id, rating)
				check(err)
			}

		}
	}
}

func ResetRatings(db *sql.DB) {
	update := `UPDATE results SET rating=0;`
	_, err := db.Exec(update)
	check(err)
}