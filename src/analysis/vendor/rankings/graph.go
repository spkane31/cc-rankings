package rankings

import (
	"database/sql"
	"fmt"
	"time"
	"os"
	"math"

	_ "github.com/lib/pq"
)

var _, _ = os.Exit, math.Round

func AddToGraph(db *sql.DB, all_results *[]int, result, runner_id int, gender string) {

	debug := false
	for i := range *all_results {
		if CheckEdgeCondition(db, (*all_results)[i], result) {
			
			race_a, race_b, time_dif := GetEdgeInformation(db, (*all_results)[i], result)

			edge := CreateEdge(db, race_a, race_b, time_dif, runner_id, gender)

			// fmt.Printf("EDGES: %v\n", NumEdges(db))

			if debug {fmt.Println(edge)}

		}
	}
}

func CreateEdge(db *sql.DB, race_a, race_b int, time_dif float64, runner_id int, gender string) int {
	if race_a == race_b {
		return -1
	} else if race_a > race_b {
		race_a, race_b = race_b, race_a
	}
	var ret int
	
	// The from_race_id will always be the smaller one	
	query := `SELECT id FROM edges WHERE (from_race_id=$1 AND to_race_id=$2 AND runner_id=$3);`

	err := db.QueryRow(query, race_a, race_b, runner_id).Scan(&ret)
	
	if err == sql.ErrNoRows {
		// This edge does not exist, need to create one
		query := `INSERT INTO edges (from_race_id, to_race_id, runner_id, difference, gender, inserted_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`
		err = db.QueryRow(query, race_a, race_b, runner_id, time_dif, gender, time.Now(), time.Now()).Scan(&ret)
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

func BuildGraph(db *sql.DB, gender string, reg_dist, extra_dist int) *Graph {

	queryCenter := `SELECT id, course, distance, average, correction_avg FROM races WHERE is_base=$1 AND gender=$2;`

	var questionable_border float64 
	if gender == "MALE" {questionable_border = 180}
	if gender == "FEMALE" {questionable_border = 112.5}

	var center Race
	err := db.QueryRow(queryCenter, true, gender).Scan(&center.id, &center.course, &center.distance, &center.average, &center.correction_avg)
	check(err)

	// Once we have the center, we can build the interconnections out
	g := NewGraph()

	query := `select from_race_id, to_race_id, count(*), sum(difference) from edges group by from_race_id, to_race_id, gender HAVING count(*) > 6 and gender=$1;`
	rows, err := db.Query(query, gender)
	check(err)
	defer rows.Close()

	question_count := 0
	for rows.Next() {
		var from_race_id int
		var to_race_id int
		var count int
		var diff float64

		err = rows.Scan(&from_race_id, &to_race_id, &count, &diff)
		check(err)

		e := diff / float64(count)

		if math.Abs(e) > questionable_border {
			question_count++
		}

		if math.Abs(e) < 400 {
			err = g.AddEdge(from_race_id, to_race_id, e)
			check(err)
		}

	}

	fmt.Printf("Questionable edges: %0.6f percent\n", 100.0 * float64(question_count) / float64(g.Length()))
	return g
}

func FindCorrections(g *Graph, base_id int, db *sql.DB) {
	// v := g.GetIthVertex(0)
	// fmt.Println(v)

	// base_id := 1010
	// fmt.Printf("Base is Vertex %v\n", base_id)

	g.Completeness(base_id)

	g.ShortestPaths(base_id, db)
}

func NumEdges(db *sql.DB) (ret int) {
	query := `SELECT count(*) from edges;`
	err := db.QueryRow(query).Scan(&ret)
	check(err)
	return
}

// func (g *Graph) UpdateRace(db *sql.DB, id int, correction float64) (err error) {
// 	update := `UPDATE races SET correction_graph=$2 WHERE id=$1;`
// 	_, err = db.Exec(update, id, correction)
// 	check(err)

// 	fmt.Println(id, correction)
// 	os.Exit(1)
// 	return
// } 