package main

import (
	"database/sql"
	"fmt"
	"strings"
	"strconv"
	"os"

	_ "github.com/lib/pq"
)

var _ = os.Exit

type Result struct {
	id int
	race_instance_id int
	distance int
	rating float64
	time string
	runner_id int
	scaled_time float64
	time_float float64
}

func CreateResult(db *sql.DB, details []string, distance, gender, course, date, race_name string) int {
	// detailts = [last, first, year, school, time]

	// First create the team w we have the runner_id, check for the team
	debug = false
	// fmt.Println(details)
	team_id, err := FindTeam(db, details[3])
	if err == sql.ErrNoRows {
		if debug {fmt.Println("No team found, creating a new one")}
		team_id = AddTeam(db, details[3])
	}
	if debug {fmt.Printf("Team ID: %d\n", team_id)}
	
	// This will create a new runner given their name, team, and year, and return the ID
	if debug {fmt.Println("\nCreating a new Result")}
	runner_id, err := FindRunner(db, details[1], details[0], details[2], gender, team_id)
	if err == sql.ErrNoRows {
		if debug {fmt.Println("No Runner found, creating a new one")}
		// Order is last name, first name
		runner_id = AddRunner(db, details[1], details[0], details[2], gender, team_id)
	}

	if debug {fmt.Printf("Runner ID: %d\n", runner_id)}

	// Now lets connect the runner to the team
	ConnectRunnerTeam(db, runner_id, team_id)

	// Now lets check if this particular course/race is in our database
	race_id, err := FindCourse(db, race_name, course, gender, distance)
	if err == sql.ErrNoRows {
		if debug {fmt.Println("This Race does not exist, creating a new RACE")}
		race_id = AddRace(db, race_name, course, gender, distance)
	}
	if debug {fmt.Printf("Race ID: %d\n", race_id)}
	
	// Once the race exists, we want to create an instance of the race which
	// requires the date, and the race to link up to
	instance_id, err := FindInstance(db, date, race_id)
	if err == sql.ErrNoRows {
		if debug {fmt.Println("New Race Instance")}
		instance_id = AddInstance(db, date, race_id)
	}
	if debug {fmt.Printf("Instance ID: %d\n", instance_id)}

	// Now we create the result and link it to the runner and the race instance
	result_id, err := FindResult(db, details[4], distance, runner_id, instance_id)
	if err == sql.ErrNoRows {
		if debug {fmt.Println("Adding Result")}
		result_id = AddResult(db, details[4], distance, runner_id, instance_id, gender)
	}
	if debug {fmt.Printf("Result ID: %d\n", result_id)}
	
	return result_id
}

func AddResult(db *sql.DB, time string, distance string, runner_id, instance_id int, gender string) int {
	var id int
	var scaled float64
	time_float := GetTime(time)
	if distance == "10000" {
		scaled = time_float / 1.268
	} else if distance == "8000" {
		scaled = time_float
	} else if distance == "5000" && gender == "WOMENS" {
		scaled = time_float
	} else if distance == "6000" && gender == "WOMENS" {
		scaled = time_float / 1.213
	} else {
		scaled = 0.0
	}
	sqlStatement := `INSERT INTO results (distance, time, runner_id, unit, race_instance_id, scaled_time, time_float) VALUES ($1, $2, $3, $4, $5, $6, $7);`
	var unit string
	d := GetDistance(distance)
	_, err := db.Exec(sqlStatement, d, time, runner_id, unit, instance_id, scaled, time_float)
	
	check(err)
	id, err = FindResult(db, time, distance, runner_id, instance_id)

	return id
}

func FindResult(db *sql.DB, time string, distance string, runner_id, instance_id int) (int, error) {
	var id int
	d := GetDistance(distance)
	queryStatement := `SELECT id FROM results WHERE (time=$1 AND distance=$2 AND runner_id=$3 AND race_instance_id=$4);`
	row := db.QueryRow(queryStatement, time, d, runner_id, instance_id)
	err := row.Scan(&id)
	return id, err
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