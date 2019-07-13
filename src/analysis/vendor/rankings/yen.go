package rankings

import (
	"fmt"
	"math"
	"sort"
)

func getPath(prev map[int]int, lastNode int) (path []int) {
	prevNode, has := prev[lastNode]
	if has == false {
		return nil
	}

	reversePath := []int{lastNode}
	for ; has != false; prevNode, has = prev[prevNode] {
		reversePath = append(reversePath, prevNode)
	}

	path = make([]int, len(reversePath))
	for idx, n := range reversePath {
		path[len(reversePath) - idx-1] = n
	}

	return
}

type potential struct {
	dist float64
	path []int
}

func (g *Graph) Yen(source, destination, K int) ([]float64, [][]int, error) {
	var err error
	var i, j, k int
	var dijkstraDist map[int]float64
	var dijkstraPrev map[int]int
	var existed bool
	var spurWeight float64
	var spurPath []int
	var potentials []potential
	distTopK := make([]float64, K)
	pathTopK := make([][]int, K)

	for i := 0; i < K; i++ {
		distTopK[i] = math.Inf(1)
	}

	dijkstraDist, dijkstraPrev, err = g.Dijkstra(source)
	if err != nil {
		return nil, nil, err
	}
	distTopK[0] = dijkstraDist[destination]
	pathTopK[0] = getPath(dijkstraPrev, destination)

	fmt.Println(distTopK, pathTopK)

	for k = 1; k < K; {
		for i = 0; i < len(pathTopK[k-1])-1; i++ {
			for j = 0; j < k; j++ {
				if isShareRootPath(pathTopK[j], pathTopK[k-1][:i+1]) {
					g.DisableEdge(pathTopK[j][i], pathTopK[j][i+1])
				}
			}
			g.DisablePath(pathTopK[k-1][:i])

			dijkstraDist, dijkstraPrev, _ = g.Dijkstra(pathTopK[k-1][i])
			if dijkstraDist[destination] != math.Inf(1) {
				spurWeight = g.GetPathWeight(pathTopK[k-1][:i+1]) + dijkstraDist[destination]
				spurPath = mergePath(pathTopK[k-1][:i], getPath(dijkstraPrev, destination))
				existed = false

				for _, each := range potentials {
					if isSamePath(each.path, spurPath) {
						existed = true
						break
					}
				}

				if !existed {
					potentials = append(potentials, potential{spurWeight, spurPath})
				}
			}

			g.Reset()
		}

		if len(potentials) == 0 {
			break
		}
		sort.Slice(potentials, func(i, j int) bool {
			return potentials[i].dist < potentials[j].dist
		})

		if len(potentials) >= K - k {
			for l := 0; k < K; l++ {
				distTopK[k] = potentials[l].dist
				pathTopK[k] = potentials[l].path
				k++
			}
			break
		} else {
			distTopK[k] = potentials[0].dist
			pathTopK[k] = potentials[0].path
			potentials = potentials[1:]
			k++
		}
	}


	return distTopK, pathTopK, nil
}

func isShareRootPath(path, rootPath []int) bool {
	if len(path) < len(rootPath) { return false }

	return isSamePath(path[:len(rootPath)], rootPath)
}

func isSamePath(path1, path2 []int) bool {
	if len(path1) != len(path2) {
		return false
	}

	for i := 0; i < len(path1); i++ {
		if path1[i] != path2[i] {
			return false
		}
	}

	return true
}

func mergePath(path1, path2 []int) []int {
	newPath := []int{}
	newPath = append(newPath, path1...)
	newPath = append(newPath, path2...)

	return newPath
}

func (g *Graph) DisableEdge(from, to int) {
	if _, has := g.egress[from][to]; has {
		g.egress[from][to].enable = false
	}
}

func (g *Graph) DisableVertex(v int) {
	for _, edge := range g.egress[v] {
		edge.enable = false
	}
}

func (g *Graph) DisablePath(path []int) {
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

func (g *Graph) ResetVertices() {
	for _, vertex := range g.vertices {
		vertex.enable = true
	}
}