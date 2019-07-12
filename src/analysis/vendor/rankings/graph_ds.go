package rankings

import (
	"fmt"
	"math"
	"os"
	"database/sql"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	// "github.com/gonum/stat/distuv"
	_ "github.com/lib/pq"
)

var _ = os.Exit

type Vertex interface {

	ID() int

	Edges() []Edge
}

type Edge interface {
	Get() (from int, to int, weight float64)
}

type Graph struct {
	vertices 	map[int]*vertex
	egress 		map[int]map[int]*edge
	ingress		map[int]map[int]*edge
}

type vertex struct {
	// self interface{}
	id 			int
	outTo		map[int]float64
	inFrom	map[int]float64
	enable 	bool
}

func (v *vertex) ID() int {
	return v.id
}

func (v *vertex) Edges() (edges []edge) {
	edges = make([]edge, len(v.outTo) + len(v.inFrom))
	i := 0
	for to, weight := range v.outTo {
		edges[i] = edge{v.id, to, weight}
		i++
	}
	for from, weight := range v.inFrom {
		edges[i] = edge{from, v.id, weight}
		i++
	}

	return
}

func NewVertex(id int) *vertex {
	out := make(map[int]float64)
	in := make(map[int]float64)
	return &vertex{id, out, in, true}
} 

type edge struct {
	// self interface{}
	from 		int
	to 			int
	weight 	float64
	// enable bool
	// changed bool
}

func (e *edge) Get() (int, int, float64) {
	return e.from, e.to, e.weight
}

func (e *edge) GetWeight() float64 {
	return e.weight
}

func NewGraph() *Graph {
	g := new(Graph)
	g.vertices = make(map[int]*vertex)
	g.egress = make(map[int]map[int]*edge)
	g.ingress = make(map[int]map[int]*edge)

	return g
}

func (g Graph) Length() int {
	return len(g.vertices)
}

func (g *Graph) GetVertex(id int) (vertex interface{}, err error) {
	if v, exists := g.vertices[id]; exists {
		vertex = v
		return
	}

	err = fmt.Errorf("Vertex %v is not found", id)
	return
}

func (g *Graph) GetEdge(from, to int) (interface{}, error) {
	if _, exists := g.vertices[from]; !exists {
		return nil, fmt.Errorf("Vertex(from) %v is not found", from)
	}

	if _, exists := g.vertices[to]; !exists {
		return nil, fmt.Errorf("Vertex(to) %v is not found", to)
	}

	if edge, exists := g.egress[from][to]; exists {
		return edge, nil
	}

	return nil, fmt.Errorf("Edge from %v to %v is not found", from, to)
}

func (g *Graph) GetEdgeWeight(from, to int) (float64, error) {

	if _, exists := g.vertices[from]; !exists {
		return math.Inf(1), fmt.Errorf("Vertex(from) %v is not found", from)
	}

	if _, exists := g.vertices[to]; !exists {
		return math.Inf(1), fmt.Errorf("Vertex(to) %v is not found", to)
	}

	if edge, exists := g.egress[from][to]; exists {
		return edge.weight, nil
	}

	return math.Inf(1), nil
}

func (g *Graph) AddVertex(id int) error {
	if _, exists := g.vertices[id]; exists {
		return fmt.Errorf("Vertex %v already exists", id)
	}

	g.vertices[id] = NewVertex(id)
	g.egress[id] = make(map[int]*edge)
	g.ingress[id] = make(map[int]*edge)

	return nil
}

func (g *Graph) AddEdge(from, to int, weight float64) error {
	if weight == math.Inf(1) {

		return fmt.Errorf("-inf weight is reserved for internal usage")
	}

	if _, exists := g.vertices[from]; !exists {
		g.AddVertex(from)
		// return fmt.Errorf("Vertex(from) %v is not found", from)
	}

	if _, exists := g.vertices[to]; !exists {
		g.AddVertex(to)
		// return fmt.Errorf("Vertex(to) %v is not found", to)
	}

	if _, exists := g.egress[from][to]; exists {
		return fmt.Errorf("Edge from %v to %v is duplicate", from, to)
	}

	g.egress[from][to] = &edge{from, to, weight}
	g.ingress[to][from] = g.egress[from][to]

	return nil
}

