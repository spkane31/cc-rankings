package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	
)

func CreateDatabase() *sql.DB{
	db := ConnectToPSQL()

	dropStatement := `DROP DATABASE rankings_test`
	_, err := db.Exec(dropStatement)
	check(err)

	createStatement := `
		CREATE DATABASE rankings_test;

		CREATE TABLE teams (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			region VARCHAR(255),
			conference VARCHAR(255)
		)

		CREATE TABLE runners (
			id SERIAL PRIMARY KEY,
			first_name VARCHAR(255) NOT NULL,
			last_name VARCHAR(255) NOT NULL,
			year VARCHAR(255),
			team_id INT,
			FOREIGN KEY (team_id) REFERENCES teams(id)
		);
	`

	_, err = db.Exec(createStatement)
	check(err)

	err = db.Ping()
	check(err)
	fmt.Println("Created Database!")

	return db
}