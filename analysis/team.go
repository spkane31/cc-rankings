package main
import (
	"database/sql"
	_ "fmt"

	_ "github.com/lib/pq"

)

func GetTeam(db *sql.DB, name string) int {
	checkStatement := `SELECT id FROM teams WHERE (name=$1);`
	row := db.QueryRow(checkStatement, name)
	var id int
	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		sqlStatement := `INSERT INTO teams (name) VALUES ($1);`
		_, err = db.Exec(sqlStatement, name)
		check(err)

		row = db.QueryRow(checkStatement, name)
		err = row.Scan(&id)
		return id
	}
	return id
}

func DeleteTeam(db *sql.DB, name string) {
	sqlStatement := `DELETE FROM teams WHERE (name=$1);`
	_, err := db.Exec(sqlStatement, name)
	check(err)
}