package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host = "localhost"
	port = 5432
	user = "postgres"
	password = "postgres"
	dbname = "rankings_test"
)

type Runner struct {
	ID 				int
	FirstName string
	LastName 	string
	Team			string
	Year			string
}

func CreateConnectionString() string {
	str := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	 host, port, user, password, dbname)
	return str
}

func ConnectToPSQL() *sql.DB {
	db, err := sql.Open("postgres", CreateConnectionString())
	check(err)
	// Probably want the below to be in a different function
	// defer db.Close()

	// Let's ping the database to check that the connection worked
	err = db.Ping()
	check(err)
	fmt.Println("Successfully connected!")
	return db
}

func InsertRunner(db *sql.DB, first, last, team, year string) {
	sqlStatement := `
	INSERT INTO runner (first_name, last_name, team, year)
	VALUES ($1, $2, $3, $4)`

	result, err := db.Exec(sqlStatement, first, last, team, year)
	check(err)
	fmt.Println(result)
}

func UpdateRunner(db *sql.DB, first, last, team, year string) {
} 

func QueryRunnerFromID(db *sql.DB, id int) Runner {
	var runner Runner
	sqlStatement := `SELECT * FROM runner WHERE id=$1;`
	row := db.QueryRow(sqlStatement, id)
	err := row.Scan(&runner.ID, &runner.FirstName, &runner.LastName, &runner.Team, &runner.Year)
	switch err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned")
		case nil:
			fmt.Println(runner)
		default:
			panic(err)
	}

	return runner
}