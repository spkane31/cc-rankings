package main

import (
	"database/sql"
	// "fmt"
	// "strconv"

	_ "github.com/lib/pq"
)

type Race struct {
	ID int
	name string
	course string
	distance int
	gender string
	correction float64
}

func FindCourse(db *sql.DB, m map[string]string, g string) (int, error) {
	var id int
	queryStatement := `SELECT id FROM races WHERE (name=$1 AND course=$2 AND gender=$3 AND distance=$4);`
	d := m["mens_distance"]
	var distance int
	if d == "10K" {
		distance = 10000
	}
	row := db.QueryRow(queryStatement, m["name"], m["course"], g, distance)
	err := row.Scan(&id)
	return id, err
}

func AddRace(db *sql.DB, m map[string]string, g string) (int) {
	id, err := FindCourse(db, m, g)
	d := m["mens_distance"]
	var distance int
	if d == "10K" {
		distance = 10000
	}
	if err == sql.ErrNoRows {
		sqlStatement := `INSERT INTO races (name, course, distance, gender) VALUES ($1, $2, $3, $4);`
		_, err := db.Exec(sqlStatement, m["name"], m["course"], distance, g)
		check(err)

		id, err = FindCourse(db, m, g)

		return id
	}

	return id
}