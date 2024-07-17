package evolution

import (
	"fmt"
	"math"
	"tkOptimizer/internal/keyboard"
)

func GetLenOr1000(num int) int {
  if num < 1000 {
    return num
  }
  return 1000
}

func (e *Evolution) GenerateKeyboardsThreads(population int) ([]*keyboard.Keyboard, error) {
  numWorkers := e.Threads
	jobs := make(chan int, population)
	results := make(chan *keyboard.Keyboard, GetLenOr1000(population))
	errorChan := make(chan error, 10)

  weights := e.KeyboardConfig.Weights
  height := e.KeyboardConfig.Height
  width := e.KeyboardConfig.Width
  charSet := e.KeyboardConfig.CharSet

	for w := 1; w <= numWorkers; w++ {
		go func() {
      for range jobs {
        var k *keyboard.Keyboard
        var err error

        if len(weights) != 0 {
          k, err = keyboard.GenerateNewWithWeights(
            height, width, weights, charSet)
          if err != nil {
            errorChan <- err
            return
          }
        } else {
          k, err = keyboard.GenerateNew(height, width, charSet)
          if err != nil {
            errorChan <- err
            return
          }
        }
        results <- k
      }
    return
		}()
	}

	for i := 0; i < population; i++ {
		jobs <- 1
	}
	close(jobs)

  newKeyboards := make([]*keyboard.Keyboard, population)
	for j := 0; j < population; j++ {
    select {
    case r := <-results:
      newKeyboards[j] = r
    case err := <-errorChan:
      close(results)
      close(errorChan)
      return nil, err
    }
	}

  close(results)
  close(errorChan)

	return newKeyboards, nil
}

func (e *Evolution) TestKeyboardsThreads(k []*keyboard.Keyboard) {
  numWorkers := e.Threads
	jobs := make(chan *keyboard.Keyboard, len(k))
	results := make(chan int, GetLenOr1000(len(k)))

	for w := 1; w <= numWorkers; w++ {
		go func() {
      for k := range jobs {
        k.TravelDistance(e.TestText)
        results <- 1
      }
    }()
	}
	for i := 0; i < len(k); i++ {
		jobs <- k[i]
	}
	close(jobs)

	for j := 0; j < len(k); j++ {
		<-results
	}
}

func recombineWorkerThreads(
	resultChan chan<- *keyboard.Keyboard,
	errorChan chan<- error,
	mutationProbability float64,
	k1 *keyboard.Keyboard,
	k2 *keyboard.Keyboard,
) {
	mergeKeyboard, err := Recombination(mutationProbability, k1, k2)
	if err != nil {
		errorChan <- err
		return
	}
	resultChan <- mergeKeyboard
}

func (e *Evolution) RecombineThreads(k []*keyboard.Keyboard) ([]*keyboard.Keyboard, error) {
  numWorkers := e.Threads
	passNumber := math.Floor(e.Percentile * float64(len(e.Keyboards)))
	passNumberInt := int(passNumber)
	if !IsEven(passNumberInt) {
		passNumberInt -= 1
	}
	if passNumberInt < 2 {
    fmt.Println("Keyboards len < 2")
		return nil, nil
	}

	semaphore := make(chan int, numWorkers)
	keyboardChan := make(chan *keyboard.Keyboard, passNumberInt/2)
	errorChan := make(chan error, 5)

	// Launch worker goroutines
	for i := 0; i < passNumberInt; i += 2 {
		semaphore <- 1 // Acquire a semaphore slot

		go func(k1, k2 *keyboard.Keyboard) {
			defer func() {
				<-semaphore // Release the semaphore slot
			}()

			mK, err := Recombination(e.MutationProbability, k1, k2)
			if err != nil {
				errorChan <- err
				return
			}
			keyboardChan <- mK
		}(e.Keyboards[i], e.Keyboards[i+1])
	}

	// Process keyboards as they become available
	nextGen := make([]*keyboard.Keyboard, passNumberInt/2)
	for i := 0; i < len(nextGen); i++ {
		select {
		case k := <-keyboardChan:
			nextGen[i] = k
		case err := <-errorChan:
			close(keyboardChan)
			close(errorChan)
			return nil, err
		}
	}

	// Close channels after all processing is done
	close(keyboardChan)
	close(errorChan)

	return nextGen, nil
}

func (e *Evolution) RunThreads() error {
  nK, err := e.GenerateKeyboardsThreads(
    PopulationSizeNext(e.Percentile, e.Population))
  if err != nil {
    return err
  }
	for math.Floor(e.Percentile * float64(len(e.Keyboards))) > 2 {
		e.TestKeyboardsThreads(nK)
    e.SortKeyboards(nK)
    e.AppendMetric(nK[0].Distance)
    e.SetKeyboards(nK)
    nK, err = e.RecombineThreads(nK)
    if err != nil {
      return err
    }
	}
  return nil
}

