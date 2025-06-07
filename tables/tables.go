package tables

import (
	"fmt"
	"strings"

	"tablemaker/output"
)

// TableConfig represents the main configuration structure for ASCII tables
type TableConfig struct {
	Type      string            `json:"type"`
	Name      string            `json:"name"`
	Headers   []string          `json:"headers"`
	Rows      [][]string        `json:"rows"`
	Alignment []string          `json:"alignment,omitempty"`
	PNG       *output.PNGConfig `json:"png,omitempty"`
}

// AlignmentType represents text alignment options
type AlignmentType string

const (
	AlignLeft   AlignmentType = "left"
	AlignCenter AlignmentType = "center"
	AlignRight  AlignmentType = "right"
)

// TableStyle defines the characters used for table borders
type TableStyle struct {
	TopLeft     string
	TopRight    string
	BottomLeft  string
	BottomRight string
	Horizontal  string
	Vertical    string
	TopJoin     string
	BottomJoin  string
	LeftJoin    string
	RightJoin   string
	Cross       string
}

// Predefined table styles
var tableStyles = map[string]TableStyle{
	"single-line-full": {
		TopLeft:     "┌",
		TopRight:    "┐",
		BottomLeft:  "└",
		BottomRight: "┘",
		Horizontal:  "─",
		Vertical:    "│",
		TopJoin:     "┬",
		BottomJoin:  "┴",
		LeftJoin:    "├",
		RightJoin:   "┤",
		Cross:       "┼",
	},
	"double-line-full": {
		TopLeft:     "╔",
		TopRight:    "╗",
		BottomLeft:  "╚",
		BottomRight: "╝",
		Horizontal:  "═",
		Vertical:    "║",
		TopJoin:     "╦",
		BottomJoin:  "╩",
		LeftJoin:    "╠",
		RightJoin:   "╣",
		Cross:       "╬",
	},
}

// TableRenderer interface for different table styles
type TableRenderer interface {
	Render(config TableConfig) string
}

// ASCIITableRenderer renders tables with configurable ASCII styles
type ASCIITableRenderer struct {
	Style TableStyle
}

func parseAlignment(align string) AlignmentType {
	switch strings.ToLower(align) {
	case "center", "centre":
		return AlignCenter
	case "right":
		return AlignRight
	default:
		return AlignLeft
	}
}

func cleanText(text string) string {
	return strings.ReplaceAll(text, "**", "")
}

func getDisplayLength(text string) int {
	// Use rune length for proper unicode support
	return len([]rune(cleanText(text)))
}

func calculateColumnWidths(config TableConfig) []int {
	if len(config.Headers) == 0 {
		return []int{}
	}

	colWidths := make([]int, len(config.Headers))

	// Check header widths using display length
	for i, header := range config.Headers {
		colWidths[i] = getDisplayLength(header)
	}

	// Check all row cell widths using display length
	for _, row := range config.Rows {
		for i, cell := range row {
			if i < len(colWidths) {
				cellLen := getDisplayLength(cell)
				if cellLen > colWidths[i] {
					colWidths[i] = cellLen
				}
			}
		}
	}

	// Add 2 spaces padding (1 left + 1 right)
	for i := range colWidths {
		colWidths[i] += 2
	}

	return colWidths
}

func getColumnAlignment(config TableConfig, columnIndex int) AlignmentType {
	if columnIndex < len(config.Alignment) {
		return parseAlignment(config.Alignment[columnIndex])
	}
	return AlignLeft
}

func formatCellContent(content string, width int, alignment AlignmentType) string {
	cleanContent := cleanText(content)
	contentLen := getDisplayLength(content)

	// Make sure we have enough space
	if width < contentLen+2 {
		width = contentLen + 2
	}

	// Calculate total spaces needed to fill the width
	spacesNeeded := width - contentLen

	switch alignment {
	case AlignCenter:
		leftSpaces := spacesNeeded / 2
		rightSpaces := spacesNeeded - leftSpaces
		return strings.Repeat(" ", leftSpaces) + cleanContent + strings.Repeat(" ", rightSpaces)
	case AlignRight:
		leftSpaces := spacesNeeded - 1
		return strings.Repeat(" ", leftSpaces) + cleanContent + " "
	default: // AlignLeft
		rightSpaces := spacesNeeded - 1
		return " " + cleanContent + strings.Repeat(" ", rightSpaces)
	}
}

func (r *ASCIITableRenderer) Render(config TableConfig) string {
	if len(config.Rows) == 0 || len(config.Headers) == 0 {
		return ""
	}

	colWidths := calculateColumnWidths(config)
	var result strings.Builder

	// Top border
	result.WriteString(r.Style.TopLeft)
	for i, width := range colWidths {
		result.WriteString(strings.Repeat(r.Style.Horizontal, width))
		if i < len(colWidths)-1 {
			result.WriteString(r.Style.TopJoin)
		}
	}
	result.WriteString(r.Style.TopRight + "\n")

	// Header row
	result.WriteString(r.Style.Vertical)
	for i, header := range config.Headers {
		alignment := getColumnAlignment(config, i)
		formattedCell := formatCellContent(header, colWidths[i], alignment)
		result.WriteString(formattedCell + r.Style.Vertical)
	}
	result.WriteString("\n")

	// Header separator
	result.WriteString(r.Style.LeftJoin)
	for i, width := range colWidths {
		result.WriteString(strings.Repeat(r.Style.Horizontal, width))
		if i < len(colWidths)-1 {
			result.WriteString(r.Style.Cross)
		}
	}
	result.WriteString(r.Style.RightJoin + "\n")

	// Data rows
	for rowIdx, row := range config.Rows {
		result.WriteString(r.Style.Vertical)
		for i, cell := range row {
			if i < len(colWidths) {
				alignment := getColumnAlignment(config, i)
				formattedCell := formatCellContent(cell, colWidths[i], alignment)
				result.WriteString(formattedCell + r.Style.Vertical)
			}
		}
		result.WriteString("\n")

		// Row separator (except for last row)
		if rowIdx < len(config.Rows)-1 {
			result.WriteString(r.Style.LeftJoin)
			for i, width := range colWidths {
				result.WriteString(strings.Repeat(r.Style.Horizontal, width))
				if i < len(colWidths)-1 {
					result.WriteString(r.Style.Cross)
				}
			}
			result.WriteString(r.Style.RightJoin + "\n")
		}
	}

	// Bottom border
	result.WriteString(r.Style.BottomLeft)
	for i, width := range colWidths {
		result.WriteString(strings.Repeat(r.Style.Horizontal, width))
		if i < len(colWidths)-1 {
			result.WriteString(r.Style.BottomJoin)
		}
	}
	result.WriteString(r.Style.BottomRight + "\n")

	return result.String()
}

func GetRenderer(tableType string) (TableRenderer, error) {
	style, exists := tableStyles[strings.ToLower(tableType)]
	if !exists {
		return nil, fmt.Errorf("unknown table type: %s. Available types: %v",
			tableType, getAvailableTypes())
	}
	return &ASCIITableRenderer{Style: style}, nil
}

func GetAvailableTypes() []string {
	return getAvailableTypes()
}

func getAvailableTypes() []string {
	types := make([]string, 0, len(tableStyles))
	for typeName := range tableStyles {
		types = append(types, typeName)
	}
	return types
}

func RegisterTableStyle(name string, style TableStyle) {
	tableStyles[strings.ToLower(name)] = style
}