func (g *Graph) UpdateEdgeWeight(from, to int, weight float64) error {
	if weight == math.Inf(-1) {
		return fmt.Errorf("-inf weight is reserved for internal usage")
	}

	if _, exists := g.vertices[from]; !exists {
		return fmt.Errorf("Vertex(from) %v is not found.", from)
	}

	if _, exists := g.vertices[to]; !exists {
		return fmt.Errorf("Vertex(to) %v is not found.", from)
	}

	if edge, exists := g.egress[from][to]; exists {
		edge.weight = weight
		return nil
	}

	return fmt.Errorf("Edge from %v to %v is not found", from, to)
}

func (g *Graph) DeleteVertex(id int) interface{} {
	if vertex, exists := g.vertices[id]; exists {
		for to := range g.egress[id] {
			delete(g.ingress[to], id)
		}

		for from := range g.ingress[id] {
			delete(g.egress[from], id)
		}

		delete(g.egress, id)
		delete(g.ingress, id)
		delete(g.vertices, id)

		return vertex
	}

	return nil
}

func (g *Graph) DeleteEdge(from, to int) interface{} {
	if _, exists := g.vertices[from]; !exists {
		return nil
	}

	if _, exists := g.vertices[to]; !exists {
		return nil
	}

	if edge, exists := g.egress[from][to]; exists {
		delete(g.egress[from], to)
		delete(g.ingress[to], from)
		return edge
	}

	return nil
}

func (g *Graph) AddVertexWithEdges(v Vertex) error {
	if _, exists := g.vertices[v.ID()]; exists {
		return fmt.Errorf("Vertex %v is duplicate", v.ID())
	}

	g.vertices[v.ID()] = NewVertex(v.ID())
	g.egress[v.ID()] = make(map[int]*edge)
	g.ingress[v.ID()] = make(map[int]*edge)

	for _, eachEdge := range v.Edges() {
		from, to, weight := eachEdge.Get()

		if weight == math.Inf(-1) {
			return fmt.Errorf("-inf weight is reserved for internal usage")
		}
		if from != v.ID() && to != v.ID() {
			return fmt.Errorf("Edge from %v to %v is unrelated to vertex %v", from, to, v.ID())
		}

		if _, exists := g.egress[to]; !exists {
			g.egress[to] = make(map[int]*edge)
		}

		if _, exists := g.ingress[from]; !exists {
			g.egress[from] = make(map[int]*edge)
		}

		if _, exists := g.ingress[from]; !exists {
			g.ingress[from] = make(map[int]*edge)
		}

		if _, exists := g.ingress[to]; !exists {
			g.ingress[to] = make(map[int]*edge)
		}

		g.egress[from][to] = &edge{from, to, weight}
		g.ingress[to][from] = g.egress[from][to]
	}
	
	return nil
}

func (g *Graph) CheckIntegrity() error {
	for from, out := range g.egress {
		if _, exists := g.vertices[from]; !exists {
			return fmt.Errorf("Vertex %v is not found.", from)
		}
		for to := range out {
			if _, exists := g.vertices[to]; !exists {
				return fmt.Errorf("Vertex %v is not found", to)
			}
		}
	}

	for to, in := range g.ingress {
		if _, exists := g.vertices[to]; !exists {
			return fmt.Errorf("Vertex %v is not found", to)
		}

		for from := range in {
			if _, exists := g.vertices[from]; !exists {
				return fmt.Errorf("Vertex %v is not found", from)
			}
		}
	}

	return nil
}

