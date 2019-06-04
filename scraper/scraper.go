package main

import (
	"fmt"
	"net/http"
	"log"
	"net/url"
	"strconv"
	"strings"
	"unicode"
	// "os"

	"github.com/PuerkitoBio/goquery"
)

// var count int
var links []string

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func processElement(index int, element *goquery.Selection) {
	href, exists := element.Attr("href")
	if exists {
		// fmt.Println(href)
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
	
	fmt.Println(len(sel.Nodes))
	
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
	log.Printf("https://wwww.tfrrs.org" + link)
	
	response, err := http.Get("https://www.tfrrs.org" + link)
	check(err)
	defer response.Body.Close()
	
	document, err := goquery.NewDocument("https://www.tfrrs.org" + link)
	check(err)
	
	// This will get the title, ie. the Race Name
	sel := document.Find("h3 .white-underline-hover")
	name := GetRaceName(sel)
	log.Println("Scraping ", name)
	sel = document.Find("div .panel-heading-normal-text")
	date, course := GetRaceDate(sel)
	fmt.Printf("\nCourse: %v\nDate: %v", course, date)
}


func main() {
	
	log.Println("Scraping TFRRS!")
	GetUrlMonthYear(11, 2018)
	log.Printf("Found %d Links!", len(links))
	ScrapePage(links[1])
}