package evolution

import (
	"math/rand"
	"testing"
	"time"
	"tkOptimizer/internal/layout"
)

func TestRecombinationSingle(t *testing.T) {
  e, err := New(
    1,
    100000,
    1.0,
    0.001,
    NewKeyboardConfig(
      8, 8,
      [][]float64{}, layout.Layout{},
      []rune("abcdefgklmnoprsthqwvx"),
    ),
    "one,two,thre,four,five,six",
  )
  if err != nil {
    t.Error(err)
    return
  }
  genK, err := GenerateKeyboards(e.KeyboardConfig, e.GetInitPopulation())
  if err != nil {
    t.Error(err)
    return
  }
  genK, err = Recombine(genK, e.MutationProbability)
  if err != nil {
    t.Error(err)
    return
  }
}

func TestRecombinationThreads(t *testing.T) {
  e, err := New(
    7,
    1000000,
    1.0,
    0.001,
    NewKeyboardConfig(
      8, 8,
      [][]float64{}, layout.Layout{},
      []rune("abcdefgklmnoprsthqwvx"),
    ),
    "one,two,thre,four,five,six",
  )
  if err != nil {
    t.Error(err)
    return
  }
  kN, err := GenerateKeyboardsThreads(e.Threads, e.KeyboardConfig, e.GetInitPopulation())
  if err != nil {
    t.Error(err)
    return
  }
  kN, err = e.RecombineThreads(e.Threads, kN)
  if err != nil {
    t.Error(err)
    return
  }
}

func TestRunSingle(t *testing.T) {
  e, err := New(
    7,
    1000000,
    0.2,
    0.05,
    NewKeyboardConfig(
      8, 8,
      [][]float64{}, layout.Layout{},
      []rune("abcdefgklmnoprsthqwvx"),
    ),
    "one,two,thre,four,five,six",
  )
  if err != nil {
    t.Error(err)
    return
  }
  genK, err := GenerateKeyboards(e.KeyboardConfig, e.GetInitPopulation())
  genK, err = Run(e, genK)
  if err != nil {
    t.Error(err)
    return
  }
  t.Error()
}

func TestRunThreads(t *testing.T) {
  e, err := FromYaml("../../test/configEvo.yml")
  //e, err := New(
  //  10,
  //  1.0,
  //  0.05,
  //  NewKeyboardConfig(
  //    8, 8, [][]float64{},
  //    []rune("abcdefgklmnoprsthqwvx"),
  //  ),
  //  "one,two,thre,four,five,six",
  //)
  //if err != nil {
  //  t.Error(err)
  //  return
  //}
  k, err := GenerateKeyboardsThreads(e.Threads, e.KeyboardConfig, e.GetInitPopulation())

  k, err = Run(e, k)
  if err != nil {
    t.Error(err)
  }
  TestKeyboards(k, e.TestText)
  t.Log("Best keyboard: ", k[0].Distance)
  t.Error()
}

func TestReproducableSeed(t *testing.T) {
  timeNow := time.Now()
  t.Log(timeNow.UnixNano())
  r := rand.New(rand.NewSource(1721169499718084584))
  for i := 0; i < 5; i++ {
    t.Log(r.Float64())
  }
  t.Error()
}
