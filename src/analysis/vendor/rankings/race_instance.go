package rankings

import (
	"database/sql"
	"time"
	"fmt"
	// "strconv"

	_ "github.com/lib/pq"
)

type Instance struct {
	id int
	date string
	race_id int 
	average float64
	std_dev float64
}

func GetInstance(db *sql.DB, date string, race_id int) (int, error) {
	var id int
	queryStatement := `SELECT id FROM race_instances WHERE (date=$1 AND race_id=$2);`
	row := db.QueryRow(queryStatement, date, race_id)
	err := row.Scan(&id)
	return id, err
}

func AddInstance(db *sql.DB, date string, race_id int, gender, distance string) int {
	var id int
	valid := false
	if gender == "MALE" {
		if (distance == "8000" || distance == "10000") {valid = true}
	} else if gender == "FEMALE" {
		if (distance == "6000" || distance == "5000") {valid = true}
	} else {
		fmt.Println("NO GENDER")
		return -1
	}
	id, err := GetInstance(db, date, race_id)
	if err == sql.ErrNoRows {
		sqlStatement := `INSERT INTO race_instances 
											(date, race_id, inserted_at, updated_at, valid) 
											VALUES ($1, $2, $3, $4, $5);`
		_, err = db.Exec(sqlStatement, date, race_id, time.Now(), time.Now(), valid)

		id, err = GetInstance(db, date, race_id)
	}
	return id
}

func FindAllInstances(db *sql.DB, id int) *[]Instance {
	var ret []Instance

	query := `SELECT * from race_instances WHERE race_id=$1;`
	rows, err := db.Query(query, id)
	check(err)
	defer rows.Close()

	for rows.Next() {
		var inst Instance
		err = rows.Scan(&inst.id, &inst.date, &inst.race_id, &inst.average, &inst.std_dev)
		ret = append(ret, inst)
	}

	return &ret
}

func GetInstanceResults(db *sql.DB, id int) *[]Result {
	var ret []Result

	query := `SELECT * FROM results WHERE race_instance_id=$1;`
	rows, err := db.Query(query, id)
	check(err)
	defer rows.Close()

	for rows.Next() {
		var result Result
		err = rows.Scan(&result.id, &result.distance, &result.unit, &result.rating, &result.time, &result.race_instance_id, &result.runner_id, &result.scaled_time, &result.time_float)
		check(err)
		ret = append(ret, result)
	}

	return &ret
}
