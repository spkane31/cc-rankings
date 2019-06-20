package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// This is the details for connecting to the Elixir dev database
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
	// Probably want the below to be in a different function
	// defer db.Close()

	// Let's ping the database to check that the connection worked
	err = db.Ping()
	check(err)
	fmt.Println("Successfully connected!")
	return db
}