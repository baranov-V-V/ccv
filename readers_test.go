package main

import (
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseItem(t *testing.T) {
	tests := []struct {
		name    string
		input   lizardItem
		want    FunctionStat
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid function with package",
			input: lizardItem{
				Name:   "pkg::subpkg::func_name at path/to/file.go:123",
				Values: []int{0, 3, 5},
			},
			want: FunctionStat{
				Name:      "func_name",
				Package:   []string{"pkg", "subpkg"},
				Line:      123,
				Length:    3,
				File:      "path/to/file.go",
				Compexity: 5,
			},
			wantErr: false,
		},
		{
			name: "valid function with parameters",
			input: lizardItem{
				Name:   "pkg::func_name(param1, param2) at path/file.go:456",
				Values: []int{0, 2, 3},
			},
			want: FunctionStat{
				Name:      "func_name",
				Package:   []string{"pkg"},
				Line:      456,
				Length:    2,
				File:      "path/file.go",
				Compexity: 3,
			},
			wantErr: false,
		},
		{
			name: "function with many packages",
			input: lizardItem{
				Name:   "pkg1::pkg2::pkg3::pkg4::pkg5::func_name(param1, param2) at path/file.go:456",
				Values: []int{0, 1, 4},
			},
			want: FunctionStat{
				Name:      "func_name",
				Package:   []string{"pkg1", "pkg2", "pkg3", "pkg4", "pkg5"},
				Line:      456,
				Length:    1,
				File:      "path/file.go",
				Compexity: 4,
			},
			wantErr: false,
		},
		{
			name: "function without package",
			input: lizardItem{
				Name:   "standalone_func at path/file.go:789",
				Values: []int{0, 0, 1},
			},
			want: FunctionStat{
				Name:      "standalone_func",
				Package:   []string{},
				Line:      789,
				Length:    0,
				File:      "path/file.go",
				Compexity: 1,
			},
			wantErr: false,
		},
		{
			name: "invalid format",
			input: lizardItem{
				Name:   "invalid_format",
				Values: []int{1, 1, 1},
			},
			wantErr: true,
			errMsg:  "invalid function format: invalid_format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseItem(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want.Name, got.Name)
			assert.Equal(t, tt.want.Package, got.Package)
			assert.Equal(t, tt.want.Line, got.Line)
			assert.Equal(t, tt.want.File, got.File)
			assert.Equal(t, tt.want.Compexity, got.Compexity)
		})
	}
}

