package plot

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"golang.org/x/exp/maps"
)

type ScatterData struct {
	Complexity float64
	Churn      uint
}

type ScatterEntry struct {
	ScatterData
	File string
}

type groupedEntry struct {
	ScatterData
	Files []string
}

type Category = string

type ScatterSeries map[Category][]opts.ScatterData

// TODO: Maybe categorizer?
type EntryMapper interface {
	Map(ScatterData) Category
	Style(Category) opts.ItemStyle
}

func groupByFile(entries []ScatterEntry) []groupedEntry {
	groups := make(map[ScatterData]groupedEntry)

	for _, entry := range entries {
		group, exists := groups[entry.ScatterData]
		if !exists {
			group = groupedEntry{ScatterData: entry.ScatterData}
		}

		group.Files = append(group.Files, entry.File)
		groups[entry.ScatterData] = group
	}

	return maps.Values(groups)
}

func formDataSeries(entries []ScatterEntry, mapper EntryMapper) ScatterSeries {
	series := make(ScatterSeries)

	groupedEntries := groupByFile(entries)

	for _, entry := range groupedEntries {
		category := mapper.Map(entry.ScatterData)
		series[category] = append(series[category], opts.ScatterData{
			Value:      []interface{}{entry.Complexity, entry.Churn, strings.Join(entry.Files, "<br/>")},
			Symbol:     "circle",
			SymbolSize: ScatterSymbolSize,
		})
	}

	return series
}

// TODO: pass EntryMapper as parameter!
func CreateScatterChart(entries []ScatterEntry, mapper EntryMapper, outputPath string) error {
	scatter := charts.NewScatter()
	scatter.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Code Complexity vs Churn",
			Top:   "0%",
			Left:  "center",
			Show:  opts.Bool(false),
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:    opts.Bool(true),
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
		// charts.WithColorsOpts(getRiskColors(riskLevels)),
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
			Width:  fmt.Sprintf("%dpx", WidthPx),
			Height: fmt.Sprintf("%dpx", HeightPx),
		}),
	)

	for category, data := range formDataSeries(entries, mapper) {
		scatter.AddSeries(category, data).SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show: opts.Bool(false),
			}),
			charts.WithItemStyleOpts(mapper.Style(category)),
		)
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return scatter.Render(f)
}
