package main

import (
	"database/sql"
	// "fmt"

	_ "github.com/lib/pq"
)

type Runner struct {
	ID 				int
	FirstName string
	LastName 	string
	Team			string
	Team_ID 	int
	Year			string
}

func FindRunner(db *sql.DB, first, last, year string) (int, error) {
	var id int
	// TODO - Add a team check here too!
	queryStatement := `SELECT id FROM runners WHERE (first_name=$1 AND last_name=$2 AND year=$3);`
	row := db.QueryRow(queryStatement, first, last, year)
	err := row.Scan(&id)
	return id, err
}

func AddRunner(db *sql.DB, first, last, year string) int {
	// This will create a new runner given their name, and return the ID

	// First we should probably check for a runner
	checkStatement := `SELECT id FROM runners WHERE (first_name=$1 AND last_name=$2 AND year=$3);`
	// row := db.QueryRow(checkStatement, first, last, year)
	// var id int
	// err := row.Scan(&id)
	id, err := FindRunner(db, first, last, year)
	if err == sql.ErrNoRows {
		// If their is no hit on the query, then we create a new runner, requery, and return the id
		sqlStatement := `INSERT INTO runners (first_name, last_name, year) VALUES ($1, $2, $3)`
	
		_, err := db.Exec(sqlStatement, first, last, year)
		check(err)
		
		row := db.QueryRow(checkStatement, first, last, year)
		err = row.Scan(&id)

		return id
	}
	return id
}

func AddYearToRunner(db *sql.DB, id int, year string) int {
	// TODO - add errors to this in case things fail

	updateStatement := `
		UPDATE users
		SET year=$1 
		WHERE id=$2;`

	res, err := db.Exec(updateStatement, year, id)
	check(err)
	count, err := res.RowsAffected()
	check(err)
	if count == 0 {
		return -1
	}
	return id
}

func ConnectRunnerTeam(db *sql.DB, runner, team int) int { 
	updateStatement := `UPDATE runners SET team_id=$1 WHERE id=$2`
	res, err := db.Exec(updateStatement, team, runner)
	check(err)
	count, err := res.RowsAffected()
	if count == 0 {
		return -1
	} else {
		return 0
	}
}