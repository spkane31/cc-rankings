package rankings

import (
	"fmt"
	"math"
)

type ID interface{}

type Vertex interface {

	ID() ID

	Edges() []Edge
}

type Edge interface {
	Get() (from ID, to ID, weight float64)
}

type Graph struct {
	vertices 	map[ID]*vertex
	egress 		map[ID]map[ID]*edge
	ingress		map[ID]map[ID]*edge
}

type vertex struct {
	self interface{}
	enable bool
}

type edge struct {
	self interface{}
	weight float64
	enable bool
	changed bool
}

func (e *edge) GetWeight() float64 {
	return e.weight
}

func NewGraph() *Graph {
	g := new(Graph)
	g.vertices = make(map[ID]*vertex)
	g.egress = make(map[ID]map[ID]*edge)
	g.ingress = make(map[ID]map[ID]*edge)

	return g
}

func (g *Graph) GetVertex(id ID) (vertex interface{}, err error) {
	if v, exists := g.vertices[id]; exists {
		vertex = v.self
		return
	}

	err = fmt.Errorf("Vertex %v is not found", id)
	return
}

func (g *Graph) GetEdge(from, to ID) (interface{}, error) {
	if _, exists := g.vertices[from]; !exists {
		return nil, fmt.Errorf("Vertex(from) %v is not found", from)
	}

	if _, exists := g.vertices[to]; !exists {
		return nil, fmt.Errorf("Vertex(to) %v is not found", to)
	}

	if edge, exists := g.egress[from][to]; exists {
		return edge.self, nil
	}

	return nil, fmt.Errorf("Edge from %v to %v is not found", from, to)
}

func (g *Graph) GetEdgeWeight(from, to ID) (float64, error) {

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

func (g *Graph) AddVertex(id ID, v interface{}) error {
	if _, exists := g.vertices[id]; exists {
		return fmt.Errorf("Vertex %v already exists", id)
	}

	g.vertices[id] = &vertex{v, true}
	g.egress[id] = make(map[ID]*edge)
	g.ingress[id] = make(map[ID]*edge)

	return nil
}

func (g *Graph) AddEdge(from, to ID, weight float64, e interface{}) error {
	if weight == math.Inf(1) {

		return fmt.Errorf("-inf weight is reserved for internal usage")
	}

	if _, exists := g.vertices[from]; !exists {
		return fmt.Errorf("Vertex(from) %v is not found", from)
	}

	if _, exists := g.vertices[to]; !exists {
		return fmt.Errorf("Vertex(to) %v is not found", to)
	}

	if _, exists := g.egress[from][to]; exists {
		return fmt.Errorf("Edge from %v to %v is duplicate", from, to)
	}

	g.egress[from][to] = &edge{e, weight, true, false}
	g.ingress[to][from] = g.egress[from][to]

	return nil
}

func (g *Graph) UpdateEdgeWeight(from, to ID, weight float64) error {
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

func (g *Graph) DeleteVertex(id ID) interface{} {
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

		return vertex.self
	}

	return nil
}

func (g *Graph) DeleteEdge(from, to ID) interface{} {
	if _, exists := g.vertices[from]; !exists {
		return nil
	}

	if _, exists := g.vertices[to]; !exists {
		return nil
	}

	if edge, exists := g.egress[from][to]; exists {
		delete(g.egress[from], to)
		delete(g.ingress[to], from)
		return edge.self
	}

	return nil
}

func (g *Graph) AddVertexWithEdges(v Vertex) error {
	if _, exists := g.vertices[v.ID()]; exists {
		return fmt.Errorf("Vertex %v is duplicate", v.ID())
	}

	g.vertices[v.ID()] = &vertex{v, true}
	g.egress[v.ID()] = make(map[ID]*edge)
	g.ingress[v.ID()] = make(map[ID]*edge)

	for _, eachEdge := range v.Edges() {
		from, to, weight := eachEdge.Get()

		if weight == math.Inf(-1) {
			return fmt.Errorf("-inf weight is reserved for internal usage")
		}
		if from != v.ID() && to != v.ID() {
			return fmt.Errorf("Edge from %v to %v is unrelated to vertex %v", from, to, v.ID())
		}

		if _, exists := g.egress[to]; !exists {
			g.egress[to] = make(map[ID]*edge)
		}

		if _, exists := g.ingress[from]; !exists {
			g.egress[from] = make(map[ID]*edge)
		}

		if _, exists := g.ingress[from]; !exists {
			g.ingress[from] = make(map[ID]*edge)
		}

		if _, exists := g.ingress[to]; !exists {
			g.ingress[to] = make(map[ID]*edge)
		}

		g.egress[from][to] = &edge{eachEdge, weight, true, false}
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

func (g *Graph) GetPathWeight(path []ID) (totalWeight float64) {
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

func (g *Graph) DisableEdge(from, to ID) {
	g.egress[from][to].enable = false
}

func (g *Graph) DisableVertex(vertex ID) {
	for _, edge := range g.egress[vertex] {
		edge.enable = false
	}
}

func (g *Graph) DisablePath(path []ID) {
	for _, vertex := range path {
		g.DisableVertex(vertex)
	}
}

func (g *Graph) Reset() {
	for _, out := range g.egress {
		for _, edge := range out {
			edge.enable = true
		}
	}
}

// func BuildGraph() Graph {

// }