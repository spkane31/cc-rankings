package rankings

import (
	"testing"
	"database/sql"
	_ "github.com/lib/pq"
)

func TestGetTime(t *testing.T) {
		 time := GetTime("1:01:01")
		 if time != (61*60 + 1) {
			t.Errorf("t = %v; want 3661", time)
		 }
}

func TestGetTime2(t *testing.T) {
	time := GetTime("1:01:01.5")
	if time != (61.0*60.0 + 1.5) {
	 t.Errorf("t = %v; want 3661.5", time)
	}
}

func FindResultsForRunner(db *sql.DB, id int) *[]Result {
	var res []Result

	queryStatement := `SELECT * FROM results WHERE runner_id=$1 ORDER BY race_instance_id;`

	rows, err := db.Query(queryStatement, id)
	check(err)
	defer rows.Close()

	for rows.Next() {
		var result Result
		err = rows.Scan(&result.id, &result.distance, &result.unit, &result.rating, &result.time, &result.race_instance_id, &result.runner_id, &result.scaled_time, &result.time_float)
		res = append(res, result)
	}

	return &res
}