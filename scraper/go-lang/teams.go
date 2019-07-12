package main

import (
	"database/sql"
	// "fmt"

	_ "github.com/lib/pq"
)

type Team struct {
	ID int
	Name string	
	Region string
	Conference string
	Runners []Runner
}

func FindTeam(db *sql.DB, name string) (int, error) {
	var id int
	queryStatement := `SELECT id FROM teams WHERE (name=$1);`
	row := db.QueryRow(queryStatement, name)
	err := row.Scan(&id)
	return id, err
}

func AddTeam(db *sql.DB, name string) int {
	var id int
	queryStatement := `SELECT id FROM teams WHERE (name=$1);`
	row := db.QueryRow(queryStatement, name)
	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		sqlStatement := `INSERT INTO teams (name) VALUES ($1);`
		_, err = db.Exec(sqlStatement, name)
		check(err)

		row := db.QueryRow(queryStatement, name)
		err = row.Scan(&id)
	}
	return id
}