TextRank in Go
==============

This is an implementation of the TextRank algorithm in Go.  TextRank is a graph based method for ranking text e.g. 
finding the most important sentences or words within a document.  I wrote this in order to learn Go. 

More information on TextRank can be found in this paper : http://acl.ldc.upenn.edu/acl2004/emnlp/pdf/Mihalcea.pdf

The Package
===========

The package consists of two parts:

* a simple graph data structure
* the TextRank Algorithm 

graph.go
--------

This is where we implement our graph data structure. 

```go
//vertices are integers
type Vertex int

//edges are always directed and weighted 
//edges are stored in lists indexed by their starting point; so we only store where the edge points to
type Edge struct {
  to Vertex
  weight float64
}

//we store two maps - one that maps inbound edges to each vertex and one outbound edges from each vertex
//we use this approach in order to have fast edges in and edges out queries for each vertex 
type Graph struct {
  vertexCount Vertex
  maxVertex int
  OutBoundEdges map[Vertex][]Edge //Out(V) queries
  InBoundEdges map[Vertex][]Edge  //In(V) queries
}
```

For textrank, i wanted fast edges in and edges out of each node lookups without using a matrix representation of the graph.  
I will be refining this to work on far larger graphs, so matrix representations are out. 

Graphs also have a few helper functions:

```go
func (g *Graph) AddVertex() (Vertex)

func (g *Graph) VertexCount() (int)

func (g *Graph) AddEdge(from Vertex, to Vertex, weight float64)

func (g *Graph) In(v Vertex) ([]Edge)

func (g *Graph) Out(v Vertex) ([]Edge)

func (g *Graph) Weight(from Vertex, to Vertex) (weight float64, e error)
```


textrank.go
-----------

The most important function is Iterate, the actual implementation of TextRank:

```go
func Iterate(d float64, oldScores []float64, g *Graph) ([]float64)
```

Iterate takes a param d (d=.85 in the paper) along with a vector of vertex scores and graph.  It returns an updated vector of scores based on the graph using the TextRank algorithm.


The following are a few helper functions:

```go
//this is useful for checking for convergence
func ScoreDiff(s1 []float64, s2 []float64) (float64)

//pair indices and score for sorting
type IndexScorePair struct {
  Index int
  Score float64
}

func Sort(pairs []IndexScorePair)

```

Example
=======

The following example uses TextRank on the sentences of a state of the union address to determine the most 'important' sentences.


```go
package main

//rough implementation of: http://acl.ldc.upenn.edu/acl2004/emnlp/pdf/Mihalcea.pdf
//this examples uses a state of the union address and creates a graph of sentences

import (
  "github.com/cioc/textrank"
  "fmt"
  "io/ioutil"
  "strings"
  "math"
  "runtime"
)

//checks if w is in sentence sentence
func wordInSentence(w string, sentence []string) (bool) {
  for _, w2 := range sentence {
    if w == w2 {
      return true
    }
  }
  return false
}

//similarity score from paper 
func sentenceSimilarity(s1 string, s2 string) (float64) {
  w1 := strings.Split(strings.ToLower(s1), " ")
  w2 := strings.Split(strings.ToLower(s2), " ")
  c := 0
  for _, w := range w1 {
    if wordInSentence(w, w2) {
      c += 1
    }
  }
  o := float64(c) / (math.Log(float64(len(w1))) + math.Log(float64(len(w2))))
  if math.IsNaN(o) {
    return float64(0)
  }
  if math.IsInf(o, 0) {
    return float64(0)
  }
  return o
}

func main() {
  //you should adjust this based on your machine
  runtime.GOMAXPROCS(8)
  content, err := ioutil.ReadFile("Barack_Obama_2013.txt")
  if err != nil {
    panic(err)
  }
  lines := strings.Split(string(content), ".")
  g := textrank.NewGraph(len(lines))

  //create a mapping from graph id to lines for output
  omap := make(map[textrank.Vertex]string)
  for _, l := range lines {
    omap[g.AddVertex()] = l
  }

  //build similarity graph between sentences
  for i, l1 := range lines {
    for k, l2 := range lines {
      g.AddEdge(textrank.Vertex(i), textrank.Vertex(k), sentenceSimilarity(l1, l2))
    }
  }

  //run textrank until convergence
  oldScores := make([]float64, len(lines), len(lines))
  c := float64(10000)
  for ; c > float64(0.0001); {
    newScores := textrank.Iterate(float64(.85), oldScores, g)
    c = textrank.ScoreDiff(oldScores, newScores)
    fmt.Printf("Cumulative diff: %v\n", c)
    for i, v := range newScores {
      oldScores[i] = v
    }
  }

  //sort scores and print out top 10 sentences
  fmt.Printf("%v\n", oldScores)
  pairs := make([]textrank.IndexScorePair, len(lines), len(lines))
  for i, v := range oldScores {
    pairs[i] = textrank.IndexScorePair{i, v}
  }
  textrank.Sort(pairs)
  for i := 0; i < 10; i++ {
    fmt.Printf("%v : %v, : %v\n", pairs[i].Index, pairs[i].Score, omap[textrank.Vertex(pairs[i].Index)])
  }
}
```
