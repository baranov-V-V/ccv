package plot

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

var (
	VeryLowRisk  uint = 10
	LowRisk      uint = 15
	MediumRisk   uint = 20
	HighRisk     uint = 25
	VeryHighRisk uint = 30
	CriticalRisk uint = 35
)

// If need to show scroll in chart
var WithScroll = false

// If need to devide risks into categories
var WithRisks = false

// Need to make it more general TODO refactor
func ValidateRiskThresholds() error {
	if VeryLowRisk >= LowRisk {
		return fmt.Errorf("Very Low Risk threshold (%d) must be less than Low Risk threshold (%d)",
			VeryLowRisk, LowRisk)
	}
	if LowRisk >= MediumRisk {
		return fmt.Errorf("Low Risk threshold (%d) must be less than Medium Risk threshold (%d)",
			LowRisk, MediumRisk)
	}
	if MediumRisk >= HighRisk {
		return fmt.Errorf("Medium Risk threshold (%d) must be less than High Risk threshold (%d)",
			MediumRisk, HighRisk)
	}
	if HighRisk >= VeryHighRisk {
		return fmt.Errorf("High Risk threshold (%d) must be less than Very High Risk threshold (%d)",
			HighRisk, VeryHighRisk)
	}
	if VeryHighRisk >= CriticalRisk {
		return fmt.Errorf("Very High Risk threshold (%d) must be less than Critical Risk threshold (%d)",
			VeryHighRisk, CriticalRisk)
	}
	return nil
}

type ChartEntry struct {
	File       string
	Complexity float64
	Churn      uint
}

type RiskLevel struct {
	Name  string
	Color string
	Min   uint
	Max   uint
}

func getRiskLevels() []RiskLevel {
	return []RiskLevel{
		{Name: "Very Low Risk", Color: "#90EE90", Min: VeryLowRisk, Max: LowRisk - 1},
		{Name: "Low Risk", Color: "#47d147", Min: LowRisk, Max: MediumRisk - 1},
		{Name: "Medium Risk", Color: "#ffd700", Min: MediumRisk, Max: HighRisk - 1},
		{Name: "High Risk", Color: "#ffa64d", Min: HighRisk, Max: VeryHighRisk - 1},
		{Name: "Very High Risk", Color: "#ff4d4d", Min: VeryHighRisk, Max: CriticalRisk - 1},
		{Name: "Critical Risk", Color: "#8b0000", Min: CriticalRisk, Max: ^uint(0)},
	}
}

func CreateComplexityChurnChart(entries []ChartEntry, outputPath string) error {
	show := true
	riskLevels := getRiskLevels()

	riskMaps := make(map[string]map[string][]string)
	for _, level := range riskLevels {
		riskMaps[level.Name] = make(map[string][]string)
	}

	// Group files by risk level and coordinates
	for _, entry := range entries {
		key := fmt.Sprintf("%f-%d", entry.Complexity, entry.Churn)
		riskScore := entry.Complexity + float64(entry.Churn)

		for _, level := range riskLevels {
			if riskScore >= float64(level.Min) && riskScore <= float64(level.Max) {
				riskMaps[level.Name][key] = append(riskMaps[level.Name][key], entry.File)
				break
			}
		}
	}

	scatter := charts.NewScatter()
	scatter.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Code Complexity vs Churn",
			Top:   "0%",
			Left:  "center",
			Show:  opts.Bool(false),
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:    &show,
			Trigger: "item",
			Formatter: opts.FuncOpts(`function(params) {
				return 'Complexity: ' + params.value[0] + 
					   '<br/>Churn: ' + params.value[1] + 
					   '<br/>Files:<br/>' + params.value[2];
			}`),
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name:  "Complexity",
			Type:  "value",
			Scale: opts.Bool(true),
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name:  "Churn",
			Type:  "value",
			Scale: opts.Bool(true),
		}),
		charts.WithColorsOpts(getRiskColors(riskLevels)),
		/*
			// Horizontal zoom slider
			charts.WithDataZoomOpts(opts.DataZoom{
				Type:       "slider",
				Start:     0,
				End:       100,
				XAxisIndex: []int{0},
			}),
			// Vertical zoom slider
			charts.WithDataZoomOpts(opts.DataZoom{
				Type:       "slider",
				Start:     0,
				End:       100,
				YAxisIndex: []int{0},
				Orient:    "vertical",
			}),
			// Inside zoom for both axes
			charts.WithDataZoomOpts(opts.DataZoom{
				Type:  "inside",
				Start: 0,
				End:   100,
			}),
		*/
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "1200px",
			Height: "800px",
		}),
	)

	// Convert maps to series data and add to chart
	for _, level := range riskLevels {
		var seriesData []opts.ScatterData
		for key, files := range riskMaps[level.Name] {
			var complexity float64
			var churn uint
			fmt.Sscanf(key, "%f-%d", &complexity, &churn)
			filesList := strings.Join(files, "<br/>")
			seriesData = append(seriesData, opts.ScatterData{
				Value:      []interface{}{complexity, churn, filesList},
				Symbol:     "circle",
				SymbolSize: 8,
			})
		}
		scatter.AddSeries(level.Name, seriesData).
			SetSeriesOptions(
				charts.WithLabelOpts(opts.Label{
					Show: opts.Bool(false),
				}),
			)
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return scatter.Render(f)
}

func getRiskColors(levels []RiskLevel) []string {
	colors := make([]string, len(levels))
	for i, level := range levels {
		colors[i] = level.Color
	}
	return colors
}
