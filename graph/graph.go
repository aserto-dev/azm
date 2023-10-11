package graph

type Graph struct {
	nodes     []string
	adjMatrix [][]string
}

func NewGraph() *Graph {
	return &Graph{
		nodes:     make([]string, 0),
		adjMatrix: make([][]string, 0),
	}
}

func (g *Graph) AddNode(nodeName string) {
	for _, vertexName := range g.nodes {
		if nodeName == vertexName {
			return
		}
	}
	g.nodes = append(g.nodes, nodeName)
	g.adjMatrix = append(g.adjMatrix, make([]string, 0))
}

func (g *Graph) AddEdge(source, dest, edgeName string) {
	sourceIndex := g.findVertexIndexByName(source)
	destIndex := g.findVertexIndexByName(dest)
	if destIndex == -1 || sourceIndex == -1 {
		return
	}

	if len(g.adjMatrix[sourceIndex]) == 0 {
		g.adjMatrix[sourceIndex] = make([]string, len(g.nodes))
	}

	g.adjMatrix[sourceIndex][destIndex] = edgeName
}

func (g *Graph) dFS(source int, visited []bool, traversal []string) []string {
	visited[source] = true
	traversal = append(traversal, g.nodes[source])
	for i := 0; i < len(g.nodes); i++ {
		if len(g.adjMatrix[source]) != 0 && g.adjMatrix[source][i] != "" && !visited[i] {
			traversal = append(traversal, g.adjMatrix[source][i])
			traversal = g.dFS(i, visited, traversal)
		}
	}
	return traversal
}

func (g *Graph) TraverseGraph(startVertex string) []string {
	visited := make([]bool, len(g.nodes))
	traversal := make([]string, 0)
	startIndex := g.findVertexIndexByName(startVertex)
	if startIndex == -1 {
		return traversal
	}

	return g.dFS(startIndex, visited, traversal)
}

func (g *Graph) findVertexIndexByName(name string) int {
	for index, vertexName := range g.nodes {
		if vertexName == name {
			return index
		}
	}
	return -1
}
