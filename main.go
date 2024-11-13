package main

import (
	"fmt"
	"github.com/baranov-V-V/ccv/pkg/plot"
	"github.com/baranov-V-V/ccv/pkg/process"
	"github.com/baranov-V-V/ccv/pkg/read"
	"os"

	"github.com/spf13/cobra"
)

// File to store the output graph
var outputFile = ""

func main() {
	rootCmd := &cobra.Command{
		Use:   "ccv [flags] churn_file complexity_file",
		Short: "Compare code complexity and churn metrics",
		Args:  cobra.ExactArgs(2),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return plot.ValidateRiskThresholds()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			churnFile := args[0]
			complexityFile := args[1]

			if process.Verbose {
				fmt.Printf("Processing files:\n  Churn: %s\n  Complexity: %s\n", churnFile, complexityFile)
			}

			// Read churn data
			cf, err := os.Open(churnFile)
			if err != nil {
				return fmt.Errorf("error opening churn file: %w", err)
			}
			defer cf.Close()

			churns, err := read.ReadChurn(cf)
			if err != nil {
				return fmt.Errorf("error reading churn data: %w", err)
			}

			// Read complexity data
			xf, err := os.Open(complexityFile)
			if err != nil {
				return fmt.Errorf("Error opening complexity file: %w\n", err)
			}
			defer xf.Close()

			lizard, err := read.ReadLizardXML(xf)
			if err != nil {
				return fmt.Errorf("Error reading complexity data: %w\n", err)
			}

			files, err := read.ParseLizard(lizard)
			if err != nil {
				return fmt.Errorf("Error parsing complexity data: %w\n", err)
			}

			// Prepare plot data
			files = process.ApplyFilters(files, process.ComplexityFilter{MinComplexity: 5}.Filter)
			entries := process.PreparePlotData(files, churns)

			// Generate plot
			if err := plot.CreateComplexityChurnChart(entries, outputFile); err != nil {
				return fmt.Errorf("error creating chart: %w\n", err)
			}

			if process.Verbose {
				fmt.Printf("Chart generated: %s\n", outputFile)
			}

			return nil
		},
	}

	flags := rootCmd.PersistentFlags()
	flags.StringVarP(&outputFile, "output", "o", "complexity_churn.html", "Output file path")
	flags.BoolVarP(&process.Verbose, "verbose", "v", false, "Enable verbose output")
	flags.StringVarP(&process.Plot, "plot-type", "t", "changes", "Specify OY plot type")
	flags.UintVar(&plot.VeryLowRisk, "very-low-risk", 10, "Very Low Risk threshold")
	flags.UintVar(&plot.LowRisk, "low-risk", 15, "Low Risk threshold")
	flags.UintVar(&plot.MediumRisk, "medium-risk", 20, "Medium Risk threshold")
	flags.UintVar(&plot.HighRisk, "high-risk", 25, "High Risk threshold")
	flags.UintVar(&plot.VeryHighRisk, "very-high-risk", 30, "Very High Risk threshold")
	flags.UintVar(&plot.CriticalRisk, "critical-risk", 35, "Critical Risk threshold")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
