package textrank

type Vertex int

type Edge struct {
  to Vertex
  weight float64
}

type Graph struct {
  vertexCount Vertex
  maxVertex int
  OutBoundEdges map[Vertex][]Edge //Out(V) queries
  InBoundEdges map[Vertex][]Edge  //In(V) queries
}

func NewGraph(maxVertex int) (*Graph) {
  return &Graph{0, maxVertex, make(map[Vertex][]Edge), make(map[Vertex][]Edge)}
}

func (g *Graph) AddVertex() (Vertex) {
  o := g.vertexCount
  g.vertexCount += 1
  return o
}

func (g *Graph) VertexCount() (int) {
  return int(g.vertexCount)
}

func (g *Graph) AddEdge(from Vertex, to Vertex, weight float64) {
  if g.OutBoundEdges[from] == nil {
    g.OutBoundEdges[from] = make([]Edge, 1, g.maxVertex + 1)
    g.OutBoundEdges[from][0] = Edge{to, weight}
  } else {
    g.OutBoundEdges[from] = append(g.OutBoundEdges[from], Edge{to, weight})
  }
  if g.InBoundEdges[to] == nil {
    g.InBoundEdges[to] = make([]Edge, 1, g.maxVertex + 1)
    g.InBoundEdges[to][0] = Edge{from, weight}
  } else {
    g.InBoundEdges[to] = append(g.InBoundEdges[to], Edge{from, weight})
  }
}

func (g *Graph) In(v Vertex) ([]Edge) {
  return g.InBoundEdges[v]
}

func (g *Graph) Out(v Vertex) ([]Edge) {
  return g.OutBoundEdges[v]
}

func (g *Graph) Weight(from Vertex, to Vertex) (weight float64, e error) {
  for _, edge := range g.OutBoundEdges[from] {
    if edge.to == to {
      return edge.weight, nil
    }
  }
  return float64(0), e
}
