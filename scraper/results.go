package main

import (
	"strings"
	"os"
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strconv"
	// "fmt"

	"github.com/PuerkitoBio/goquery"
)

func RemoveLeadingSpaces(s string) string {
	// ret := s
	count := 0
	for i := range s {
		if s[i] == 32 || s[i] == 10 {
			count++
			// ret = ret[i:len(ret)]
		} else {
			return s[count-1:len(s)]
		}
	}

	return s[count:len(s)]
}

func ScrapeResults(doc *goquery.Document) (*[][]string, *[][]string) { 
	var m_results [][]string
	var w_results [][]string
	var name string
	var year string
	var time string
	var team string
	sel := doc.Find(".color-xc")
	if len(sel.Nodes) > 4 {
		return &m_results, &w_results
	}
	womens_results := sel.Eq(1)
	row := womens_results.Find("tr")
	for i := range row.Nodes {
		cells := row.Eq(i).Find("td")

		// The individual cells start with a '\n', and have an extra space at the end, filtering
		// this out for neater CSV files.
		name = cells.Eq(1).Text()
		name = RemoveLeadingSpaces(name)
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
		var last string
		var first string
		if len(n) == 1 {

		} else {
			last = n[0]
			first = n[1][1:]
		}
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
	data := make(map[string]string)
	data["mens_results"] = filepath.Join(path, "mens.csv")
	data["womens_results"] = filepath.Join(path, "womens.csv")
	data["mens_count"] = strconv.Itoa(len(*mens))
	data["womens_count"] = strconv.Itoa(len(*womens))
	data["course"] = course
	data["date"] = date
	data["name"] = name
	WriteJSON(data, filepath.Join(path, "raceSummary.json"))
}

func WriteJSON(Summary map[string]string, path string) {
	file, _ := json.MarshalIndent(Summary, "", "  ")
	_ = ioutil.WriteFile(path, file, 0644)
}