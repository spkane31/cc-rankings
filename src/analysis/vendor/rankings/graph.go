package rankings

import (
	"database/sql"
	"fmt"
	"time"
	// "os"
	// "strings"
	// "strconv"
	// "math"

	_ "github.com/lib/pq"
)

func AddToGraph(db *sql.DB, all_results *[]int, result int) {

	debug := false
	for i := range *all_results {
		if (*all_results)[i] == result {

		} else 
		if CheckResultsYears(db, (*all_results)[i], result) {
			race_a, race_b, time_dif := GetEdgeInformation(db, (*all_results)[i], result)

			edge := UpdateEdge(db, race_a, race_b, time_dif)
			if debug {fmt.Println(edge)}

		}
	}

	MarkResultAsAdded(db, result)

}

func UpdateEdge(db *sql.DB, race_a, race_b int, time_dif float64) int {
	if race_a == race_b {
		return -1
	} else if race_a > race_b {
		race_a, race_b = race_b, race_a
	}

	var ret int
	
	// The from_race_id will always be the smaller one	
	query := `SELECT id, count, total_time FROM edges WHERE (from_race_id=$1 AND to_race_id=$2);`
	id, count, total_time := 0, 0, 0.0
	row := db.QueryRow(query, race_a, race_b)
	err := row.Scan(&id, &count, &total_time)
	if err == sql.ErrNoRows || id == 0 {
		// This edge does not exist, need to create one
		query := `INSERT INTO edges (from_race_id, to_race_id, total_time, count, inserted_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`
		err = db.QueryRow(query, race_a, race_b, time_dif, 1, time.Now(), time.Now()).Scan(&ret)
		check(err)
	} else {
		// This edge does exist, update it
	
		update := `UPDATE edges SET count=$2, total_time=$3 WHERE id=$1 RETURNING id;`
		count++
		total_time += time_dif
		err = db.QueryRow(update, id, count, total_time).Scan(&ret)
		check(err)		
	}

	return ret
}