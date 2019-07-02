package main

import (
	"database/sql"
	"fmt"
	"os"
	// "strings"
	// "strconv"
	"math"
	"time"
	_ "github.com/lib/pq"
)

var _, _, _, _ = fmt.Println, os.Exit, math.Round, time.Now


type Pair struct {
	to, from int
}

type Edge struct {
	to, from int
	count int
	total_time float64
}

func AddToGraph(results *[]Result, g *map[Pair]*Edge) {
	for i := 0; i < len(*results)-1; i++ {
		for j := i+1; j < len(*results); j++ {

			if CheckDistances((*results)[i], (*results)[j], "WOMENS") {
				p := Pair{ (*results)[i].race_instance_id, (*results)[j].race_instance_id }
				e, has := (*g)[p]
				
				if !has {
					var e Edge
					// New Edge, define basic parts
					e.to = (*results)[i].race_instance_id
					e.from = (*results)[j].race_instance_id
					e.count = 1
					e.total_time = GetTime((*results)[i].time) - GetTime((*results)[j].time)

					// Insert the edge
					(*g)[p] = &e
				} else {
					// Already exists, increase count, add to total time
					e.count += 1
					e.total_time += GetTime((*results)[i].time) - GetTime((*results)[j].time)
				}
			}

		}
	}
}

func PrintGraph(db *sql.DB, g *map[Pair]*Edge) {

	for p, e := range (*g) {
		PrintEdgeNames(db, p.to, p.from)
		w := fmt.Sprintf("%.3f", e.total_time / float64(e.count))
		fmt.Printf("Connections: %v\t Weight: %v\n", e.count, w)
	}

}

func CheckDistances(a, b Result, gender string) bool {
	var ret bool = false

	if gender == "MENS" {
		if a.distance == 10000 || a.distance == 8000 {
			if b.distance == 10000 || b.distance == 80000 {
				return true
			}
		}
	}

	if gender == "WOMENS" {
		if a.distance == 6000 || a.distance == 5000 {
			if b.distance == 6000 || b.distance == 5000 {
				return true
			}
		}
	}

	return ret
}

func PrintEdgeNames(db *sql.DB, i, j int) {
	instanceStatement := `SELECT race_id FROM race_instances WHERE id=$1;`
	row := db.QueryRow(instanceStatement, i)
	var race_id int
	var to string
	var from string
	err := row.Scan(&race_id)
	check(err)
	raceStatement := `SELECt name FROM races WHERE id=$1;`
	row = db.QueryRow(raceStatement, race_id)
	err = row.Scan(&to)
	check(err)

	row = db.QueryRow(instanceStatement, j)
	err = row.Scan(&race_id)
	check(err)
	row = db.QueryRow(raceStatement, race_id)
	err = row.Scan(&from)
	check(err)

	fmt.Printf("\n%v -> %v\n", to[0:len(to)-1], from[0:len(from)-1])
}