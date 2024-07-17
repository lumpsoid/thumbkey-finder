package evolution

import (
	"math/rand"
	"testing"
	"time"
)

func TestRecombinationSingle(t *testing.T) {
  e, err := New(
    1,
    100000,
    1.0,
    0.001,
    NewKeyboardConfig(
      8, 8, [][]float64{},
      []rune("abcdefgklmnoprsthqwvx"),
    ),
    "one,two,thre,four,five,six",
  )
  if err != nil {
    t.Error(err)
    return
  }
  err = e.GenerateKeyboards(e.Population)
  if err != nil {
    t.Error(err)
    return
  }
  err = e.Recombine()
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
      8, 8, [][]float64{},
      []rune("abcdefgklmnoprsthqwvx"),
    ),
    "one,two,thre,four,five,six",
  )
  if err != nil {
    t.Error(err)
    return
  }
  kN, err := e.GenerateKeyboardsThreads(e.Population)
  if err != nil {
    t.Error(err)
    return
  }
  kN, err = e.RecombineThreads(kN)
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
      8, 8, [][]float64{},
      []rune("abcdefgklmnoprsthqwvx"),
    ),
    "one,two,thre,four,five,six",
  )
  if err != nil {
    t.Error(err)
    return
  }
  err = e.Run()
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

  err = e.RunThreads()
  if err != nil {
    t.Error(err)
  }
  e.TestKeyboards()
  t.Log("Best keyboard: ", e.Keyboards[0].Distance)
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
