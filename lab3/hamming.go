package main

import (
	"fmt"
	"hash"
	"log"
	"math"
	"math/bits"
	"runtime"
	"sync"

	"gonum.org/v1/gonum/stat"
)

const (
	ExpectedValueHamming  = 128.0 // n * p = 256 * 0.5
	ExpectedStdDevHamming = 8.0   // sqrt(n * p * (1-p)) = sqrt(256 * 0.5 * 0.5) = 8
)

func hammingTest(newHash func() hash.Hash) int {
	distances := getHammingDistancesForHash(newHash)
	mean := stat.Mean(distances, nil)

	zValue := (mean - ExpectedValueHamming) * math.Sqrt(Size) / ExpectedStdDevHamming

	if math.Abs(zValue) <= CriticalValue {
		return 1
	} else {
		return 0
	}
}

func getHammingDistancesForHash(newHash func() hash.Hash) []float64 {
	distances := make([]float64, Size)
	jobs := make(chan int, Size)

	var wg sync.WaitGroup
	numWorkers := runtime.NumCPU()

	for range numWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for i := range jobs {
				h := newHash()
				data1 := generateRandomData(64)
				data2 := changeRandomBitInData(data1)

				h.Write(data1)
				hash1 := h.Sum(nil)
				h.Reset()
				h.Write(data2)
				hash2 := h.Sum(nil)

				distance, err := calculateHammingDistance(hash1, hash2)
				if err != nil {
					log.Printf("error calculating hamming distance for index %d: %v", i, err)
					distances[i] = 0
					continue
				}
				distances[i] = float64(distance)
			}
		}()
	}

	for i := range Size {
		jobs <- i
	}
	close(jobs)

	wg.Wait()
	return distances
}

func calculateHammingDistance(hash1, hash2 []byte) (int, error) {
	if len(hash1) != len(hash2) {
		return 0, fmt.Errorf("hashes must have the same length")
	}

	distance := 0
	for i := range len(hash1) {
		xorRes := hash1[i] ^ hash2[i]
		distance += bits.OnesCount8(xorRes)
	}

	return distance, nil
}
