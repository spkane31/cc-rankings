package main

import (
	"strings"
	"fmt"
	"os"
	"encoding/csv"
	"path/filepath"

	"github.com/PuerkitoBio/goquery"
)

type RaceSummary struct {
	name 						string `json:"name"`
	date 						string `json:"date"`
	course			 		string `json:"course"`
	mens_results 		string `json:"mens_results"`
	womens_results 	string `json:"womens_results"`
	mens_count 			int `json:"mens_count"`
	womens_count 		int `json:"womens_count"`
}

func ScrapeResults(doc *goquery.Document) (*[][]string, *[][]string) { 
	var m_results [][]string
	var w_results [][]string
	var name string
	var year string
	var time string
	var team string
	sel := doc.Find(".color-xc")

	womens_results := sel.Eq(1)
	row := womens_results.Find("tr")
	for i := range row.Nodes {
		cells := row.Eq(i).Find("td")

		// The individual cells start with a '\n', and have an extra space at the end, filtering
		// this out for neater CSV files.
		name = cells.Eq(1).Text()
		name = name[1:len(name)-1]
		year = cells.Eq(2).Text()
		year = year[1:len(year)-1]
		team = cells.Eq(3).Text()
		team = team[1:len(team)-1]
		time = cells.Eq(5).Text()
		time = time[1:len(time)-1]
		
		// The name is "last, first". Turning this into two different vars
		n := strings.Split(name, ",")
		last := n[0]
		first := n[1][1:]
		e := []string{last, first, team, year, time}
		w_results = append(w_results, e)

	}

	mens_results := sel.Eq(3)
	row = mens_results.Find("tr")
	for i := range row.Nodes {
		cells := row.Eq(i).Find("td")

		// The individual cells start with a '\n', and have an extra space at the end, filtering
		// this out for neater CSV files.
		name = cells.Eq(1).Text()
		name = name[1:len(name)-1]
		year = cells.Eq(2).Text()
		year = year[1:len(year)-1]
		team = cells.Eq(3).Text()
		team = team[1:len(team)-1]
		time = cells.Eq(5).Text()
		time = time[1:len(time)-1]
		
		// The name is "last, first". Turning this into two different vars
		n := strings.Split(name, ",")
		last := n[0]
		first := n[1][1:]
		e := []string{last, first, team, year, time}
		m_results = append(m_results, e)

	}

	// records := [][]string{
	// 	{last, first, year, team, time},
	// // 	{last, first, year, team, time},
	// // }

	// w := csv.NewWriter(os.Stdout)

	// for _, record := range(records) {
	// 	if err := w.Write(record); err != nil {
	// 		log.Fatalln("error writing record to csv:", err)
	// 	}
	// }

	// w.Flush()

	// if err := w.Error(); err != nil {
	// 	log.Fatal(err)
	// }
	return &m_results, &w_results
}

func WriteResults(mens, womens *[][]string, name, date, course string) {
	path := filepath.Join(HomePath, "RaceResults")
	path = filepath.Join(path, strings.Replace(name[0:len(name)-1], " ", "", -1))
	// path := HomePath + "RaceResults" + strings.Replace(name, " ", "", -1)
	// fmt.Println(path)
	os.MkdirAll(path, os.ModePerm)

	m_file, err := os.Create(filepath.Join(path, "mens.csv"))
	check(err)
	w_file, err := os.Create(filepath.Join(path, "womens.csv"))
	check(err)

	writer := csv.NewWriter(m_file)
	defer writer.Flush()

	for _, value := range *mens {
		err := writer.Write(value)
		check(err)
	}
	
	writer = csv.NewWriter(w_file)
	defer writer.Flush()

	for _, value := range *mens {
		err := writer.Write(value)
		check(err)
	}

	// Now to create the raceSummary.json file with more details

	Summary := RaceSummary{
		mens_results: filepath.Join(path, "mens.csv"),
		womens_results: filepath.Join(path, "womens.csv"),
		mens_count: len(*mens),
		womens_count: len(*womens),
		course: course,
		date: date,
		name: name,
	}
	fmt.Println(Summary)

}