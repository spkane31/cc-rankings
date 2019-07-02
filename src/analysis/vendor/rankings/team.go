package rankings

import (
	"database/sql"
	// "fmt"
	"time"
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
	query := `SELECT id FROM teams WHERE (name=$1);`
	row := db.QueryRow(query, name)
	var id int
	err := row.Scan(&id)
	return id, err
}

func AddTeam(db *sql.DB, name string) int {
	var id int
	id, err := FindTeam(db, name)
	if err == sql.ErrNoRows {
		insert := `INSERT INTO teams (name, inserted_at, updated_at) VALUES ($1, $2, $3);`
		_, err = db.Exec(insert, name, time.Now(), time.Now())
		check(err)

		id, err = FindTeam(db, name)
	}
	return id
}

func DeleteTeam(db *sql.DB, name string) {
	sqlStatement := `DELETE FROM teams WHERE (name=$1);`
	_, err := db.Exec(sqlStatement, name)
	check(err)
}