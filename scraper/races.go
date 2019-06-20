package main

import (
	"database/sql"
	"fmt"
	"os"
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

func GetDistance(d string) int {
	if d == "10K" {return 10000}
	if d == "8K" {return 8000}
	if d == "6K" {return 6000}
	if d == "5K" {return 5000}

	return -1
}

func FindCourse(db *sql.DB, name, course, gender, d string) (int, error) {
	var id int
	queryStatement := `SELECT id FROM races WHERE (name=$1 AND course=$2 AND gender=$3 AND distance=$4);`
	// d := m["mens_distance"]
	distance := GetDistance(d)
	if distance == -1 {
		fmt.Printf("Distance not recognized")
		os.Exit(1)
	}
	row := db.QueryRow(queryStatement, name, course, gender, distance)
	err := row.Scan(&id)
	return id, err
}

func AddRace(db *sql.DB, name, course, gender, d string) (int) {
	id, err := FindCourse(db, name, course, gender, d)
	distance := GetDistance(d)
	if err == sql.ErrNoRows {
		sqlStatement := `INSERT INTO races (name, course, distance, gender) VALUES ($1, $2, $3, $4);`
		_, err := db.Exec(sqlStatement, name, course, distance, gender)
		check(err)

		id, err = FindCourse(db, name, course, gender, d)
		return id
	}
	return id
}