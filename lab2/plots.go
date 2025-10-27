package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

const (
	plotDir    = "results/plots/"
	keygenDir  = "results/keygen/"
	encryptDir = "results/encryption/"
	decryptDir = "results/decryption/"
	plotSize   = 8 * vg.Inch

	mbDivider    = 1024 * 1024
	kbDivider    = 1024
	pointsLimit  = 8
	pointsLimit4 = 4
)

type PlotSeries struct {
	Name    string
	GetData func() (plotter.XYs, error)
}

type PlotConfig struct {
	Title    string
	XLabel   string
	YLabel   string
	Filepath string
	Series   []PlotSeries
}

func drawAndSavePlot(title, xLabel, yLabel, filepath string, args ...interface{}) {
	p := plot.New()
	p.Title.Text = title
	p.X.Label.Text = xLabel
	p.Y.Label.Text = yLabel

	if len(args) == 0 {
		log.Printf("  [!] Skipping plot %s: no data series provided", title)
		return
	}

	err := plotutil.AddLinePoints(p, args...)
	if err != nil {
		log.Printf("  [!] Error adding line points for %s: %v", title, err)
	}

	if err := p.Save(plotSize, plotSize, filepath); err != nil {
		log.Printf("  [!] ERROR saving plot %s: %v", filepath, err)
	}
}

