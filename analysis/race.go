package main

import (
	"database/sql"
	"fmt"
	"os"
	"encoding/csv"
	// "bufio"
	// "io"

	_ "github.com/lib/pq"
)

var _, _, _ = fmt.Println, os.Exit, csv.NewReader

func UpdateRace(db *sql.DB, id int, correction float64) {
	update := `UPDATE races SET correction=$2 WHERE id=$1;`
	_, err := db.Exec(update, id, correction)
	check(err)
}

func FindRace(db *sql.DB, id int) int {
	var ret int
	query := `SELECT * FROM races WHERE id=$1;`
	row := db.QueryRow(query, id)
	err := row.Scan(&ret)
	check(err)
	return ret
}