func (g *Graph) GetPathWeight(path []int) (totalWeight float64) {
	if len(path) == 0 {
		return math.Inf(-1)
	}

	if _, exists := g.vertices[path[0]]; !exists {
		return math.Inf(-1)
	}

	for i := 0; i < len(path)-1; i++ {
		if _, exists := g.vertices[path[i+1]]; !exists {
			return math.Inf(-1)
		}

		if edge, exists := g.egress[path[i]][path[i+1]]; exists {
			totalWeight += edge.GetWeight()
		} else {
			return math.Inf(-1)
		}
	}

	return totalWeight
}

func (g Graph) Print() {
	for from, _ := range g.vertices {
		fmt.Printf("ID: %v Connects to: ", from)

		for to, e := range g.egress[from] {
			fmt.Printf("%v (%0.2f)\t", to, e.weight)
		}
		fmt.Printf("\n")
	}
}

func (g *Graph) PrintVertex(id int) {
	fmt.Printf("ID: %v Connects to: ", id)

	for to, e := range g.egress[id] {
		fmt.Printf("%v (%0.2f)\t", to, e.weight)
	}
	fmt.Printf("\n")
}

func (g *Graph) GetIthVertex(i int) *vertex {
	return g.vertices[i]
}

func (g *Graph) minDistance(dist map[int]float64, sptSet map[int]bool) int {
	min := math.Inf(1)
	min_index := -1

	for key, val := range dist {
		if sptSet[key] == false && math.Abs(val) <= min {
			min = dist[key]
			min_index = key
		}
	}

	return min_index
}

// Finds the shortest distances between two points 
// Since edges are seen as both positive and negative, the minimum between
// two edges is the values closest to zero (absolute value)
func (g *Graph) Dijkstra(source int) (dist map[int]float64, prev map[int]int, err error) {
	if _, exists := g.vertices[source]; !exists {
		return nil, nil, fmt.Errorf("Vertex %v does not exist.", source)
	}

	dist = make(map[int]float64)
	prev = make(map[int]int)
	heap := NewFibHeap()

	for id := range g.vertices {
		prev[id] = -1
		if id != source {
			dist[id] = math.Inf(1)
			heap.Insert(id, math.Inf(1))
		} else {
			dist[id] = 0
			heap.Insert(id, 0)
		}
	}

	for heap.Num() != 0 {
		min, _ := heap.ExtractMin()
		for to, edge := range g.egress[min.(int)] {
			if dist[min.(int)] + edge.weight < dist[to] {
				heap.DecreaseKey(to, dist[min.(int)]+edge.weight)
				prev[to] = min.(int)
				dist[to] = dist[min.(int)] + edge.weight
			}
		}
	}

	return
}

func (g *Graph) ShortestPaths(base int, db *sql.DB) {
	inf_count := 0
	max_correction := 0.0

	v := make(plotter.Values, len(g.vertices)-1)

	i := 0
	for id := range g.vertices {
		if id != base {
			dist, _, err := g.Dijkstra(id)
			if id == 1128 {
				fmt.Println("HERE")
				fmt.Println(dist[base])
			}
			check(err)
			if dist[base] == math.Inf(1) {
				inf_count++
			} else {
				if math.Abs(dist[base]) > math.Abs(max_correction) {max_correction = dist[base]}
				
				v[i] = dist[base]
				i++

				if id == 1128 {
					fmt.Println(dist[base])
					os.Exit(1)
				}

				// Update the race in the database
				UpdateRace(db, id, dist[base])


			}
			// os.Exit(1)
		}
	}

	p, err := plot.New()
	check(err)

	h, err := plotter.NewHist(v, int(math.Sqrt(float64(len(g.vertices)-1))))
	check(err)

	p.Add(h)
	save_file := fmt.Sprintf("hist%v.png", base)
	err = p.Save(8*vg.Inch, 8*vg.Inch, save_file)
	check(err)

	fmt.Printf("Valid Vertices: %0.4f %%\n", 100 * float64(inf_count)/float64(g.Length()))
	fmt.Printf("Max Correction: %0.4f\n", max_correction)
}