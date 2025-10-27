package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"slices"
	"time"
)

type BenchmarkResult struct {
	mean         time.Duration
	median       time.Duration
	percentile95 time.Duration
	totalTime    time.Duration
}

type AlgorithmKeyGenResult struct {
	bits    int
	keysNum int
	result  *BenchmarkResult
}

type AlgorithmEncryptResult struct {
	bytes  int
	result *BenchmarkResult
}

type EncryptFunc func([]byte) (time.Duration, time.Duration)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "benchmarks" {
		fmt.Println("Starting cryptographic benchmarks...")
		fmt.Println()

		fmt.Println("=== KEY GENERATION BENCHMARKS ===")
		rsaKeyGen, aesKeyGen, desKeyGen := benchmarkKeyGeneration()

		fmt.Println("\n=== ENCRYPTION/DECRYPTION BENCHMARKS ===")
		encryptResults := benchmarkEncryptionDecryption()

		fmt.Println("=== EXPORT ===")

		fmt.Println("Exporting key gen results...")
		keyGenExports := map[string][]*AlgorithmKeyGenResult{
			"results/keygen/rsa2048.csv": filterAKGRByBits(rsaKeyGen, 2048),
			"results/keygen/rsa3072.csv": filterAKGRByBits(rsaKeyGen, 3072),
			"results/keygen/aes128.csv":  filterAKGRByBits(aesKeyGen, 128),
			"results/keygen/aes256.csv":  filterAKGRByBits(aesKeyGen, 256),
			"results/keygen/des192.csv":  desKeyGen,
		}
		for path, data := range keyGenExports {
			exportAKGR(data, path)
		}

		fmt.Println("Exporting encryption results...")
		encryptExports := map[string][]*AlgorithmEncryptResult{
			"results/encryption/rsa2048.csv": encryptResults["RSA-2048"].encrypt,
			"results/encryption/aes128.csv":  encryptResults["AES-128-GCM"].encrypt,
			"results/encryption/aes256.csv":  encryptResults["AES-256-GCM"].encrypt,
			"results/encryption/3des192.csv": encryptResults["3DES-CBC"].encrypt,
		}
		for path, data := range encryptExports {
			exportAER(data, path)
		}

		fmt.Println("Exporting decryption results...")
		decryptExports := map[string][]*AlgorithmEncryptResult{
			"results/decryption/rsa2048.csv": encryptResults["RSA-2048"].decrypt,
			"results/decryption/aes128.csv":  encryptResults["AES-128-GCM"].decrypt,
			"results/decryption/aes256.csv":  encryptResults["AES-256-GCM"].decrypt,
			"results/decryption/3des192.csv": encryptResults["3DES-CBC"].decrypt,
		}
		for path, data := range decryptExports {
			exportAER(data, path)
		}
	}

	fmt.Println("Drawing plots...")
	drawAll()
}

