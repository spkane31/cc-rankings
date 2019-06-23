package main

import (
	"database/sql"
	"fmt"
	"os"
	"encoding/csv"
	"bufio"
	"io"

	_ "github.com/lib/pq"
)

var _, _, _ = fmt.Println, os.Exit, csv.NewReader

type Runner struct {
	id int
	first_name string
	last_name string
	year string
	gender string
	team_id string
	results []Result
}

func FindConnections(db *sql.DB) {
	conns := 0
	csvFile, err := os.Open("../scraper/RaceResults/NCAADICrossCountryChampionship/mens.csv")
	check(err)
	query := `SELECT id FROM runners WHERE (first_name=$1 AND last_name=$2 AND gender='MALE');`
	
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var id int
	for {
		line, err := reader.Read() 
		if err == io.EOF {
			break
		} else {
			row := db.QueryRow(query, line[1], line[0])
			err := row.Scan(&id)
			check(err)

			if CheckMultipleConnections(db, id) {
				conns++
			} else {
				fmt.Println(line)
			}

		}
	}

	fmt.Println(conns)

	// query := `SELECT COUNT(*) FROM runners;`
	rows, err := db.Query(query)
	var count int
	check(err)
	for rows.Next() {
		err = rows.Scan(&count)
		check(err)
	}
	fmt.Println(count)
	for i := 1; i <= count; i++ {
		if CheckMultipleConnections(db, i) {
			conns++
		}
	}
	fmt.Println(conns)
	os.Exit(1)

}

func CheckMultipleConnections(db *sql.DB, id int) bool {
	var results []Result
	query := `SELECT * FROM results WHERE runner_id=$1;`
	rows, err := db.Query(query, id)
	check(err)
	for rows.Next() {
		var result Result
		err = rows.Scan(&result.id, &result.distance, &result.unit, &result.rating, &result.time, &result.race_instance_id, &result.runner_id)
		results = append(results, result)
	}

	return len(results) > 1
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

	mens, womens := BuildRunnerConnections(db, &allResults)
	// os.Exit(1)
	Analyze(db, *mens)
	Analyze(db, *womens)
}

func BuildRunnerConnections(db *sql.DB, results *[]Result) (*[]Runner, *[]Runner) {
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

	// This can be made quicker, but now I want to organize into two lists, one for men and one for women.
	// TODO - Make this more efficient
	var mens_runners []Runner
	var womens_runners []Runner

	for _, runner := range runners {
		if runner.gender == "MENS" {
			mens_runners = append(mens_runners, runner)
		} else {
			womens_runners = append(womens_runners, runner)
		}
	}

	// for i := range runners {
	// 	if len(runners[i].results) < 2 {
	// 		os.Exit(1)
	// 	}
	// }
	fmt.Println(len(mens_runners), len(womens_runners))
	return &mens_runners, &womens_runners
}


func AddRunner(first, last, team string, db *sql.DB) int {
	checkStatement := `SELECT id FROM runners WHERE (first_name=$1 AND last_name=$2 AND team_id=$3);`
	
	team_id := GetTeam(db, team)
	row := db.QueryRow(checkStatement, first, last, team_id)
	var id int
	err := row.Scan(&id)
	
	if err == sql.ErrNoRows {
		sqlStatement := `INSERT INTO runners (first_name, last_name) VALUES ($1, $2);`

		_, err := db.Exec(sqlStatement, first, last)
		check(err)
		
		row = db.QueryRow(checkStatement, first, last, team_id)
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
	err := row.Scan(&ret.id, &ret.first_name, &ret.last_name, &ret.year, &ret.team_id, &ret.gender)
	return &ret, err
}