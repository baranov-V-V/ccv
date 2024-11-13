package plot

import (
	"testing"

	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/stretchr/testify/assert"
)

func TestGroupByFile(t *testing.T) {
	tests := []struct {
		name    string
		entries []ScatterEntry
		want    []groupedEntry
	}{
		{
			name:    "empty entries returns empty result",
			entries: []ScatterEntry{},
			want:    []groupedEntry{},
		},
		{
			name: "single file entry returns single group",
			entries: []ScatterEntry{
				{
					ScatterData: ScatterData{Complexity: 10.0, Churn: 5},
					File:        "file1.go",
				},
			},
			want: []groupedEntry{
				{
					ScatterData: ScatterData{Complexity: 10.0, Churn: 5},
					Files:       []string{"file1.go"},
				},
			},
		},
		{
			name: "multiple files with same metrics group together",
			entries: []ScatterEntry{
				{
					ScatterData: ScatterData{Complexity: 10.0, Churn: 5},
					File:        "file1.go",
				},
				{
					ScatterData: ScatterData{Complexity: 10.0, Churn: 5},
					File:        "file2.go",
				},
				{
					ScatterData: ScatterData{Complexity: 10.0, Churn: 5},
					File:        "file3.go",
				},
			},
			want: []groupedEntry{
				{
					ScatterData: ScatterData{Complexity: 10.0, Churn: 5},
					Files:       []string{"file1.go", "file2.go", "file3.go"},
				},
			},
		},
		{
			name: "multiple files multiple metrics group",
			entries: []ScatterEntry{
				{
					ScatterData: ScatterData{Complexity: 10.0, Churn: 5},
					File:        "file1.go",
				},
				{
					ScatterData: ScatterData{Complexity: 9.0, Churn: 3},
					File:        "file2.go",
				},
				{
					ScatterData: ScatterData{Complexity: 10.0, Churn: 5},
					File:        "file3.go",
				},
				{
					ScatterData: ScatterData{Complexity: 7.0, Churn: 3},
					File:        "file4.go",
				},
				{
					ScatterData: ScatterData{Complexity: 9.0, Churn: 5},
					File:        "file5.go",
				},
			},
			want: []groupedEntry{
				{
					ScatterData: ScatterData{Complexity: 10.0, Churn: 5},
					Files:       []string{"file1.go", "file3.go"},
				},
				{
					ScatterData: ScatterData{Complexity: 9.0, Churn: 3},
					Files:       []string{"file2.go"},
				},
				{
					ScatterData: ScatterData{Complexity: 7.0, Churn: 3},
					Files:       []string{"file4.go"},
				},
				{
					ScatterData: ScatterData{Complexity: 9.0, Churn: 5},
					Files:       []string{"file5.go"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := groupByFile(tt.entries)
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}

type stubMapper struct{}

func (m *stubMapper) Map(data ScatterData) Category {
	if data.Complexity >= 10.0 {
		return "critical"
	}
	if data.Complexity >= 5.0 {
		return "warning"
	}
	return "normal"
}

func (m *stubMapper) Style(_ Category) opts.ItemStyle {
	return opts.ItemStyle{}
}

func TestFormDataSeries(t *testing.T) {
	tests := []struct {
		name    string
		entries []ScatterEntry
		want    ScatterSeries
	}{
		{
			name:    "empty entries returns empty series",
			entries: []ScatterEntry{},
			want:    ScatterSeries{},
		},
		{
			name: "entries are grouped by same metrics and mapped to categories",
			entries: []ScatterEntry{
				{
					ScatterData: ScatterData{Complexity: 12.0, Churn: 5},
					File:        "critical1.go",
				},
				{
					ScatterData: ScatterData{Complexity: 12.0, Churn: 5},
					File:        "critical2.go",
				},
				{
					ScatterData: ScatterData{Complexity: 7.0, Churn: 3},
					File:        "warning.go",
				},
				{
					ScatterData: ScatterData{Complexity: 3.0, Churn: 1},
					File:        "normal1.go",
				},
				{
					ScatterData: ScatterData{Complexity: 1.0, Churn: 1},
					File:        "normal2.go",
				},
				{
					ScatterData: ScatterData{Complexity: 2.0, Churn: 1},
					File:        "normal3.go",
				},
			},
			want: ScatterSeries{
				"critical": {
					{
						Value:      []interface{}{12.0, uint(5), "critical1.go<br/>critical2.go"},
						Symbol:     "circle",
						SymbolSize: ScatterSymbolSize,
					},
				},
				"warning": {
					{
						Value:      []interface{}{7.0, uint(3), "warning.go"},
						Symbol:     "circle",
						SymbolSize: ScatterSymbolSize,
					},
				},
				"normal": {
					{
						Value:      []interface{}{3.0, uint(1), "normal1.go"},
						Symbol:     "circle",
						SymbolSize: ScatterSymbolSize,
					},
					{
						Value:      []interface{}{1.0, uint(1), "normal2.go"},
						Symbol:     "circle",
						SymbolSize: ScatterSymbolSize,
					},
					{
						Value:      []interface{}{2.0, uint(1), "normal3.go"},
						Symbol:     "circle",
						SymbolSize: ScatterSymbolSize,
					},
				},
			},
		},
	}

	mapper := &stubMapper{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formDataSeries(tt.entries, mapper)
			assert.Equal(t, tt.want, got)
		})
	}
}
