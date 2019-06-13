package main 

import (
	"database/sql"
	// "fmt"
	// "strconv"

	_ "github.com/lib/pq"
)

type Race_Instance struct {
	ID int
	date string
	race_id int
}

func FindInstance(db *sql.DB, date string, race_id int) (int, error) {
	var id int
	queryStatement := `SELECT id FROM race_instances WHERE (date=$1 AND race_id=$2);`
	row := db.QueryRow(queryStatement, date, race_id)
	err := row.Scan(&id)
	return id, err
}

func AddInstance(db *sql.DB, date string, race_id int) int {
	var id int
	id, err := FindInstance(db, date, race_id)
	if err == sql.ErrNoRows {
		sqlStatement := `INSERT INTO race_instances (date, race_id) VALUES ($1, $2);`
		_, err = db.Exec(sqlStatement, date, race_id)

		id, err = FindInstance(db, date, race_id)
	}
	return id
}

