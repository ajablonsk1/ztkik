package main

import (
	"crypto/sha256"
	"crypto/sha3"
	"fmt"
	"hash"

	"github.com/magical/go-ascon"
)

const (
	CriticalValue = 1.96
	Size          = 10_000
	TestRuns      = 1000
)

func main() {
	hashes := map[string]func() hash.Hash{
		"SHA256":   sha256.New,
		"SHA3-256": func() hash.Hash { return sha3.New256() },
		"ASCON":    func() hash.Hash { return ascon.NewHash256() },
	}

	for name, hashFunc := range hashes {
		// Test Dystansu Hamminga
		testHashFuncMultipleTimes(fmt.Sprintf("Hamming Distance Test: %s", name), hashFunc, hammingTest, TestRuns)

		// Test Predykcji Bitów
		testHashFuncMultipleTimes(fmt.Sprintf("Bit Prediction Test: %s", name), hashFunc, bitsPredictionTest, TestRuns*256)
	}
}

func testHashFuncMultipleTimes(name string, newHash func() hash.Hash, testFunc func(func() hash.Hash) int, testRuns int) {
	passed := 0
	for i := 0; i < TestRuns; i++ {
		passed += testFunc(newHash)
	}
	printTestResults(name, passed, testRuns)
}

func printTestResults(name string, passCount, testRuns int) {
	passRate := float64(passCount) / float64(testRuns) * 100

	fmt.Printf("\nResults for %s:\n", name)
	fmt.Printf("  Tests passed: %d / %d (%.2f%%)\n", passCount, testRuns, passRate)
}
