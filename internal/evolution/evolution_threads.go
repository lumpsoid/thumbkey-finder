package evolution

import (
	"fmt"
	"tkOptimizer/internal/keyboard"
)

func GetLenOr1000(num int) int {
  if num < 1000 {
    return num
  }
  return 1000
}

func GenerateKeyboardsThreads(
  threads int, 
  config *KeyboardConfiguration, 
  population int,
  placeThreshold float64,
) ([]*keyboard.Keyboard, error) {
	jobs := make(chan int, population)
	results := make(chan *keyboard.Keyboard, GetLenOr1000(population))
	errorChan := make(chan error, 10)

  weights := config.Weights
  height := config.Height
  width := config.Width
  charSet := config.CharSet

	for w := 1; w <= threads; w++ {
		go func() {
      for range jobs {
        var k *keyboard.Keyboard
        var err error

        if len(weights) != 0 {
          k, err = keyboard.GenerateNewWithWeights(
            height, width, weights, charSet, placeThreshold)
          if err != nil {
            errorChan <- err
            return
          }
        } else {
          k, err = keyboard.GenerateNew(height, width, charSet, placeThreshold)
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

func TestKeyboardsThreads(threads int, k []*keyboard.Keyboard, testText string) {
	jobs := make(chan *keyboard.Keyboard, len(k))
	results := make(chan int, GetLenOr1000(len(k)))

	for w := 1; w <= threads; w++ {
		go func() {
      for k := range jobs {
        k.TravelDistance(testText)
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
  placeThreshold float64,
	k1 *keyboard.Keyboard,
	k2 *keyboard.Keyboard,
) {
	mergeKeyboard, err := Recombination(mutationProbability, placeThreshold, k1, k2)
	if err != nil {
		errorChan <- err
		return
	}
	resultChan <- mergeKeyboard
}

func RecombineThreads(
  threads int, 
  mutationProbability float64,
  placeThreshold float64,
  k []*keyboard.Keyboard,
) ([]*keyboard.Keyboard, error) {
  kLen := len(k)
	if !IsEven(kLen) {
		kLen -= 1
	}
	if kLen < 2 {
    fmt.Println("Keyboards len < 2")
		return k, nil
	}

	semaphore := make(chan int, threads)
	keyboardChan := make(chan *keyboard.Keyboard, kLen/2)
	errorChan := make(chan error, 5)

	// Launch worker goroutines
	for i := 0; i < kLen; i += 2 {
		semaphore <- 1 // Acquire a semaphore slot

		go func(k1, k2 *keyboard.Keyboard) {
			defer func() {
				<-semaphore // Release the semaphore slot
			}()

			mK, err := Recombination(mutationProbability, placeThreshold, k1, k2)
			if err != nil {
				errorChan <- err
				return
			}
			keyboardChan <- mK
		}(k[i], k[i+1])
	}

	// Process keyboards as they become available
	nextGen := make([]*keyboard.Keyboard, kLen/2)
	for j := 0; j < len(nextGen); j++ {
		select {
		case k := <-keyboardChan:
			nextGen[j] = k
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