func exportToCSV(filepath string, records [][]string) {
	f, err := os.Create(filepath)
	if err != nil {
		log.Fatalf("error during file creation: %v", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	if err := w.WriteAll(records); err != nil {
		log.Fatalln("error writing csv:", err)
	}

	if err = w.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}
}

func exportAKGR(results []*AlgorithmKeyGenResult, filepath string) {
	records := [][]string{{"Number of keys", "Mean", "Median", "95. Percentile", "Total time"}}
	for _, result := range results {
		records = append(records, []string{
			fmt.Sprint(result.keysNum),
			result.result.mean.String(),
			result.result.median.String(),
			result.result.percentile95.String(),
			result.result.totalTime.String(),
		})
	}
	exportToCSV(filepath, records)
}

func exportAER(results []*AlgorithmEncryptResult, filepath string) {
	records := [][]string{{"Bytes", "Mean", "Median", "95. Percentile"}}
	for _, result := range results {
		records = append(records, []string{
			fmt.Sprint(result.bytes),
			result.result.mean.String(),
			result.result.median.String(),
			result.result.percentile95.String(),
		})
	}
	exportToCSV(filepath, records)
}

func filterAKGRByBits(results []*AlgorithmKeyGenResult, bits int) []*AlgorithmKeyGenResult {
	res := []*AlgorithmKeyGenResult{}
	for _, kgt := range results {
		if kgt.bits == bits {
			res = append(res, kgt)
		}
	}
	return res
}

func runKeyGenBenchmark(name string, keyNums []int, bitsToTest []int, measureFunc func(nums, bits int) []time.Duration) []*AlgorithmKeyGenResult {
	results := make([]*AlgorithmKeyGenResult, 0)
	for _, nums := range keyNums {
		for _, bits := range bitsToTest {
			fmt.Printf("Calculating for %s; %d bits; %d keys\n", name, bits, nums)
			durations := measureFunc(nums, bits)
			result := calculateBenchmarkResult(durations)
			results = append(results, &AlgorithmKeyGenResult{bits, nums, result})
		}
	}
	return results
}

func benchmarkKeyGeneration() ([]*AlgorithmKeyGenResult, []*AlgorithmKeyGenResult, []*AlgorithmKeyGenResult) {
	keyNums := []int{1, 10, 100, 1000}
	rsaBits := []int{2048, 3072}
	aesBits := []int{128, 256}
	desBits := []int{192}

	rsaKeyGenResults := runKeyGenBenchmark("RSA", keyNums, rsaBits, measureTimeRSA)
	aesKeyGenResults := runKeyGenBenchmark("AES", keyNums, aesBits, measureTimeAES)
	desKeyGenResults := runKeyGenBenchmark("3DES", keyNums, desBits, measureTime3DES)

	return rsaKeyGenResults, aesKeyGenResults, desKeyGenResults
}

func benchmarkEncryptionDecryption() map[string]struct {
	encrypt []*AlgorithmEncryptResult
	decrypt []*AlgorithmEncryptResult
} {
	iterations := 32

	symmetricDataSizes := []int{128, 512, 2 * 1024, 8 * 1024, 32 * 1024, 1024 * 1024, 4 * 1024 * 1024, 16 * 1024 * 1024}
	rsaDataSizes := []int{16, 32, 64, 128, 190}

	algorithms := map[string]struct {
		measureFunc EncryptFunc
		dataSizes   []int
	}{
		"RSA-2048": {
			measureFunc: measureEncryptRSA,
			dataSizes:   rsaDataSizes,
		},
		"AES-128-GCM": {
			measureFunc: func(data []byte) (time.Duration, time.Duration) {
				return measureEncryptAESGCM(data, 128)
			},
			dataSizes: symmetricDataSizes,
		},
		"AES-256-GCM": {
			measureFunc: func(data []byte) (time.Duration, time.Duration) {
				return measureEncryptAESGCM(data, 256)
			},
			dataSizes: symmetricDataSizes,
		},
		"3DES-CBC": {
			measureFunc: func(data []byte) (time.Duration, time.Duration) {
				return measureEncrypt3DESCBC(data, 192)
			},
			dataSizes: symmetricDataSizes,
		},
	}

	results := make(map[string]struct {
		encrypt []*AlgorithmEncryptResult
		decrypt []*AlgorithmEncryptResult
	})

	for name, algo := range algorithms {
		results[name] = struct {
			encrypt []*AlgorithmEncryptResult
			decrypt []*AlgorithmEncryptResult
		}{
			encrypt: make([]*AlgorithmEncryptResult, 0, len(algo.dataSizes)),
			decrypt: make([]*AlgorithmEncryptResult, 0, len(algo.dataSizes)),
		}
	}

	for algoName, algo := range algorithms {
		fmt.Printf("Running benchmarks for %s...\n", algoName)

		for _, bytes := range algo.dataSizes {
			data := make([]byte, bytes)
			_, err := rand.Read(data)
			if err != nil {
				log.Fatalf("Error generating random data: %v", err)
			}

			encryptResult, decryptResult := runBenchmark(data, algo.measureFunc, iterations)

			r := results[algoName]
			r.encrypt = append(r.encrypt, &AlgorithmEncryptResult{bytes, encryptResult})
			r.decrypt = append(r.decrypt, &AlgorithmEncryptResult{bytes, decryptResult})
			results[algoName] = r
		}
	}

	return results
}

func runBenchmark(data []byte, measureFunc EncryptFunc, iterations int) (*BenchmarkResult, *BenchmarkResult) {
	encryptDurations := make([]time.Duration, 0, iterations-2)
	decryptDurations := make([]time.Duration, 0, iterations-2)

	for i := range iterations {
		encryptDuration, decryptDuration := measureFunc(data)

		if i < 2 {
			continue
		}

		encryptDurations = append(encryptDurations, encryptDuration)
		decryptDurations = append(decryptDurations, decryptDuration)
	}

	encryptResult := calculateBenchmarkResult(encryptDurations)
	decryptResult := calculateBenchmarkResult(decryptDurations)

	return encryptResult, decryptResult
}

func calculateBenchmarkResult(durations []time.Duration) *BenchmarkResult {
	durationsLen := len(durations)
	if durationsLen == 0 {
		log.Println("Warning: calculateBenchmarkResult received empty durations slice. Returning zero result.")
		return &BenchmarkResult{}
	}

	totalTime := time.Duration(0)
	for _, duration := range durations {
		totalTime += duration
	}

	sorted := make([]time.Duration, durationsLen)
	copy(sorted, durations)
	slices.Sort(sorted)

	return &BenchmarkResult{
		mean:         totalTime / time.Duration(durationsLen),
		median:       calculatePercentile(sorted, 50),
		percentile95: calculatePercentile(sorted, 95),
		totalTime:    totalTime,
	}
}

func calculatePercentile(durations []time.Duration, percentile float64) time.Duration {
	n := len(durations)
	if n == 0 {
		return 0
	}

	position := (percentile / 100) * float64(n+1)

	_, decimalPart := math.Modf(position)
	integerPart := int(position)

	if integerPart < 1 {
		return durations[0]
	}
	if integerPart >= n {
		return durations[n-1]
	}
	lowerValue := float64(durations[integerPart-1])
	upperValue := float64(durations[integerPart])

	interpolated := lowerValue + decimalPart*(upperValue-lowerValue)

	return time.Duration(interpolated)
}

func measureKeyGen(keyNums int, keyGenFunc func() error) []time.Duration {
	// warm up
	if err := keyGenFunc(); err != nil {
		log.Fatalf("Error during warm-up key generation: %v\n", err)
	}

	generateKeyTimes := make([]time.Duration, 0, keyNums)

	for range keyNums {
		start := time.Now()
		if err := keyGenFunc(); err != nil {
			log.Fatalf("Error during key generation: %v\n", err)
		}
		generateKeyTimes = append(generateKeyTimes, time.Since(start))
	}

	return generateKeyTimes
}

func measureTimeRSA(keyNums, bits int) []time.Duration {
	return measureKeyGen(keyNums, func() error {
		_, err := rsa.GenerateKey(rand.Reader, bits)
		return err
	})
}

func measureTimeAES(keyNums, bits int) []time.Duration {
	keySize := bits / 8
	return measureKeyGen(keyNums, func() error {
		key := make([]byte, keySize)
		if _, err := rand.Read(key); err != nil {
			return fmt.Errorf("error creating random secret key for AES: %v", err)
		}
		_, err := aes.NewCipher(key)
		return err
	})
}

func measureTime3DES(keyNums, bits int) []time.Duration {
	keySize := bits / 8
	return measureKeyGen(keyNums, func() error {
		key := make([]byte, keySize)
		if _, err := rand.Read(key); err != nil {
			return fmt.Errorf("error creating random secret key for 3DES: %v", err)
		}
		_, err := des.NewTripleDESCipher(key)
		return err
	})
}

func measureEncryptRSA(data []byte) (time.Duration, time.Duration) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Error generating RSA key: %v", err)
	}
	publicKey := &privateKey.PublicKey

	startEncrypt := time.Now()
	ciphertext, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		publicKey,
		data,
		nil,
	)
	encryptDuration := time.Since(startEncrypt)
	if err != nil {
		log.Fatalf("Error encrypting: %v", err)
	}

	startDecrypt := time.Now()
	_, err = rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		privateKey,
		ciphertext,
		nil,
	)
	decryptDuration := time.Since(startDecrypt)
	if err != nil {
		log.Fatalf("Error decrypting: %v", err)
	}

	return encryptDuration, decryptDuration
}

