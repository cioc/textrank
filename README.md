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
