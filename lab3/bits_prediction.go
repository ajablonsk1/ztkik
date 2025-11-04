package main

import (
	"hash"
	"math"
	"runtime"
	"sync"
)

const (
	ExpectedValueBitPrediction  = 0.5 // p_0
	ExpectedStdDevBitPrediction = 0.5 // sqrt(p_0 * (1-p_0))
)

func bitsPredictionTest(newHash func() hash.Hash) int {
	probabilities := getBitProbabilitiesForHash(newHash)

	passed := 0
	for _, probability := range probabilities {
		zValue := (probability - ExpectedValueBitPrediction) * math.Sqrt(Size) / ExpectedStdDevBitPrediction
		if math.Abs(zValue) <= CriticalValue {
			passed++
		}
	}
	return passed
}

func getBitProbabilitiesForHash(newHash func() hash.Hash) []float64 {
	type result struct {
		localCounts [256]uint64
	}

	results := make(chan result, runtime.NumCPU())
	jobs := make(chan int, Size)

	var wg sync.WaitGroup
	numWorkers := runtime.NumCPU()

	for range numWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var localCounts [256]uint64

			for range jobs {
				h := newHash()
				data := generateRandomData(64)
				h.Write(data)
				hash := h.Sum(nil)

				for i, hashByte := range hash {
					for j := range 8 {
						if (hashByte & (1 << j)) != 0 {
							localCounts[8*i+j]++
						}
					}
				}
			}

			results <- result{localCounts: localCounts}
		}()
	}

	go func() {
		for i := range Size {
			jobs <- i
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	bitCounts := make([]uint64, 256)
	for res := range results {
		for i := range 256 {
			bitCounts[i] += res.localCounts[i]
		}
	}

	probabilities := make([]float64, 256)
	for k := range 256 {
		probabilities[k] = float64(bitCounts[k]) / float64(Size)
	}

	return probabilities
}
