package read

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type FileStat struct {
	Path      string
	Functions FunctionsStat
}

type FilesStat = []*FileStat

type FunctionStat struct {
	File      string
	Package   []string
	Name      string
	Line      uint
	Length    uint
	Compexity uint
}

type FunctionsStat = []FunctionStat

type lizardItem struct {
	Name   string `xml:"name,attr"`
	Values []int  `xml:"value"`
}

type lizardFile struct {
	Items []lizardItem `xml:"item"`
}

type lizardMeasure struct {
	Items []lizardItem `xml:"item"`
	Sums  []struct {
		Label string `xml:"label,attr"`
		Value int    `xml:"value,attr"`
	} `xml:"sum"`
	Type string `xml:"type,attr"`
}

type lizardXML struct {
	XMLName  xml.Name        `xml:"cppncss"`
	Measures []lizardMeasure `xml:"measure"`
}

var re = regexp.MustCompile(`(.*?)\s+at\s+(.*?):(\d+)`)

func parseItem(item lizardItem) (FunctionStat, error) {
	re := regexp.MustCompile(`(.*?)\s+at\s+(.*?):(\d+)`)
	matches := re.FindStringSubmatch(item.Name)

	if len(matches) != 4 {
		return FunctionStat{}, fmt.Errorf("invalid function format: %s", item.Name)
	}

	funcParts := strings.Split(matches[1], "::")
	name := funcParts[len(funcParts)-1]

	if idx := strings.Index(name, "("); idx != -1 {
		name = name[:idx]
	}

	lineNum, _ := strconv.ParseUint(matches[3], 10, 32)

	return FunctionStat{
		File:      matches[2],
		Name:      name,
		Package:   funcParts[:len(funcParts)-1],
		Line:      uint(lineNum),
		Length:    uint(item.Values[1]),
		Compexity: uint(item.Values[2]),
	}, nil
}

func ReadLizardXML(r io.Reader) (*lizardXML, error) {
	var lizard lizardXML
	if err := xml.NewDecoder(r).Decode(&lizard); err != nil {
		return nil, err
	}

	return &lizard, nil
}

func ParseLizard(lizard *lizardXML) (FilesStat, error) {
	type Filename = string
	fileMap := make(map[Filename]*FileStat)

	filesIdx := slices.IndexFunc(lizard.Measures, func(l lizardMeasure) bool {
		return l.Type == "File"
	})

	for _, file := range lizard.Measures[filesIdx].Items {
		fileMap[file.Name] = new(FileStat)
		fileMap[file.Name].Path = file.Name
	}

	funcIdx := slices.IndexFunc(lizard.Measures, func(l lizardMeasure) bool {
		return l.Type == "Function"
	})

	for _, function := range lizard.Measures[funcIdx].Items {
		stat, err := parseItem(function)
		if err != nil {
			return nil, err
		}

		fileMap[stat.File].Functions = append(fileMap[stat.File].Functions, stat)
	}

	// Convert map to slice
	result := make([]*FileStat, 0, len(fileMap))
	for _, fileStat := range fileMap {
		result = append(result, fileStat)
	}

	return result, nil
}

type ChurnChunk struct {
	File    string `json:"path"`
	Churn   uint   `json:"changes"`
	Added   uint   `json:"additions"`
	Removed uint   `json:"deletions"`
	Commits uint   `json:"commits"`
}

type churnJSON struct {
	Files []*ChurnChunk `json:"files"`
}

func ReadChurn(r io.Reader) ([]*ChurnChunk, error) {
	var data churnJSON
	if err := json.NewDecoder(r).Decode(&data); err != nil {
		return nil, err
	}

	return data.Files, nil
}
