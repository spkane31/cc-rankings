package main

import (
	"fmt"
	"log"
	"os"
)

var _ = os.Exit

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	fmt.Println("Establishing DB connection...")
	db := ConnectToPSQL()

	PlotAllRaces(db)
	os.Exit(1)

	FindAllConnections(db)
	log.Println("Finished")
}