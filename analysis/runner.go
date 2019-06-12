package main

import (
	"database/sql"
	_ "fmt"

	_ "github.com/lib/pq"
)

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