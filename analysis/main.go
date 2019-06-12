package main

import (
	"fmt"
	"log"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	fmt.Println("Establishing DB connection...")
	db := ConnectToPSQL()

	id := AddRunner("Sean", "Kane", db)
	DeleteRunner(db, "sean", "kane")
	fmt.Println(id)
	team := AddTeam(db, "UC RC")
	fmt.Println(team)
}