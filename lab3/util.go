package main

import (
	"crypto/rand"
	mathrand "math/rand/v2"
	"slices"
)

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
