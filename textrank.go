package textrank

import (
  "math"
  "sort"
)

func Iterate(d float64, oldScores []float64, g *Graph) ([]float64) {
  ci := make(chan int)
  newScores := make([]float64, len(oldScores), len(oldScores))
  for i := range oldScores {
    go func(j int) {
      s := float64(0)
      inBound := g.In(Vertex(j))
      for vin, vinVertex := range inBound {
        outBound := g.Out(Vertex(vin))
        denom := float64(0)
        for _, out := range outBound {
          denom += out.weight
        }
        w, _ := g.Weight(vinVertex.to, Vertex(j))
        if denom != 0 && !math.IsInf(denom, 0) {
          s += (w / denom) * oldScores[vinVertex.to]
        }
      }
      newScores[j] = (1 - d) + (d * s)
      ci <- 1
    }(i)
  }
  for j := 0 ; j < len(oldScores); j++ {
    <-ci
  }
  close(ci)
  return newScores
}

func ScoreDiff(s1 []float64, s2 []float64) (float64) {
  cumulative := float64(0)
  for i, v1 := range s1 {
    cumulative += (v1 - s2[i]) * (v1 - s2[i])
  }
  return cumulative
}

type IndexScorePair struct {
  Index int
  Score float64
}

type indexScorePairSorter struct {
  pairs []IndexScorePair
  by func(p1, p2 *IndexScorePair) bool
}

func (s *indexScorePairSorter) Len() int {
  return len(s.pairs)
}

func (s *indexScorePairSorter) Swap(i, j int) {
  s.pairs[i], s.pairs[j] = s.pairs[j], s.pairs[i]
}

func (s *indexScorePairSorter) Less(i, j int) bool {
  return s.by(&s.pairs[i], &s.pairs[j])
}

func Sort(pairs []IndexScorePair) {
  sorter := &indexScorePairSorter{
    pairs: pairs,
    by: func(p1, p2 *IndexScorePair) bool {
      return p1.Score > p2.Score
    },
  }
  sort.Sort(sorter)
}
