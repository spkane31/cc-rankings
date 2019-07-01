package rankings

import (
	"fmt"
	"database/sql"

	_ "github.com/lib/pq"
)

var _ = fmt.Printf

type Runner struct {
	ID int
	FirstName string
	LastName string
	Team_ID int
	Year string
	results []Result
}

func GetRunnerID(db *sql.DB, first, last, year, gender string, team_id int) (int, error) {
	var id int

	query := `SELECT id FROM runners WHERE (first_name=$1 AND last_name=$2, AND year=$3, AND team_id=$4, AND gender=$5);`
	row := db.QueryRow(query, first, last, year, team_id, gender)
	err := row.Scan(&id)
	return id, err
}

func AddRunner(db *sql.DB, first, last, year, gender string, team_id int) int {
	var id int
	
	id, err := GetRunnerID(db, first, last, year, gender, team_id)
	if err == sql.ErrNoRows {
		 insert := `INSERT INTO runners (first_name, last_name, year, team_id, gender) VALUES ($1, $2, $3, $4, $5);`
		 _, err := db.Exec(insert, first, last, year, team_id, gender)
		 check(err)

		 id, err = GetRunnerID(db, first, last, year, gender, team_id)
		 check(err)
	}
	return id
}

func AddYearToRunner(db *sql.DB, id int, year string) int {
	update := `UPDATE runners SET year=$1 WHERE id=$2;`

	res, err := db.Exec(update, year, id)
	check(err)
	count, err := res.RowsAffected()
	check(err)
	if count == 0 { return -1 }
	return id
}

func ConnectRunnerTeam(db *sql.DB, runner_id, team_id int)  {
	update := `UPDATE runners SET team_id=$1 WHERE id=$2;`
	_, err := db.Exec(update, team_id, runner_id)
	check(err) 
}

func DeleteRunner(db *sql.DB, first, last string) {
	delete := `DELETE FROM runners WHERE (first_name=$1 AND last_name=$2);`
	_, err := db.Exec(delete, first, last)
	check(err)
}

type Year int

const (
	FR Year = 1 + iota
	SO
	JR
	SR
)

var years = [...]string{
	"FR",
	"SO",
	"JR",
	"SR",
	"N/A",
}

func (y Year) String() string {
	return years[y-1]
}