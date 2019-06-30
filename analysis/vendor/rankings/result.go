package rankings

import (
	"database/sql"
	// "fmt"
	// "os"
	"strings"
	"strconv"
	// "math"
	// "time"
	_ "github.com/lib/pq"
)

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

func FindResultByID(db *sql.DB, id int) Result {
	query := `SELECT * FROM results WHERE (id=$1);`

	row := db.QueryRow(query, id)
	var ret Result
	err := row.Scan(&ret.id, &ret.distance, &ret.unit, &ret.rating, &ret.time, &ret.race_instance_id, &ret.runner_id, &ret.scaled_time, &ret.time_float)
	check(err)
	return ret
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

func FindResultsForRunner(db * sql.DB, id int) *[]Result {
	var results []Result

	query := `SELECT * FROM results WHERE runner_id=$1 ORDER BY race_instance_id;`

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