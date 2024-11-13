package plot

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func readCSVToChartEntries(filepath string) ([]ScatterEntry, error) {
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

	entries := make([]ScatterEntry, 0, len(records))
	for _, record := range records {
		complexity, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, err
		}

		churn, err := strconv.ParseUint(record[2], 10, 64)
		if err != nil {
			return nil, err
		}

		entry := ScatterEntry{
			File:        record[0],
			ScatterData: ScatterData{Complexity: complexity, Churn: uint(churn)},
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func createTestChart(t *testing.T, entries []ScatterEntry, outputPath string) {
	t.Helper()
	err := CreateScatterChart(entries, NewRisksMapper(), outputPath)
	if err != nil {
		t.Fatalf("Failed to create chart: %v", err)
	}

	_, err = os.Stat(outputPath)
	if err != nil {
		t.Fatalf("Output file was not created: %v", err)
	}
}

var CSVDataDir = "../../test/data/"
var OutputDir = "../../test/charts/"

func TestCreateScatterChart200(t *testing.T) {
	entries, err := readCSVToChartEntries(CSVDataDir+"plot_200.csv")
	assert.NoError(t, err)

	createTestChart(t, entries, OutputDir+"scatter-200.html")
}

func TestCreateScatterChart2000(t *testing.T) {
	entries, err := readCSVToChartEntries(CSVDataDir+"plot_2000.csv")
	assert.NoError(t, err)

	createTestChart(t, entries, OutputDir+"scatter-2000.html")
}

func TestCreateScatterChart10SameValues(t *testing.T) {
	entries, err := readCSVToChartEntries(CSVDataDir + "plot_10-same.csv")
	assert.NoError(t, err)

	createTestChart(t, entries, OutputDir+"scatter-10-same.html")
}
