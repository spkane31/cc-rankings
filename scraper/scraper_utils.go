package main

import (
	"strings"
	"os"
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"fmt"
	"net/http"
	"log"
	"net/url"
	"unicode"
	"reflect"

	"github.com/PuerkitoBio/goquery"
)

var _, _, _ = fmt.Println, reflect.TypeOf, unicode.IsSymbol


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

func DetermineRaces(str string) []string {
	// Takes the raw string and returns a list of the distances and genders
	var ret []string

	raw := strings.Split(str, " ")
	for i := 0; i < len(raw); i++ {
		if len(raw[i]) > 0 {
			ret = append(ret, strings.ToUpper(raw[i]))
		}
	}

	var final []string
	for i := range ret {
		var s strings.Builder
		for _, j := range ret[i] {
			if unicode.IsNumber(j) || unicode.IsLetter(j) {
				s.WriteString(string(j))
			}
		}
		final = append(final, s.String())
	}
	final = final[1:]

	ret = []string{}
	for i := 0; i < len(final)-1; i+=2 {
		var temp strings.Builder
		temp.WriteString(final[i])
		temp.WriteString(" ")
		temp.WriteString(final[i+1])
		ret = append(ret, temp.String())
	}
	return ret
}

func ScrapeResults(doc *goquery.Document) (*[][]string, *[][]string, *[]string) { 
	var m_results [][]string
	var w_results [][]string
	var name string
	var year string
	var time string
	var team string

	races_elem := doc.Find("ol")
	races := races_elem.Text()

	r := DetermineRaces(races)

	sel := doc.Find(".color-xc")
	if len(sel.Nodes) > 4 {
		return &m_results, &w_results, &r
	}
	
	var womens_results *goquery.Selection
	var mens_results *goquery.Selection
	fmt.Println(r)
	if r[0][:5] == "WOMEN" {
		womens_results = sel.Eq(1)
		mens_results = sel.Eq(3)
	} else {
		// Mens
		womens_results = sel.Eq(3)
		mens_results = sel.Eq(1)
	}

	womens_results = sel.Eq(1)
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

	mens_results = sel.Eq(3)
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
	
	return &m_results, &w_results, &r
}

func WriteResults(mens, womens *[][]string, name, date, course string, races *[]string) {
	path := filepath.Join(HomePath, "RaceResults")
	path = filepath.Join(path, strings.Replace(name[0:len(name)-1], " ", "", -1))
	
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


	var mens_dist string
	var womens_dist string
	if (*races)[0][0:5] == "WOMEN" && len(*races) > 1{
		womens_dist = strings.Split((*races)[0], " ")[1]
		mens_dist = strings.Split((*races)[1], " ")[1]
	} else if len(*races) > 1 {
		womens_dist = strings.Split((*races)[1], " ")[1]
		mens_dist = strings.Split((*races)[0], " ")[1]
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
	data["mens_distance"] = mens_dist
	data["womens_distance"] = womens_dist
	WriteJSON(data, filepath.Join(path, "raceSummary.json"))
}

func WriteJSON(Summary map[string]string, path string) {
	file, _ := json.MarshalIndent(Summary, "", "  ")
	_ = ioutil.WriteFile(path, file, 0644)
}

func processElement(index int, element *goquery.Selection) {
	href, exists := element.Attr("href")
	if exists {
		ParseUrl(href)
	}
}

func ParseUrl(s string) {
	parsedUrl, err := url.Parse(s)
	check(err)
	path := parsedUrl.Path

	if strings.Contains(path, "/results/") {
		links = append(links, parsedUrl.Path)	
	}
}

func GetUrlMonthYear(month, year int) {
	y := strconv.Itoa(year)
	m := strconv.Itoa(month)

	response, err := http.PostForm(
		"https://www.tfrrs.org/results_search.html",
		url.Values{
			"sport": {"xc"},
			"state": {""},
			"month": {m},
			"year": {y},
		},
	)

	check(err)
	defer response.Body.Close()
	document, err := goquery.NewDocumentFromReader(response.Body)
	check(err)
	document.Find("a").Each(processElement)
}

func RemoveLeadTrailSpaces(s string) string {
	var ret string
	for i, r := range s {
		if !unicode.IsSpace(r) {
			ret = s[i:]
			break
		}
	}
	// The last two characters are spaces, I think
	return string(ret[0:len(ret)-3])
}

func TrimString(s string) string {
	var ret string
	for _, r := range s {
		if r != 10 {
			ret = ret + string(r)
		} else {
			ret += " "
		}
	}
	RemoveLeadTrailSpaces(ret)
	return ret
}

func GetRaceName(sel *goquery.Selection) string {
	// The race name is padded with spaces, and I want to remove those in the beginning 
	// and at the end which required some odd programming
	var trimmed string
	var letterOcurred bool
	var index int
	
	// fmt.Println(len(sel.Nodes))
	
	for i := range sel.Nodes {
		name := sel.Eq(i).Text()
		
		for i, r := range name {
			
			if unicode.IsLetter(r) {
				trimmed = trimmed + string(r)
				letterOcurred = true
			} else if unicode.IsNumber(r) {
				trimmed = trimmed + string(r)
				letterOcurred = true
			} else if unicode.IsSpace(r) {
				// If there has already been a letter, and now there's a double space, get rid of this
				// TODO - Look to see if there's a way to get the race name from the links page
				if index + 1 == i && letterOcurred{
					break
				}
				index = i
				if letterOcurred {
					trimmed = trimmed + string(r)
				}
			}
		}
		
	}
	return trimmed
}

func GetRaceDate(sel *goquery.Selection) (string, string) {
	var date string
	var course string
	date = sel.Eq(0).Text()
	course = sel.Eq(2).Text()

	course = TrimString(course)
	
	return date, course
}

func ScrapePage(link string) {
	// This function will scrape a race page and create the CSV, and the JSON File
	// log.Printf("https://wwww.tfrrs.org" + link)
	
	response, err := http.Get("https://www.tfrrs.org" + link)

	check(err)
	defer response.Body.Close()
	
	document, err := goquery.NewDocument("https://www.tfrrs.org" + link)
	check(err)

	mens_results, womens_results, races := ScrapeResults(document)
	if len(*mens_results) == 0 || len(*womens_results) == 0 {
		// This is a scenario with funky formatting on the page, for now I am ignoring it
		return 
	}
	
	// This will get the title, ie. the Race Name
	sel := document.Find("h3 .white-underline-hover")
	name := GetRaceName(sel)
	log.Println("Scraping ", name)
	sel = document.Find("div .panel-heading-normal-text")
	date, course := GetRaceDate(sel)

	WriteResults(mens_results, womens_results, name, date, course, races)
}
