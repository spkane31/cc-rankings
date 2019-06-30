package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

var _ = os.Exit

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	fmt.Printf("%v: Establishing DB connection...\n", time.Now().Format("01-02-2006, 15:04:05"))
	db := ConnectToPSQL()

	ComputeAverages(db)
	UpdateCorrectionAvg(db)
	os.Exit(1)

	g := make(map[Pair]*Edge)

	FindAllConnections(db, &g)

	PlotAllRaces(db)

	fmt.Printf("%v: Finished!\n", time.Now().Format("01-02-2006, 15:04:05"))
}