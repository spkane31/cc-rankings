package main

import (
	"database/sql"
	"fmt"
	"os"
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
	correction float64
	is_base bool
}

type Instance struct {
	id int
	date string
	race_id int 
	average float64
	std_dev float64
}

func UpdateRace(db *sql.DB, id int, correction float64) {
	update := `UPDATE races SET correction=$2 WHERE id=$1;`
	_, err := db.Exec(update, id, correction)
	check(err)
}

func FindRace(db *sql.DB, id int) int {
	var ret int
	query := `SELECT * FROM races WHERE id=$1;`
	row := db.QueryRow(query, id)
	err := row.Scan(&ret)
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
		err = rows.Scan(&result.id, &result.distance, &result.unit, &result.rating, &result.time, &result.race_instance_id, &result.runner_id)
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