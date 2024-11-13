package plot

import (
	"fmt"

	"github.com/go-echarts/go-echarts/v2/opts"
)

type RiskLevel struct {
	Name  string
	Color string
	Min   uint
	Max   uint
}

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

// deprecated TODO: delete!
func getRiskColors(levels []RiskLevel) []string {
	colors := make([]string, len(levels))
	for i, level := range levels {
		colors[i] = level.Color
	}
	return colors
}

type RisksMapper struct {
	levels []RiskLevel
}

func NewRisksMapper() *RisksMapper {
	return &RisksMapper{
		levels: getRiskLevels(),
	}
}

var _ EntryMapper = (*RisksMapper)(nil)
	func (rm *RisksMapper) Map(data ScatterData) Category {
		riskScore := data.Complexity + float64(data.Churn)
	
		for _, level := range rm.levels {
			if riskScore >= float64(level.Min) && riskScore <= float64(level.Max) {
				return level.Name
			}
		}
	
		return "Unknown"
	}

	func (rm *RisksMapper) Style(category Category) opts.ItemStyle {
		for _, level := range rm.levels {
			if level.Name == category {
				return opts.ItemStyle{
					Color: level.Color,
				}
			}
		}
	
		return opts.ItemStyle{}
	}
