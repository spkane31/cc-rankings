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
	dbname = "rankings_dev"
)

func CreateConnectionString() string {
	str := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	 host, port, user, password, dbname)
	return str
}

func ConnectToPSQL() *sql.DB {
	db, err := sql.Open("postgres", CreateConnectionString())
	check(err)

	// Let's ping the database to check that the connection worked
	err = db.Ping()
	check(err)
	fmt.Println("Successfully connected!")
	return db
}

func InsertRunner(db *sql.DB, first, last, team, year string) {
	sqlStatement := `
	INSERT INTO runners (first_name, last_name, year)
	VALUES ($1, $2, $3)`

	result, err := db.Exec(sqlStatement, first, last, year)
	check(err)
	fmt.Println(result)
}

func UpdateRunner(db *sql.DB, first, last, team, year string) {
} 

func QueryRunnerFromID(db *sql.DB, id int) Runner {
	var runner Runner
	sqlStatement := `SELECT id, first_name, last_name, year FROM runners WHERE id=$1;`
	row := db.QueryRow(sqlStatement, id)
	err := row.Scan(&runner.ID, &runner.FirstName, &runner.LastName, &runner.Year)
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

func GetRunners(db *sql.DB, n int) *[]Runner {
	var runners []Runner
	sqlStatement := `SELECT * FROM runners LIMIT $1`

	rows, err := db.Query(sqlStatement, n)
	check(err)

	defer rows.Close()
	for rows.Next() {
		var r Runner 	
		err := rows.Scan(&r.ID, &r.FirstName, &r.LastName, &r.Team, &r.Year)
		check(err)
		runners = append(runners, r)
	}

	// Check for errors encountered during iteration
	err = rows.Err()
	check(err)

	return &runners
}