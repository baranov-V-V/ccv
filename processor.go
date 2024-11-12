package main

import (
	"fmt"
)

type ComplexityInputType int

const (
	Lizard ComplexityInputType = iota
	ClangTidy
)

type ChurnInputType int

const (
	ModifiedScript ChurnInputType = iota
)

var Version = "0.0.1"

var Verbose = false

var ComplexityInput ComplexityInputType = Lizard

var ChurnInput ChurnInputType = ModifiedScript

type PlotType = string

const (
	Commits PlotType = "commits"
	Changes PlotType = "changes"
)

var Plot = Commits

// type FilesFilter func(files FilesStat) FilesStat

// Place where it is used
type FilesFilter interface {
	Filter(files FilesStat) FilesStat
}

type ComplexityFilter struct {
	minComlexity uint
}

func (f ComplexityFilter) Filter(files FilesStat) FilesStat {
	result := make(FilesStat, 0, len(files))

	for _, file := range files {
		filteredFuncs := make([]FunctionStat, 0)
		for _, fn := range file.Functions {
			if fn.Compexity >= f.minComlexity {
				filteredFuncs = append(filteredFuncs, fn)
			}
		}

		if len(filteredFuncs) > 0 {
			newFile := &FileStat{
				Path:      file.Path,
				Functions: filteredFuncs,
			}
			result = append(result, newFile)
		}
	}

	return result
}

type FilesFilterFunc func(files FilesStat) FilesStat

func ApplyFilters(files FilesStat, filters ...FilesFilterFunc) FilesStat {
	result := files

	for _, filter := range filters {
		result = filter(result)
	}

	return result
}

type FileComplexity struct {
	File       string
	Complexity float64
}

// Calculates average complexity bases on functions in file: sum(funcComplexity) / funcCount
func avgComplexity(files FilesStat) []FileComplexity {
	result := make([]FileComplexity, 0, len(files))

	for _, file := range files {
		if len(file.Functions) == 0 {
			continue
		}

		var totalComplexity float64 = 0
		for _, fn := range file.Functions {
			totalComplexity += float64(fn.Compexity)
		}

		complexity := totalComplexity / float64(len(file.Functions))
		if Verbose {
			fmt.Printf("File: %s, Complexity: %f\n", file.Path, complexity)
		}

		result = append(result, FileComplexity{
			File:       file.Path,
			Complexity: complexity,
		})
	}

	return result
}

// Skip file if it is not found in chunk or files, first goes over all churns
// Matches based on filename
func PreparePlotData(files FilesStat, churns []*ChurnChunk) []ChartEntry {
	result := make([]ChartEntry, 0)

	// Calculate average complexity for each file
	fileComplexities := avgComplexity(files)

	// Create map for quick churn lookup
	churnMap := make(map[string]*ChurnChunk)
	for _, churn := range churns {
		churnMap[churn.File] = churn
	}

	// Match files with churns and create chart entries
	for _, fc := range fileComplexities {
		churn, exists := churnMap[fc.File]

		if !exists {
			continue
		}

		entry := ChartEntry{
			File:       fc.File,
			Complexity: fc.Complexity,
		}

		switch Plot {
		case Commits:
			entry.Churn = churn.Commits
		case Changes:
			entry.Churn = churn.Churn
		default:
			panic("Unknown plot type")
		}

		result = append(result, entry)
	}

	return result
}
