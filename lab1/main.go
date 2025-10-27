package main

import (
	"fmt"
	"sort"
	"strings"
)

type FreqPair struct {
	letter rune
	count  int
}

func main() {
	ciphertext := "YAE VCCX ZA FWSG ZRC I IVX WZ QEUZ NC SAFPWQC OWZR ZRC LCVMZR AH ZRC ILFRINCZ"

	fmt.Println("=== START ===")
	fmt.Println()
	fmt.Println("Ciphertext:", ciphertext)
	fmt.Println()

	freq := analyzeFrequency(ciphertext)
	displayFrequency(freq)

	fmt.Println()
	fmt.Println("=== KEY BREAKING ATTEMPTS ===")
	fmt.Println()

	mostCommon := getMostCommon(freq, 2)
	fmt.Printf("Most frequent letters in ciphertext: %c (%d occurrences), %c (%d occurrences)\n\n",
		mostCommon[0].letter, mostCommon[0].count,
		mostCommon[1].letter, mostCommon[1].count)

	englishCommon := []rune{'E', 'T'}

	attempts := 0
	for i := range 2 {
		for j := range 2 {
			attempts++
			c1 := mostCommon[i].letter
			c2 := mostCommon[j].letter
			p1 := englishCommon[0]
			p2 := englishCommon[1]

			if i == j {
				attempts--
				continue
			}

			fmt.Printf("Attempt %d: Assuming %c→%c and %c→%c\n", attempts, p1, c1, p2, c2)

			a, b, valid := solveAffineSystem(p1, c1, p2, c2)

			if !valid {
				fmt.Println("System of equations has no solution or a is not invertible modulo 26")
				fmt.Println()
				continue
			}

			fmt.Printf("  Key: a=%d, b=%d\n", a, b)

			aInv := modInverse(a, 26)
			if aInv == -1 {
				fmt.Println("No modular inverse for a\n")
				continue
			}

			plaintext := decryptAffine(ciphertext, aInv, b)
			fmt.Printf("  Plaintext: %s\n", plaintext)

			if looksLikeEnglish(plaintext) {
				fmt.Println("SUCCESS!")
				fmt.Printf("Encryption key: a=%d, b=%d\n", a, b)
				fmt.Printf("Decryption key: a_inv=%d, b=%d\n", aInv, b)
				fmt.Printf("Plaintext: %s\n", plaintext)
				return
			} else {
				fmt.Println("FAIL!")
			}
		}
	}

	fmt.Println("Failed to find key with basic assumptions.")
	fmt.Println("Trying other combinations...")

	bruteForceAttack(ciphertext, freq)
}

func analyzeFrequency(text string) map[rune]int {
	freq := make(map[rune]int)
	for _, ch := range text {
		if ch >= 'A' && ch <= 'Z' {
			freq[ch]++
		}
	}
	return freq
}

func displayFrequency(freq map[rune]int) {
	fmt.Println("=== LETTER FREQUENCY HISTOGRAM ===\n")

	pairs := make([]FreqPair, 0, len(freq))
	total := 0
	for letter, count := range freq {
		pairs = append(pairs, FreqPair{letter, count})
		total += count
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].count > pairs[j].count
	})

	for _, pair := range pairs {
		percentage := float64(pair.count) / float64(total) * 100
		bar := strings.Repeat("█", pair.count)
		fmt.Printf("%c: %3d (%.1f%%) %s\n", pair.letter, pair.count, percentage, bar)
	}
}

func getMostCommon(freq map[rune]int, n int) []FreqPair {
	pairs := make([]FreqPair, 0, len(freq))
	for letter, count := range freq {
		pairs = append(pairs, FreqPair{letter, count})
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].count > pairs[j].count
	})

	if len(pairs) > n {
		pairs = pairs[:n]
	}
	return pairs
}

func solveAffineSystem(p1, c1, p2, c2 rune) (int, int, bool) {
	x1 := int(p1 - 'A')
	y1 := int(c1 - 'A')
	x2 := int(p2 - 'A')
	y2 := int(c2 - 'A')

	deltaX := mod(x1-x2, 26)
	deltaY := mod(y1-y2, 26)

	deltaXInv := modInverse(deltaX, 26)
	if deltaXInv == -1 {
		return 0, 0, false
	}

	a := mod(deltaY*deltaXInv, 26)

	if gcd(a, 26) != 1 {
		return 0, 0, false
	}

	b := mod(y1-a*x1, 26)

	return a, b, true
}

func decryptAffine(ciphertext string, aInv, b int) string {
	var result strings.Builder
	for _, ch := range ciphertext {
		if ch >= 'A' && ch <= 'Z' {
			y := int(ch - 'A')
			x := mod(aInv*(y-b), 26)
			result.WriteRune(rune('A' + x))
		} else {
			result.WriteRune(ch)
		}
	}
	return result.String()
}

func looksLikeEnglish(text string) bool {
	words := strings.Fields(text)
	commonWords := map[string]bool{
		"THE": true, "AND": true, "TO": true, "OF": true, "A": true,
		"IN": true, "IS": true, "IT": true, "YOU": true, "THAT": true,
		"WAS": true, "FOR": true, "ON": true, "ARE": true, "WITH": true,
		"AS": true, "I": true, "BE": true, "THIS": true, "AT": true,
	}

	validWords := 0
	for _, word := range words {
		if commonWords[word] {
			validWords++
		}
	}

	return validWords >= 2 || (len(words) > 5 && float64(validWords)/float64(len(words)) > 0.3)
}

func bruteForceAttack(ciphertext string, freq map[rune]int) {
	mostCommon := getMostCommon(freq, 5)
	englishCommon := []rune{'E', 'T', 'A', 'O', 'I'}

	for i := 0; i < len(mostCommon) && i < 5; i++ {
		for j := 0; j < len(mostCommon) && j < 5; j++ {
			if i == j {
				continue
			}
			for pi := range englishCommon {
				for pj := range englishCommon {
					if pi == pj {
						continue
					}

					a, b, valid := solveAffineSystem(
						englishCommon[pi], mostCommon[i].letter,
						englishCommon[pj], mostCommon[j].letter,
					)

					if !valid {
						continue
					}

					aInv := modInverse(a, 26)
					if aInv == -1 {
						continue
					}

					plaintext := decryptAffine(ciphertext, aInv, b)

					if looksLikeEnglish(plaintext) {
						fmt.Printf("\nKEY FOUND!\n")
						fmt.Printf("Assumption: %c→%c and %c→%c\n",
							englishCommon[pi], mostCommon[i].letter,
							englishCommon[pj], mostCommon[j].letter)
						fmt.Printf("Key: a=%d, b=%d\n", a, b)
						fmt.Printf("Plaintext: %s\n", plaintext)
						return
					}
				}
			}
		}
	}
}

func mod(a, m int) int {
	result := a % m
	if result < 0 {
		result += m
	}
	return result
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func modInverse(a, m int) int {
	if gcd(a, m) != 1 {
		return -1
	}

	for x := 1; x < m; x++ {
		if mod(a*x, m) == 1 {
			return x
		}
	}
	return -1
}
