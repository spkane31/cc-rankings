package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var _, _ = fmt.Println, os.Exit

type Runner struct {
	id int
	first_name string
	last_name string
	year string
	team_id string
	results []Result
}

type Race struct {
	id int
	name string
	course string
	distance int
	gender string
	correction float64
}

type Instance struct {
	id int
	date string
	race_id int
}

func FindAllConnections(db *sql.DB) {
	queryStatement := `SELECT runner_id FROM results GROUP BY runner_id HAVING COUNT(runner_id) > 1;`
	rows, err := db.Query(queryStatement)
	var allResults []Result
	check(err)
	var id int
	var count int
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&id)
		count++
		results := FindResultsForRunner(db, id)
		for _, result := range (*results) {
			allResults = append(allResults, result)
		}	
	}
	fmt.Printf("%v Connections!\n", count)

	connections := BuildRunnerConnections(db, &allResults)
	// os.Exit(1)
	Analyze(db, *connections)
}

func BuildRunnerConnections(db *sql.DB, results *[]Result) *[]Runner {
	// We want to build an array Runners that has all the results for each runner contained
	fmt.Println("Building Runner Connections!")
	var runners []Runner

	// Find the runner for the first result
	result := (*results)[0]
	newRunner, err := FindRunner(db, result.runner_id)
	check(err)

	newRunner.results = append(newRunner.results, result)
	
	runners = append(runners, *newRunner)

	var inResults bool
	for i := 1; i < len(*results); i++ {
		// var newRunner Runner
		result := (*results)[i]
		
		inResults = false
		// Check to see if the runner is already in our list
		for i, runner := range runners {
			if runner.id == result.runner_id {
				newRes := append(runner.results, result)
				runner.results = newRes
				runners[i] = runner
				inResults = true
				break
			}
		}

		if !inResults {
			// New Runner, find the runner, add the result and add to the runners list
			newRunner, err = FindRunner(db, result.runner_id) 
			check(err)
			newRunner.results = append(newRunner.results, result)
			runners = append(runners, *newRunner)
		}
	}

	// for i := range runners {
	// 	if len(runners[i].results) < 2 {
	// 		os.Exit(1)
	// 	}
	// }
	return &runners
}


func AddRunner(first, last string, db *sql.DB) int {
	checkStatement := `SELECT id FROM runners WHERE (first_name=$1 AND last_name=$2);`
	row := db.QueryRow(checkStatement, first, last)
	var id int
	err := row.Scan(&id)
	
	if err == sql.ErrNoRows {
		sqlStatement := `INSERT INTO runners (first_name, last_name) VALUES ($1, $2);`

		_, err := db.Exec(sqlStatement, first, last)
		check(err)
		
		row = db.QueryRow(checkStatement, first, last)
		err = row.Scan(&id)
		return id
	}
	return id
}

func DeleteRunner(db *sql.DB, first, last string) {
	sqlStatement := `DELETE FROM runners WHERE (first_name=$1 AND last_name=$2);`
	_, err := db.Exec(sqlStatement, first, last)
	check(err)
}

func FindRunner(db *sql.DB, id int) (*Runner, error) {
	var ret Runner
	queryStatement := `SELECT * FROM runners WHERE (id=$1);`
	row := db.QueryRow(queryStatement, id)
	err := row.Scan(&ret.id, &ret.first_name, &ret.last_name, &ret.year, &ret.team_id)
	return &ret, err
}