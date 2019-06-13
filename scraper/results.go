package main

import (
	"database/sql"
  "fmt"

	_ "github.com/lib/pq"
)

type Result struct {
	id int
	race_instance_id int
	distance int
	rating float64
	time string
	runner_id int
}

func CreateResult(db *sql.DB, gender string, details []string, raceSummary map[string]string) int {
	// This will create a new runner given their name, and return the ID
	if debug {fmt.Println("\nCreating a new Result")}
	var id int
	runner_id, err := FindRunner(db, details[1], details[0], details[3])
	if err == sql.ErrNoRows {
		if debug {fmt.Println("No Runner found, creating a new one")}
		// Order is last name, first name
		runner_id = AddRunner(db, details[1], details[0], details[3])
	}

	if runner_id == 0 {
		fmt.Println("Could not find runner")
		fmt.Println(details)
	}
	if debug {fmt.Printf("Runner ID: %d\n", runner_id)}

	// NOw we have the runner_id, check for the team
	team_id, err := FindTeam(db, details[2])
	if err == sql.ErrNoRows {
		if debug {fmt.Println("No team found, creating a new one")}
		team_id = AddTeam(db, details[2])
	}
	if debug {fmt.Printf("Team ID: %d\n", team_id)}

	// Now lets connect the runner to the team
	ConnectRunnerTeam(db, runner_id, team_id)

	// Now lets check if this particular course/race is in our database
	race_id, err := FindCourse(db, raceSummary, gender)
	if err == sql.ErrNoRows {
		if debug {fmt.Println("This Race does not exist, creating a new RACE")}
		race_id = AddRace(db, raceSummary, gender)
	}
	if debug {fmt.Printf("Race ID: %d\n", race_id)}
	
	// Once the race exists, we want to create an instance of the race which
	// requires the date, and the race to link up to
	instance_id, err := FindInstance(db, raceSummary["date"], race_id)
	if err == sql.ErrNoRows {
		if debug {fmt.Println("New Race Instance")}
		instance_id = AddInstance(db, raceSummary["date"], race_id)
	}
	if debug {fmt.Printf("Instance ID: %d\n", instance_id)}

	// Now we create the result and link it to the runner and the race instance
	result_id, err := FindResult(db, details[4], raceSummary["mens_distance"], runner_id, instance_id)
	if err == sql.ErrNoRows {
		if debug {fmt.Println("Adding Result")}
		result_id = AddResult(db, details[4], raceSummary["mens_distance"], runner_id, instance_id)
	}
	if debug {fmt.Printf("Result ID: %d\n", result_id)}
	return id
}

func AddResult(db *sql.DB, time string, distance string, runner_id, instance_id int) int {
	var id int
	sqlStatement := `INSERT INTO results (distance, time, runner_id, unit, race_instance_id) VALUES ($1, $2, $3, $4, $5);`
	var d int
	var unit string
	if distance == "10K" {
		d = 10000
		unit = "METERS"
	} else if distance == "8K" {
		d = 8000
		unit = "METERS"
	} else if distance == "6K" {
		d = 6000
		unit = "METERS"
	}

	_, err := db.Exec(sqlStatement, d, time, runner_id, unit, instance_id)
	
	check(err)
	id, err = FindResult(db, time, distance, runner_id, instance_id)

	return id
}

func FindResult(db *sql.DB, time string, distance string, runner_id, instance_id int) (int, error) {
	var id int
	var d int
	// var unit string
	if distance == "10K" {
		d = 10000
		// unit = "METERS"
	} else if distance == "8K" {
		d = 8000
		// unit = "METERS"
	} else if distance == "6K" {
		d = 6000
		// unit = "METERS"
	}
	queryStatement := `SELECT id FROM results WHERE (time=$1 AND distance=$2 AND runner_id=$3 AND race_instance_id=$4);`
	row := db.QueryRow(queryStatement, time, d, runner_id, instance_id)
	err := row.Scan(&id)
	return id, err
}