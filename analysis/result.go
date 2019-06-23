package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"strconv"
	"math"
	_ "github.com/lib/pq"
)

var _, _, _ = fmt.Println, os.Exit, math.Round

type Pair struct {
	to, from int
}

type Edge struct {
	to, from int
	count int
	total_time float64
}

type Result struct {
	id int
	distance int
	unit string
	rating float64
	time string
	race_instance_id int
	runner_id int
}

func FindResult(db *sql.DB, id int) Result {
	queryStatement := `SELECT * FROM results WHERE (id=$1);`

	row := db.QueryRow(queryStatement, id)
	var r Result
	err := row.Scan(&r.id, &r.distance, &r.unit, &r.rating, &r.time, &r.race_instance_id, &r.runner_id)
	check(err)
	return r
}

func Analyze(db *sql.DB, runners []Runner) {
	fmt.Println("Analyzing Results")
	graph := make(map[Pair]*Edge)
	for _, runner := range runners {
		for i := 0; i < len(runner.results)-1; i++ {
			for j := i+1; j < len(runner.results); j++ {

				if runner.results[i].distance == runner.results[j].distance {
					// var e *Edge
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
						// This connection already exists: Need to modify it
						dif := GetTime(runner.results[i].time) - GetTime(runner.results[j].time)
						e.count += 1
						e.total_time += dif
					}
					// graph[p] = e
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
	for i := len(t)-1; i >= 0; mult *= 60 {
		f, _ := strconv.ParseFloat(t[i], 16)
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
		err = rows.Scan(&result.id, &result.distance, &result.unit, &result.rating, &result.time, &result.race_instance_id, &result.runner_id)
		res = append(res, result)
	}

	return &res
}

func PrintEdgeNames(db *sql.DB, i, j int) {
	instanceStatement := `SELECT race_id FROM race_instances WHERE id=$1;`
	row := db.QueryRow(instanceStatement, i)
	var race_id int
	var to string
	var from string
	err := row.Scan(&race_id)
	check(err)
	raceStatement := `SELECt name FROM races WHERE id=$1;`
	row = db.QueryRow(raceStatement, race_id)
	err = row.Scan(&to)
	check(err)

	row = db.QueryRow(instanceStatement, j)
	err = row.Scan(&race_id)
	check(err)
	row = db.QueryRow(raceStatement, race_id)
	err = row.Scan(&from)
	check(err)

	fmt.Printf("\n%v -> %v\n", to[0:len(to)-1], from[0:len(from)-1])
}