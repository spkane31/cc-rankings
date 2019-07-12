package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"strconv"
	"math"
	// "time"
	_ "github.com/lib/pq"
)

var _, _, _ = fmt.Println, os.Exit, math.Round


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

func FindResult(db *sql.DB, id int) Result {
	queryStatement := `SELECT * FROM results WHERE (id=$1);`

	row := db.QueryRow(queryStatement, id)
	var r Result
	err := row.Scan(&r.id, &r.distance, &r.unit, &r.rating, &r.time, &r.race_instance_id, &r.runner_id, &r.scaled_time, &r.time_float)
	check(err)
	return r
}

func Analyze(db *sql.DB, runners []Runner) {
	// fmt.Printf("%v: Analyzing Results\n", time.Now().Format("01-02-2006, 15:04:05"))
	graph := make(map[Pair]*Edge)
	for _, runner := range runners {
		for i := 0; i < len(runner.results)-1; i++ {
			for j := i+1; j < len(runner.results); j++ {

				if runner.results[i].distance == runner.results[j].distance {
					p := Pair{runner.results[i].race_instance_id, runner.results[j].race_instance_id}
					e, has := graph[p]
					if has == false {
						var e *Edge = new(Edge)
						fmt.Println("NEW EDGE")
						fmt.Printf("%v -> %v\n", runner.results[i].race_instance_id, runner.results[j].race_instance_id)
						e.to = runner.results[i].race_instance_id
						(*e).from = runner.results[j].race_instance_id
						e.count = 1
						e.total_time = 	GetTime(runner.results[i].time) - GetTime(runner.results[j].time)
						graph[p] = e
					} else {
						dif := GetTime(runner.results[i].time) - GetTime(runner.results[j].time)
						e.count += 1
						e.total_time += dif
					}
				}

			}
		}
	}
	for p, e := range graph {
		PrintEdgeNames(db, p.to, p.from)
		w := fmt.Sprintf("%.2f", e.total_time / float64(e.count))
		fmt.Printf("Total Connections: %v\tTotal Time Difference: %v\tWeight: %v\n", e.count, e.total_time, w)
		UpdateRace(db, p.to, e.total_time / float64(e.count))
	}
}

func GetTime(time string) float64 {
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

func FindResultsForRunner(db *sql.DB, id int) *[]Result {
	var res []Result

	queryStatement := `SELECT * FROM results WHERE runner_id=$1 ORDER BY race_instance_id;`

	rows, err := db.Query(queryStatement, id)
	check(err)
	defer rows.Close()

	for rows.Next() {
		var result Result
		err = rows.Scan(&result.id, &result.distance, &result.unit, &result.rating, &result.time, &result.race_instance_id, &result.runner_id, &result.scaled_time, &result.time_float)
		res = append(res, result)
	}

	return &res
}

func GetRaceResults(db *sql.DB, id int) *[]Result {
	var ret []Result

	queryStatement := `SELECT * FROM results WHERE race_id=$1;`
	rows, err := db.Query(queryStatement, id)
	check(err)
	defer rows.Close()

	for rows.Next() {
		var result Result
		err = rows.Scan(&result.id, &result.distance, &result.unit, &result.rating, &result.time, &result.race_instance_id, &result.runner_id, &result.scaled_time, &result.time_float)
		ret = append(ret, result)
	}

	return &ret
}

func FilterDNFs(results *[]Result) []Result {
	var ret []Result

	for _, each := range *results {
		if GetTime(each.time) != 0 {
			ret = append(ret, each)
		}
	}

	return ret
}