package rankings

import (
	// "math/rand"
	"fmt"
	"os"
	"database/sql"
	_ "github.com/lib/pq"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	// "gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	// "gonum.org/v1/gonum/stat"
	// "github.com/gonum/stat/distuv"
)

var _ = os.Exit

func MakePlots(db *sql.DB) {
	PlotHistogramEdges(db)
}

func PlotHistogramEdges(db * sql.DB) {
	query := `select count, total_time FROM edges WHERE count > 10;`
	rows, err := db.Query(query)
	check(err)
	defer rows.Close()

	num_bins := 100
	v := make(plotter.Values, num_bins)
	vals := make([]float64, 5986)
	i := 0
	max := 0.0
	min := 100.0
	for rows.Next() {
		var count int
		var total_time float64
		var weight float64
		err = rows.Scan(&count, &total_time)
		check(err)
		weight = total_time / float64(count)

		vals[i] = weight
		// fmt.Println(v[i])
		i++
		if weight > max {max = weight}
		if weight < min {min = weight}
	}

	fmt.Println(min, max)
	r := max - min
	bin := int(r / float64(num_bins))
	fmt.Printf("BIN WIDTH: %v\n", bin)
	bin_x_values := make([]int, num_bins)
	bin_x_values[0] = int(min)

	for i := 1; i < len(bin_x_values); i++ {
		bin_x_values[i] = bin_x_values[i-1] + bin
	}

	fmt.Println(bin_x_values)
	
	for _, val := range vals{
		fmt.Printf("Val: %v\tVal-min: %v\n", val, val-min)
		// b := int(val * r / float64(num_bins)) //int(val - min) % bin

		for i, x := range bin_x_values {
			if int(val) < x {
				v[i]++
				break
			}
		}
		// fmt.Println(b)
		// if b > num_bins || b < 0 {
		// 	fmt.Errorf("ERROR")
		// 	os.Exit(1)
		// }
		// fmt.Println(v)
	}

	p, err := plot.New()
	check(err)
	p.Title.Text = "Histogram of weights"

	h, err := plotter.NewHist(v, 16)
	check(err)

	// h.Normalize(1)
	p.Add(h)

	err = p.Save(8 * vg.Inch, 8 * vg.Inch, "hist.png")
	check(err)
}