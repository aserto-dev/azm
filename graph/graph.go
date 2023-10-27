package graph

type Relations []string

type Graph struct {
	nodes     []string
	adjMatrix [][]Relations
}

func NewGraph() *Graph {
	return &Graph{
		nodes:     make([]string, 0),
		adjMatrix: make([][]Relations, 0),
	}
}

func (g *Graph) AddNode(nodeName string) {
	for _, vertexName := range g.nodes {
		if nodeName == vertexName {
			return
		}
	}
	g.nodes = append(g.nodes, nodeName)
	g.adjMatrix = append(g.adjMatrix, make([]Relations, 0))
}

func (g *Graph) AddEdge(source, dest, edgeName string) {
	sourceIndex := g.findVertexIndexByName(source)
	destIndex := g.findVertexIndexByName(dest)
	if destIndex == -1 || sourceIndex == -1 {
		return
	}

	if len(g.adjMatrix[sourceIndex]) == 0 {
		g.adjMatrix[sourceIndex] = make([]Relations, len(g.nodes))
	}

	g.adjMatrix[sourceIndex][destIndex] = append(g.adjMatrix[sourceIndex][destIndex], edgeName)
}

func (g *Graph) dFS(source, dest int, visited []map[int]bool, traversal [][]string, index int) [][]string {
	if len(traversal) < index+1 {
		traversal = append(traversal, make([]string, 0))
	}
	traversal[index] = append(traversal[index], g.nodes[source])

	if source == dest {
		return traversal
	} else {
		for i := 0; i < len(g.nodes); i++ {
			if len(g.adjMatrix[source]) > i && len(g.adjMatrix[source][i]) != 0 {
				if len(visited[i]) == 0 {
					visited[i] = make(map[int]bool, len(g.adjMatrix[source][i]))
				}
				for j := 0; j < len(g.adjMatrix[source][i]); j++ {
					if visited[i][j] {
						continue
					}
					visited[i][j] = true
					if len(traversal) < index+1 {
						traversal = append(traversal, make([]string, 0))
					}
					traversal[index] = append(traversal[index], g.adjMatrix[source][i][j])
					traversal = g.dFS(i, dest, visited, traversal, index)
					index++
				}
				visited[i] = make(map[int]bool, 0)
			}
		}
	}

	return traversal
}

func (g *Graph) FindPaths(startVertex, destVertx string) [][]string {
	visited := make([]map[int]bool, len(g.nodes))
	traversal := make([][]string, 0)
	startIndex := g.findVertexIndexByName(startVertex)
	destIndex := g.findVertexIndexByName(destVertx)
	if startIndex == -1 || destIndex == -1 {
		return traversal
	}

	traversal = g.dFS(startIndex, destIndex, visited, traversal, 0)
	for i, paths := range traversal {
		if paths[0] != startVertex {
			traversal[i] = append([]string{startVertex}, paths...)
		}

	}

	return traversal
}

func (g *Graph) findVertexIndexByName(name string) int {
	for index, vertexName := range g.nodes {
		if vertexName == name {
			return index
		}
	}
	return -1
}
