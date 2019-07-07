package rankings

import (
	"fmt"
	"math"
	"os"
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

	for v := 0; v < g.Length(); v++ {
		if sptSet[v] == false && math.Abs(dist[v]) <= min {
			min = dist[v]
			min_index = v
		}
	}

	return min_index
}


// Finds the shortest distances between two points 
// Since edges are seen as both positive and negative, the minimum between
// two edges is the values closest to zero (absolute value)
func (g *Graph) Dijkstra(start, target int) (cost float64) {
	
	dist := make(map[int]float64)
	// dist := make([]float64, g.Length())

	sptSet := make(map[int]bool)
	// sptSet := make([]bool, g.Length())

	for i := range dist {
		dist[i] = math.Inf(1)
		sptSet[i] = false
	}

	dist[start] = 0

	for count := 0; count < g.Length() - 1; count++ {
		u := g.minDistance(dist, sptSet)

		sptSet[u] = true
		
		for v := 0; v < g.Length(); v++ {
			_, has := g.egress[u][v]
			if !sptSet[v] && has && dist[u] != math.Inf(1) && dist[u] + g.egress[u][v].weight < dist[v] {
				dist[v] = dist[u] + g.egress[u][v].weight
			}
		}

		fmt.Println(dist[u])
		return
	}

	return
}

func (g *Graph) ShortestPaths(base int) {
	for id := range g.vertices {
		if id != base {
			fmt.Println(id, base)
			g.PrintVertex(id)
			g.Dijkstra(id, base)
			os.Exit(1)
		}
	}
}


// func (g *Graph) DisableEdge(from, to int) {
// 	g.egress[from][to].enable = false
// }

// func (g *Graph) DisableVertex(vertex int) {
// 	for _, edge := range g.egress[vertex] {
// 		edge.enable = false
// 	}
// }

// func (g *Graph) DisablePath(path []int) {
// 	for _, vertex := range path {
// 		g.DisableVertex(vertex)
// 	}
// }

// func (g *Graph) Reset() {
// 	for _, out := range g.egress {
// 		for _, edge := range out {
// 			edge.enable = true
// 		}
// 	}
// }

// func BuildGraph() Graph {

// }