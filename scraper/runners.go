package main

import (
	"database/sql"
  "fmt"

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

func CreateRunner(db *sql.DB, first, last string) int {
	// This will create a new runner given their name, and return the ID

	// First we should probably check for a runner
	checkStatement := `SELECT id FROM runners WHERE first_name=$1, last_name=$2;`
	row := db.QueryRow(checkStatement, first, last)
	var id int
	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		// If their is no hit on the query, then we create a new runner, requery, and return the id
		sqlStatement := `INSERT INTO runners (first_name, last_name) VALUES ($1, $2)`
	
		result, err := db.Exec(sqlStatement, first, last)
		check(err)
		fmt.Println(result)
		// return 
		row := db.QueryRow(checkStatement, first, last)
		err = row.Scan(&id)

		return id
	}
	fmt.Println(id)
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