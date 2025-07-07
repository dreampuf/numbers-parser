package numbers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ExcelColumnToNumber converts Excel column letters to column number (0-based)
func ExcelColumnToNumber(col string) (int, error) {
	col = strings.ToUpper(col)
	result := 0
	
	for _, char := range col {
		if char < 'A' || char > 'Z' {
			return 0, fmt.Errorf("invalid column letter: %c", char)
		}
		result = result*26 + int(char-'A') + 1
	}
	
	return result - 1, nil // Convert to 0-based
}

// NumberToExcelColumn converts column number (0-based) to Excel column letters
func NumberToExcelColumn(col int) string {
	if col < 0 {
		return ""
	}
	
	result := ""
	col++ // Convert to 1-based for calculation
	
	for col > 0 {
		col--
		result = string(rune('A'+col%26)) + result
		col = col / 26
	}
	
	return result
}

// ExcelAddressToCoords converts Excel address like "B5" to row, col coordinates (0-based)
func ExcelAddressToCoords(addr string) (int, int, error) {
	re := regexp.MustCompile(`^([A-Z]+)(\d+)$`)
	matches := re.FindStringSubmatch(strings.ToUpper(addr))
	if len(matches) != 3 {
		return 0, 0, fmt.Errorf("invalid Excel address: %s", addr)
	}
	
	// Parse column
	col, err := ExcelColumnToNumber(matches[1])
	if err != nil {
		return 0, 0, err
	}
	
	// Parse row (1-based to 0-based)
	row, err := strconv.Atoi(matches[2])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid row number in address %s: %w", addr, err)
	}
	row-- // Convert to 0-based
	
	return row, col, nil
}

// CoordsToExcelAddress converts row, col coordinates (0-based) to Excel address
func CoordsToExcelAddress(row, col int) string {
	if row < 0 || col < 0 {
		return ""
	}
	return NumberToExcelColumn(col) + strconv.Itoa(row+1)
}

// IsValidExcelAddress checks if a string is a valid Excel cell address
func IsValidExcelAddress(addr string) bool {
	re := regexp.MustCompile(`^[A-Z]+[1-9]\d*$`)
	return re.MatchString(strings.ToUpper(addr))
}

// ParseRange parses an Excel range like "A1:C3" into start and end coordinates
func ParseRange(rangeStr string) (startRow, startCol, endRow, endCol int, err error) {
	parts := strings.Split(rangeStr, ":")
	if len(parts) != 2 {
		return 0, 0, 0, 0, fmt.Errorf("invalid range format: %s", rangeStr)
	}
	
	startRow, startCol, err = ExcelAddressToCoords(parts[0])
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid start address in range: %w", err)
	}
	
	endRow, endCol, err = ExcelAddressToCoords(parts[1])
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid end address in range: %w", err)
	}
	
	// Ensure start <= end
	if startRow > endRow {
		startRow, endRow = endRow, startRow
	}
	if startCol > endCol {
		startCol, endCol = endCol, startCol
	}
	
	return startRow, startCol, endRow, endCol, nil
}

// FormatCellValue formats a cell value for display
func FormatCellValue(value interface{}) string {
	if value == nil {
		return ""
	}
	
	switch v := value.(type) {
	case string:
		return v
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%.2f", v)
	case bool:
		if v {
			return "TRUE"
		}
		return "FALSE"
	default:
		return fmt.Sprintf("%v", v)
	}
}

// GetCellRange returns all cells in a specified range
func (t *Table) GetCellRange(rangeStr string) ([]*Cell, error) {
	startRow, startCol, endRow, endCol, err := ParseRange(rangeStr)
	if err != nil {
		return nil, err
	}
	
	var cells []*Cell
	for row := startRow; row <= endRow; row++ {
		for col := startCol; col <= endCol; col++ {
			if row < t.NumRows && col < t.NumCols {
				cells = append(cells, t.cells[row][col])
			}
		}
	}
	
	return cells, nil
}

// SetCellRange sets values for all cells in a range
func (t *Table) SetCellRange(rangeStr string, values [][]interface{}) error {
	startRow, startCol, endRow, endCol, err := ParseRange(rangeStr)
	if err != nil {
		return err
	}
	
	for i, row := range values {
		targetRow := startRow + i
		if targetRow > endRow {
			break
		}
		
		for j, value := range row {
			targetCol := startCol + j
			if targetCol > endCol {
				break
			}
			
			if err := t.Write(targetRow, targetCol, value); err != nil {
				return err
			}
		}
	}
	
	return nil
}

// GetColumnWidth returns the width of a column (placeholder implementation)
func (t *Table) GetColumnWidth(col int) float64 {
	// This would need to be implemented based on the actual Numbers format
	return 100.0 // Default width
}

// SetColumnWidth sets the width of a column (placeholder implementation)
func (t *Table) SetColumnWidth(col int, width float64) error {
	// This would need to be implemented based on the actual Numbers format
	return nil
}

// GetRowHeight returns the height of a row (placeholder implementation)
func (t *Table) GetRowHeight(row int) float64 {
	// This would need to be implemented based on the actual Numbers format
	return 20.0 // Default height
}

// SetRowHeight sets the height of a row (placeholder implementation)
func (t *Table) SetRowHeight(row int, height float64) error {
	// This would need to be implemented based on the actual Numbers format
	return nil
}