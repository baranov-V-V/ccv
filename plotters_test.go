package main

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"
)

func readCSVToChartEntries(filepath string) ([]ChartEntry, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	entries := make([]ChartEntry, 0, len(records))
	for _, record := range records {
		complexity, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, err
		}

		churn, err := strconv.ParseUint(record[2], 10, 64)
		if err != nil {
			return nil, err
		}

		entry := ChartEntry{
			File:   record[0],
			Complexity: complexity,
			Churn:      uint(churn),
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func createTestChart(t *testing.T, entries []ChartEntry, outputPath string) {
	t.Helper()
	err := CreateComplexityChurnChart(entries, outputPath)
	if err != nil {
		t.Fatalf("Failed to create chart: %v", err)
	}

	_, err = os.Stat(outputPath)
	if err != nil {
		t.Fatalf("Output file was not created: %v", err)
	}
}

func TestCreateScatterChart200(t *testing.T) {
	entries, err := readCSVToChartEntries("test/data/plot_200.csv")
	if err != nil {
		t.Fatalf("Failed to read CSV data: %v", err)
	}

	createTestChart(t, entries, "test/charts/scatter-200.html")
}

func TestCreateScatterChart2000(t *testing.T) {
	entries, err := readCSVToChartEntries("test/data/plot_2000.csv")
	if err != nil {
		t.Fatalf("Failed to read CSV data: %v", err)
	}

	createTestChart(t, entries, "test/charts/scatter-2000.html")
}

func TestCreateScatterChart10SameValues(t *testing.T) {
	entries, err := readCSVToChartEntries("test/data/plot_10-same.csv")
	if err != nil {
		t.Fatalf("Failed to read CSV data: %v", err)
	}

	createTestChart(t, entries, "test/charts/scatter-10-same.html")
}