package main

import (
	"math/rand"
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
	"os"
	"strings"
	"sort"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

var _, _ = fmt.Println, os.Exit

func PlotAllRaces(db *sql.DB) {
	races := GetAllRacesByGender(db, "MENS")
	for i := range *races {
		insts := FindAllInstances(db, (*races)[i].id)
		var results []Result
		for i := range *insts {
			r := GetInstanceResults(db, (*insts)[i].id)
			for _, each := range *r {
				results = append(results, each)
			}
		}

		sort.SliceStable(results, func(i, j int) bool { return GetTime(results[i].time) < GetTime(results[j].time) })

		results = FilterDNFs(&results)

		for _, each := range results {
			if GetTime(each.time) == 0 {fmt.Println(each)}
		}

		data := CreateXYData(&results)
		os.Mkdir("plots", os.ModePerm)

		safe_name := strings.Replace((*races)[i].name, " ", "", -1)
		safe_name = strings.Replace(safe_name, "/", "", -1)

		f_name := "plots/" + safe_name
		fmt.Printf("Creating %v plot\n", (*races)[i].name)
		Plot(data, (*races)[i].name, f_name)
	}
}

func Plot(data plotter.XYs, title, saveFile string) {
	p, err := plot.New()
	check(err)

	p.Title.Text = title
	p.X.Label.Text = "Place"
	p.Y.Label.Text = "Time (seconds)"
	err = plotutil.AddLinePoints(p,
		title, data,
	)
	check(err)
	check(err)

	err = p.Save(8*vg.Inch, 8*vg.Inch, saveFile+".png")
	check(err)

}




func PlotRace(db *sql.DB) {
	name := "NCAA Division I Cross Country Championships"
	course := "E.P.TomSawyerParkLousivilleKY"
	// distance := 10000

	race := FindRaceByCourseName(db, course, name, "MENS")
	insts := FindAllInstances(db, race.id)
	results := GetInstanceResults(db, (*insts)[0].id)

	data := CreateXYData(results)
	
	p, err := plot.New()
	check(err)

	p.Title.Text = "NCAA DI Cross Country Championships 2015"
	p.X.Label.Text = "Place"
	p.Y.Label.Text = "Time (seconds)"

	err = plotutil.AddLinePoints(p,
				"NCAA DI", data,
				// "Second", randomPoints(15),
				// "Third", randomPoints(15),
	)
	check(err)

	err = p.Save(8*vg.Inch, 8*vg.Inch, "plot.png")
	check(err)
}

func CreateXYData(results *[]Result) plotter.XYs {
	pts := make(plotter.XYs, len(*results))

	for i := range *results {
		pts[i].X = float64(i+1)
		pts[i].Y = GetTime((*results)[i].time)
	}

	return pts
}

func randomPoints(n int) plotter.XYs {
	pts := make(plotter.XYs, n)
	for i := range pts{
		if i == 0 {
			pts[i].X = rand.Float64()
		} else {
			pts[i].X = pts[i-1].X + rand.Float64()
		}
		pts[i].Y = pts[i].Y + 10 * rand.Float64()
	}
	return pts
}