func drawAll() {
	plots := []PlotConfig{
		{
			Title: "Asymmetric algorithms comparison", XLabel: "Number of keys", YLabel: "Total Time (s)",
			Filepath: plotDir + "keygen_asymmetric.png",
			Series: []PlotSeries{
				{"RSA 2048", func() (plotter.XYs, error) { return getPointsKeyGen(keygenDir + "rsa2048.csv") }},
				{"RSA 3072", func() (plotter.XYs, error) { return getPointsKeyGen(keygenDir + "rsa3072.csv") }},
			},
		},
		{
			Title: "Symmetric algorithms comparison", XLabel: "Number of keys", YLabel: "Total Time (s)",
			Filepath: plotDir + "keygen_symmetric.png",
			Series: []PlotSeries{
				{"AES 128", func() (plotter.XYs, error) { return getPointsKeyGen(keygenDir + "aes128.csv") }},
				{"AES 256", func() (plotter.XYs, error) { return getPointsKeyGen(keygenDir + "aes256.csv") }},
				{"DES 192", func() (plotter.XYs, error) { return getPointsKeyGen(keygenDir + "des192.csv") }},
			},
		},
		{
			Title: "All algorithms", XLabel: "Number of keys", YLabel: "Total Time (s)",
			Filepath: plotDir + "keygen_all.png",
			Series: []PlotSeries{
				{"RSA 2048", func() (plotter.XYs, error) { return getPointsKeyGen(keygenDir + "rsa2048.csv") }},
				{"RSA 3072", func() (plotter.XYs, error) { return getPointsKeyGen(keygenDir + "rsa3072.csv") }},
				{"AES 128", func() (plotter.XYs, error) { return getPointsKeyGen(keygenDir + "aes128.csv") }},
				{"AES 256", func() (plotter.XYs, error) { return getPointsKeyGen(keygenDir + "aes256.csv") }},
				{"DES 192", func() (plotter.XYs, error) { return getPointsKeyGen(keygenDir + "des192.csv") }},
			},
		},
		{
			Title: "Encryption - All algorithms", XLabel: "Size (MBs)", YLabel: "Mean Time (ms)",
			Filepath: plotDir + "encryption_all.png",
			Series: []PlotSeries{
				{"RSA 2048", func() (plotter.XYs, error) {
					return getPointsEncryptionTime(encryptDir+"rsa2048.csv", pointsLimit, mbDivider, time.Millisecond)
				}},
				{"AES 128", func() (plotter.XYs, error) {
					return getPointsEncryptionTime(encryptDir+"aes128.csv", pointsLimit, mbDivider, time.Millisecond)
				}},
				{"AES 256", func() (plotter.XYs, error) {
					return getPointsEncryptionTime(encryptDir+"aes256.csv", pointsLimit, mbDivider, time.Millisecond)
				}},
				{"DES 192", func() (plotter.XYs, error) {
					return getPointsEncryptionTime(encryptDir+"3des192.csv", pointsLimit, mbDivider, time.Millisecond)
				}},
			},
		},
		{
			Title: "Encryption all algorithms (4 points)", XLabel: "Size (KBs)", YLabel: "Mean Time (μs)",
			Filepath: plotDir + "encryption_all_4points.png",
			Series: []PlotSeries{
				{"RSA 2048", func() (plotter.XYs, error) {
					return getPointsEncryptionTime(encryptDir+"rsa2048.csv", pointsLimit4, kbDivider, time.Microsecond)
				}},
				{"AES 128", func() (plotter.XYs, error) {
					return getPointsEncryptionTime(encryptDir+"aes128.csv", pointsLimit4, kbDivider, time.Microsecond)
				}},
				{"AES 256", func() (plotter.XYs, error) {
					return getPointsEncryptionTime(encryptDir+"aes256.csv", pointsLimit4, kbDivider, time.Microsecond)
				}},
				{"DES 192", func() (plotter.XYs, error) {
					return getPointsEncryptionTime(encryptDir+"3des192.csv", pointsLimit4, kbDivider, time.Microsecond)
				}},
			},
		},
		{
			Title: "Decryption - All algorithms", XLabel: "Size (MBs)", YLabel: "Mean Time (ms)",
			Filepath: plotDir + "decryption_all.png",
			Series: []PlotSeries{
				{"RSA 2048", func() (plotter.XYs, error) {
					return getPointsEncryptionTime(decryptDir+"rsa2048.csv", pointsLimit, mbDivider, time.Millisecond)
				}},
				{"AES 128", func() (plotter.XYs, error) {
					return getPointsEncryptionTime(decryptDir+"aes128.csv", pointsLimit, mbDivider, time.Millisecond)
				}},
				{"AES 256", func() (plotter.XYs, error) {
					return getPointsEncryptionTime(decryptDir+"aes256.csv", pointsLimit, mbDivider, time.Millisecond)
				}},
				{"DES 192", func() (plotter.XYs, error) {
					return getPointsEncryptionTime(decryptDir+"3des192.csv", pointsLimit, mbDivider, time.Millisecond)
				}},
			},
		},
		{
			Title: "Decryption all algorithms (4 points)", XLabel: "Size (KBs)", YLabel: "Mean Time (μs)",
			Filepath: plotDir + "decryption_all_4points.png",
			Series: []PlotSeries{
				{"RSA 2048", func() (plotter.XYs, error) {
					return getPointsEncryptionTime(decryptDir+"rsa2048.csv", pointsLimit4, kbDivider, time.Microsecond)
				}},
				{"AES 128", func() (plotter.XYs, error) {
					return getPointsEncryptionTime(decryptDir+"aes128.csv", pointsLimit4, kbDivider, time.Microsecond)
				}},
				{"AES 256", func() (plotter.XYs, error) {
					return getPointsEncryptionTime(decryptDir+"aes256.csv", pointsLimit4, kbDivider, time.Microsecond)
				}},
				{"DES 192", func() (plotter.XYs, error) {
					return getPointsEncryptionTime(decryptDir+"3des192.csv", pointsLimit4, kbDivider, time.Microsecond)
				}},
			},
		},
		{
			Title: "Throughput encryption all algorithms", XLabel: "Size (MBs)", YLabel: "Throughput (MB/s)",
			Filepath: plotDir + "encryption_all_throughput.png",
			Series: []PlotSeries{
				{"RSA 2048", func() (plotter.XYs, error) { return getPointsEncryptionThroughput(encryptDir + "rsa2048.csv") }},
				{"AES 128", func() (plotter.XYs, error) { return getPointsEncryptionThroughput(encryptDir + "aes128.csv") }},
				{"AES 256", func() (plotter.XYs, error) { return getPointsEncryptionThroughput(encryptDir + "aes256.csv") }},
				{"DES 192", func() (plotter.XYs, error) { return getPointsEncryptionThroughput(encryptDir + "3des192.csv") }},
			},
		},
		{
			Title: "Throughput decryption all algorithms", XLabel: "Size (MBs)", YLabel: "Throughput (MB/s)",
			Filepath: plotDir + "decryption_all_throughput.png",
			Series: []PlotSeries{
				{"RSA 2048", func() (plotter.XYs, error) { return getPointsEncryptionThroughput(decryptDir + "rsa2048.csv") }},
				{"AES 128", func() (plotter.XYs, error) { return getPointsEncryptionThroughput(decryptDir + "aes128.csv") }},
				{"AES 256", func() (plotter.XYs, error) { return getPointsEncryptionThroughput(decryptDir + "aes256.csv") }},
				{"DES 192", func() (plotter.XYs, error) { return getPointsEncryptionThroughput(decryptDir + "3des192.csv") }},
			},
		},
	}

	for _, config := range plots {
		log.Printf("Drawing plot: %s", config.Title)

		args := make([]interface{}, 0, len(config.Series)*2)

		for _, series := range config.Series {
			data, err := series.GetData()
			if err != nil {
				log.Printf("  [!] Skipping series '%s' for plot '%s': %v", series.Name, config.Title, err)
				continue
			}
			args = append(args, series.Name, data)
		}

		drawAndSavePlot(config.Title, config.XLabel, config.YLabel, config.Filepath, args...)
	}
}

