package main

import (
	"math/rand"
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
	"os"
	"strings"
	"sort"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	// "gonum.org/v1/gonum/stat"
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

		lin_reg := CreateLinearRegression(&results)

		data := CreateXYData(&results)
		os.Mkdir("plots", os.ModePerm)

		safe_name := strings.Replace((*races)[i].name, " ", "", -1)
		safe_name = strings.Replace(safe_name, "/", "", -1)

		f_name := "plots/" + safe_name
		fmt.Printf("Creating %v plot\n", (*races)[i].name)
		Plot(data, lin_reg, (*races)[i].name, f_name)
	}
}

func CreateLinearRegression(results *[]Result) plotter.XYs {
	size := len(*results)
	sum_x, sum_y, sum_xx, sum_xy := 0.0, 0.0, 0.0, 0.0

	mean := 0.0
	S := 0.0

	for i := range *results {	
		sum_x += float64(i+1)
		sum_xx += float64((i+1) * (i+1))

		t := GetTime((*results)[i].time)

		sum_y += t
		sum_xy += (float64(i+1) * t)
		prev_mean := mean
		mean = mean + (t - mean) / float64(i+1)
		S = S + (t - mean) * (t - prev_mean)
	}

	m := (float64(size) * sum_xy - sum_x*sum_y) / (float64(size) * sum_xx - sum_x * sum_x)
	b := (sum_y / float64(size)) - (m * sum_x / float64(size))

	fmt.Printf("m = %v, b = %v, mean = %v, st_dev = %v\n", m, b, mean, math.Sqrt(S / float64(size)))

	pts := make(plotter.XYs, size)

	for i := range *results {	
		pts[i].X = float64(i)
		pts[i].Y = m*float64(i) + b
	}
	return pts
}

func Plot(data, lin_reg plotter.XYs, title, saveFile string) {
	p, err := plot.New()
	check(err)

	p.Title.Text = title
	p.X.Label.Text = "Place"
	p.Y.Label.Text = "Time (seconds)"
	err = plotutil.AddLinePoints(p,
		title, data,
		"Best Fit", lin_reg,
	)
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

type xy struct {
	x []float64
	y []float64
}

func (d xy) Len() int {
	return len(d.x)
}

func (d xy) XY(i int) (x, y float64) {
	x = d.x[i]
	y = d.y[i]
	return
}