func measureEncryptAESGCM(data []byte, keySizeBits int) (time.Duration, time.Duration) {
	key := make([]byte, keySizeBits/8)
	if _, err := rand.Read(key); err != nil {
		log.Fatalf("Error generating AES key: %v", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf("Error creating AES cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatalf("Error creating GCM: %v", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		log.Fatalf("Error generating nonce: %v", err)
	}

	startEncrypt := time.Now()
	ciphertext := gcm.Seal(nil, nonce, data, nil)
	encryptDuration := time.Since(startEncrypt)

	startDecrypt := time.Now()
	_, err = gcm.Open(nil, nonce, ciphertext, nil)
	decryptDuration := time.Since(startDecrypt)
	if err != nil {
		log.Fatalf("Error decrypting AES-GCM: %v", err)
	}

	return encryptDuration, decryptDuration
}

func measureEncrypt3DESCBC(data []byte, keySizeBits int) (time.Duration, time.Duration) {
	key := make([]byte, keySizeBits/8)
	if _, err := rand.Read(key); err != nil {
		log.Fatalf("Error generating 3DES key: %v", err)
	}

	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		log.Fatalf("Error creating 3DES cipher: %v", err)
	}

	iv := make([]byte, block.BlockSize())
	if _, err := rand.Read(iv); err != nil {
		log.Fatalf("Error generating IV: %v", err)
	}

	startEncrypt := time.Now()
	paddedData := pkcs7Pad(data, block.BlockSize())
	ciphertext := make([]byte, len(paddedData))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, paddedData)
	encryptDuration := time.Since(startEncrypt)

	startDecrypt := time.Now()
	plaintext := make([]byte, len(ciphertext))
	modeDecrypt := cipher.NewCBCDecrypter(block, iv)
	modeDecrypt.CryptBlocks(plaintext, ciphertext)
	_, err = pkcs7Unpad(plaintext, block.BlockSize())
	decryptDuration := time.Since(startDecrypt)
	if err != nil {
		log.Fatalf("Error removing padding: %v", err)
	}

	return encryptDuration, decryptDuration
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padText := make([]byte, padding)
	for i := range padText {
		padText[i] = byte(padding)
	}
	return append(data, padText...)
}

func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("data is empty")
	}
	padding := int(data[len(data)-1])
	if padding > blockSize || padding == 0 {
		return nil, fmt.Errorf("invalid padding: %d", padding)
	}
	if len(data) < padding {
		return nil, fmt.Errorf("invalid padding: data shorter than padding size")
	}
	for i := len(data) - padding; i < len(data); i++ {
		if data[i] != byte(padding) {
			return nil, fmt.Errorf("invalid padding byte")
		}
	}
	return data[:len(data)-padding], nil
}