func getPointsEncryptionTime(filepath string, pointsLimit int, xDivider float64, yUnit time.Duration) (plotter.XYs, error) {
	pts := make(plotter.XYs, 0)
	res, err := readCsvFile(filepath)
	if err != nil {
		return nil, err
	}

	for i, row := range res {
		if i+1 > pointsLimit {
			break
		}
		if len(row) < 2 {
			return nil, fmt.Errorf("invalid row in %s: expected at least 2 columns, got %d", filepath, len(row))
		}

		point := plotter.XY{}

		totalBytes, err := strconv.ParseFloat(row[0], 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing bytes in %s (row %d): %w", filepath, i, err)
		}

		mean, err := time.ParseDuration(row[1])
		if err != nil {
			return nil, fmt.Errorf("error parsing mean in %s (row %d): %w", filepath, i, err)
		}

		point.X = totalBytes / xDivider
		point.Y = float64(mean / yUnit)
		pts = append(pts, point)
	}
	return pts, nil
}

func getPointsKeyGen(filepath string) (plotter.XYs, error) {
	pts := make(plotter.XYs, 0)
	res, err := readCsvFile(filepath)
	if err != nil {
		return nil, err
	}

	for i, row := range res {
		if len(row) < 5 {
			return nil, fmt.Errorf("invalid row in %s: expected at least 5 columns, got %d", filepath, len(row))
		}
		point := plotter.XY{}

		keysNum, err := strconv.ParseFloat(row[0], 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing keys num in %s (row %d): %w", filepath, i, err)
		}

		totalTime, err := time.ParseDuration(row[4])
		if err != nil {
			return nil, fmt.Errorf("error parsing total time in %s (row %d): %w", filepath, i, err)
		}

		point.X = keysNum
		point.Y = totalTime.Seconds()
		pts = append(pts, point)
	}
	return pts, nil
}

func readCsvFile(filepath string) ([][]string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error while opening a file %s: %w", filepath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error while reading csv file %s: %w", filepath, err)
	}

	if len(records) < 1 {
		return nil, fmt.Errorf("empty csv file: %s", filepath)
	}
	return records[1:], nil
}

func getPointsEncryptionThroughput(filepath string) (plotter.XYs, error) {
	pts := make(plotter.XYs, 0)
	res, err := readCsvFile(filepath)
	if err != nil {
		return nil, err
	}

	for i, row := range res {
		if len(row) < 2 {
			return nil, fmt.Errorf("invalid row in %s: expected at least 2 columns, got %d", filepath, len(row))
		}

		point := plotter.XY{}

		totalBytes, err := strconv.ParseFloat(row[0], 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing bytes in %s (row %d): %w", filepath, i, err)
		}

		mean, err := time.ParseDuration(row[1])
		if err != nil {
			return nil, fmt.Errorf("error parsing mean in %s (row %d): %w", filepath, i, err)
		}

		sizeMB := totalBytes / mbDivider
		timeS := mean.Seconds()

		point.X = sizeMB
		if timeS > 0 {
			point.Y = sizeMB / timeS
		} else {
			point.Y = 0
		}
		pts = append(pts, point)
	}
	return pts, nil
}
