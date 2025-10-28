package main

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha3"
	"fmt"
	"hash"
	"log"
	"math"
	"math/bits"
	mathrand "math/rand/v2"
	"slices"

	"github.com/magical/go-ascon"
	"gonum.org/v1/gonum/stat"
)

const (
	ExpectedValue  = 128
	ExpectedStdDev = 8.0
	CriticalValue  = 1.96
	Size           = 1000
	TestRuns       = 10000
)

func main() {
	testHashMultipleTimes("SHA256", sha256.New)
	testHashMultipleTimes("SHA3-256", func() hash.Hash { return sha3.New256() })
	testHashMultipleTimes("ASCON", func() hash.Hash { return ascon.NewHash256() })
}

func testHashMultipleTimes(name string, newHash func() hash.Hash) {
	passCount := 0
	zScores := make([]float64, TestRuns)

	for i := range TestRuns {
		hash := newHash()
		passed, zScore := hammingTestWithScore(hash)
		zScores[i] = zScore

		if passed {
			passCount++
		}
	}

	passRate := float64(passCount) / float64(TestRuns) * 100
	meanZScore := stat.Mean(zScores, nil)
	stdDevZScore := stat.StdDev(zScores, nil)

	fmt.Printf("\nResults for %s:\n", name)
	fmt.Printf("  Tests passed: %d/%d (%.2f%%)\n", passCount, TestRuns, passRate)
	fmt.Printf("  Mean Z-score: %.4f\n", meanZScore)
	fmt.Printf("  StdDev Z-score: %.4f\n", stdDevZScore)
}

func hammingTestWithScore(hash hash.Hash) (bool, float64) {
	distances := getHammingDistancesForHash(hash)
	mean := stat.Mean(distances, nil)
	zhd := (mean - ExpectedValue) * math.Sqrt(Size) / ExpectedStdDev
	return math.Abs(zhd) <= CriticalValue, zhd
}

func getHammingDistancesForHash(hash hash.Hash) []float64 {
	distances := make([]float64, Size)
	for i := range Size {
		data1 := generateRandomData(64)
		data2 := changeRandomBitInData(data1)

		hash.Write(data1)
		hash1 := hash.Sum(nil)

		hash.Reset()

		hash.Write(data2)
		hash2 := hash.Sum(nil)

		distance, err := calculateHammingDistance(hash1, hash2)
		if err != nil {
			log.Fatalf("error during calculating hamming distance: %v", err)
		}

		distances[i] = float64(distance)
	}

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

func generateRandomData(bytes int) []byte {
	res := make([]byte, bytes)
	rand.Read(res)
	return res
}

func changeRandomBitInData(data []byte) []byte {
	result := slices.Clone(data)
	byteIndex := mathrand.IntN(len(data))
	bitIndex := mathrand.IntN(8)
	result[byteIndex] = result[byteIndex] ^ (1 << bitIndex)
	return result
}