func TestReadLizardXML(t *testing.T) {
	xmlData := `<?xml version="1.0" ?>
<cppncss>
    <measure type="Function">
        <labels>
            <label>Nr.</label>
            <label>NCSS</label>
            <label>CCN</label>
        </labels>
        <item name="pkg::some_func(...) at src/file1.cpp:1">
            <value>11</value>
            <value>111</value>
            <value>1111</value>
        </item>
        <item name="pkg::sub_pkg::SomeFunc(...) at path/to/src/file2.cpp:3">
            <value>22</value>
            <value>222</value>
            <value>2222</value>
        </item>
        <average label="NCSS" value="16.2"/>
        <average label="CCN" value="4.4"/>
    </measure>
    <measure type="File">
        <labels>
            <label>Nr.</label>
            <label>NCSS</label>
            <label>CCN</label>
            <label>Functions</label>
        </labels>
        <item name="src/file1.cpp">
            <value>1</value>
            <value>2</value>
            <value>3</value>
            <value>4</value>
        </item>
        <item name="path/to/src/file2.cpp">
            <value>5</value>
            <value>6</value>
            <value>7</value>
            <value>8</value>
        </item>
        <average label="NCSS" value="1"/>
        <average label="CCN" value="2"/>
        <average label="Functions" value="3"/>
        <sum label="NCSS" value="4"/>
        <sum label="CCN" value="5"/>
        <sum label="Functions" value="6"/>
    </measure>
</cppncss>`

	reader := strings.NewReader(xmlData)
	got, err := ReadLizardXML(reader)

	assert.NoError(t, err)
	assert.NotNil(t, got)

	// Test Functions section
	idx := slices.IndexFunc(got.Measures, func(l lizardMeasure) bool {
		return l.Type == "Function"
	})
	assert.Greater(t, idx, -1)

	assert.Nil(t, got.Measures[idx].Sums)

	assert.Len(t, got.Measures[idx].Items, 2)
	assert.Contains(t, got.Measures[idx].Items, lizardItem{
		Name:   "pkg::some_func(...) at src/file1.cpp:1",
		Values: []int{11, 111, 1111},
	})
	assert.Contains(t, got.Measures[idx].Items, lizardItem{
		Name:   "pkg::sub_pkg::SomeFunc(...) at path/to/src/file2.cpp:3",
		Values: []int{22, 222, 2222},
	})

	// Test Files section
	idx = slices.IndexFunc(got.Measures, func(l lizardMeasure) bool {
		return l.Type == "File"
	})
	assert.Greater(t, idx, -1)

	assert.Len(t, got.Measures[idx].Items, 2)
	assert.Contains(t, got.Measures[idx].Items, lizardItem{
		Name:   "src/file1.cpp",
		Values: []int{1, 2, 3, 4},
	})
	assert.Contains(t, got.Measures[idx].Items, lizardItem{
		Name:   "path/to/src/file2.cpp",
		Values: []int{5, 6, 7, 8},
	})
}

func TestParseLizard(t *testing.T) {
	lizard := &lizardXML{
		Measures: []lizardMeasure{
			{
				Type: "File",
				Items: []lizardItem{
					{
						Name:   "src/file1.cpp",
						Values: []int{1, 2, 3, 4},
					},
				},
			},
			{
				Type: "Function",
				Items: []lizardItem{
					{
						Name:   "pkg::some_func(...) at src/file1.cpp:1",
						Values: []int{1, 2, 3},
					},
					{
						Name:   "pkg::sub_pkg::SomeFunc(...) at src/file1.cpp:3",
						Values: []int{4, 5, 6},
					},
				},
			},
		},
	}

	got, err := ParseLizard(lizard)
	assert.NoError(t, err)
	assert.Len(t, got, 1) // One file with two functions

	file := got[0]
	assert.Equal(t, "src/file1.cpp", file.Path)
	assert.Len(t, file.Functions, 2)

	// Check functions using Contains
	assert.Contains(t, file.Functions, FunctionStat{
		Name:      "some_func",
		Package:   []string{"pkg"},
		Line:      1,
		Length:    2,
		File:      "src/file1.cpp",
		Compexity: 3,
	})

	assert.Contains(t, file.Functions, FunctionStat{
		Name:      "SomeFunc",
		Package:   []string{"pkg", "sub_pkg"},
		Line:      3,
		Length:    5,
		File:      "src/file1.cpp",
		Compexity: 6,
	})
}

func TestReadChurn(t *testing.T) {
	jsonData := `{
        "files": [
            {
                "path": "src/file1.cpp",
                "changes": 150,
                "additions": 100,
                "deletions": 50,
                "commits": 10
            },
            {
                "path": "path/to/src/file2.cpp",
                "changes": 75,
                "additions": 45,
                "deletions": 30,
                "commits": 5
            }
        ]
    }`

	reader := strings.NewReader(jsonData)
	got, err := ReadChurn(reader)

	assert.NoError(t, err)
	assert.Len(t, got, 2)

	assert.Contains(t, got, &ChurnChunk{
		File:    "src/file1.cpp",
		Churn:   150,
		Added:   100,
		Removed: 50,
		Commits: 10,
	})

	assert.Contains(t, got, &ChurnChunk{
		File:    "path/to/src/file2.cpp",
		Churn:   75,
		Added:   45,
		Removed: 30,
		Commits: 5,
	})
}
