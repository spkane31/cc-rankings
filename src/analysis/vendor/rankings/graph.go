package rankings

import (
	"database/sql"
	"fmt"
	"time"
	"os"
	// "strings"
	// "strconv"
	// "math"

	_ "github.com/lib/pq"
)

var _ = os.Exit

func AddToGraph(db *sql.DB, all_results *[]int, result int) {

	debug := false
	for i := range *all_results {
		if CheckEdgeCondition(db, (*all_results)[i], result) {
			race_a, race_b, time_dif := GetEdgeInformation(db, (*all_results)[i], result)

			edge := UpdateEdge(db, race_a, race_b, time_dif)
			if debug {fmt.Println(edge)}

			if len(*all_results) == 2 {
				MarkResultAsAdded(db, (*all_results)[0])
			}

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

func CheckEdgeCondition(db *sql.DB, result_a, result_b int) bool {
	if result_a == result_b {
		return false
	}

	if !CheckResultsYears(db, result_a, result_b) {
		return false
	}

	query := `SELECT gender, distance, time_float FROM results WHERE id=$1;`

	var gender_a string
	var dist_a int
	var time_a float64
	err := db.QueryRow(query, result_a).Scan(&gender_a, &dist_a, &time_a)
	check(err)

	var gender_b string
	var dist_b int
	var time_b float64
	err = db.QueryRow(query, result_b).Scan(&gender_b, &dist_b, &time_b)
	check(err)

	if time_b == 0 || time_a == 0 {
		return false
	}
	
	if gender_a != gender_b {
		return false
	} else if gender_a == "MALE" {
		if dist_a == 8000 || dist_a == 10000 {
			if dist_b == 8000 || dist_b == 10000 {
				return true
			} else {
				return false
			}
		} else {
			return false
		}
	} else if gender_b == "FEMALE" {
		if dist_a == 5000 || dist_a == 6000 {
			if dist_b == 5000 || dist_b == 6000 {
				return true
			} else {
				return false
			}
		} else {
			return false
		}
	}

	return false
}

func BuildGraph(db *sql.DB) *Graph {
	queryCenter := `SELECT id, course, distance, average, correction_avg FROM races WHERE is_base=$1 AND gender=$2;`

	var center Race
	err := db.QueryRow(queryCenter, true, "MALE").Scan(&center.id, &center.course, &center.distance, &center.average, &center.correction_avg)
	check(err)

	// Once we have the center, we can build the interconnections out
	g := NewGraph()

	query := `SELECT id FROM races WHERE gender=$1 AND (distance=$2 OR distance=$3);`
	rows, err := db.Query(query, "MALE", 8000, 10000)
	check(err)
	defer rows.Close()

	edqeQuery := `SELECT to_race_id, total_time, count FROM edges WHERE from_race_id=$1;`

	for rows.Next() {
		var from_race_id int
		err = rows.Scan(&from_race_id)
		check(err)

		edges, err := db.Query(edqeQuery, from_race_id)
		check(err)

		from_vertex, err := g.GetVertex(from_race_id)
		check(err)

		for edges.Next() {
			var to_race_id int
			var total_time float64
			var weight float64

			err = rows.Scan(&to_race_id, &total_time, &weight)
			check(err)

			to_vertex, err := g.GetVertex(to_race_id)
			if err != nil {
				to_vertex = g.AddVertex(from_race_id)
			}


			err = g.AddEdge(from_race_id, to_race_id, )

		}

	}
	

	return g
}