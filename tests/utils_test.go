package tests

import (
	"testing"

	"github.com/numbers-parser-go/pkg/numbers"
)

func TestExcelColumnToNumber(t *testing.T) {
	tests := []struct {
		col      string
		expected int
	}{
		{"A", 0},
		{"B", 1},
		{"Z", 25},
		{"AA", 26},
		{"AB", 27},
		{"AZ", 51},
		{"BA", 52},
		{"ZZ", 701},
	}

	for _, test := range tests {
		result, err := numbers.ExcelColumnToNumber(test.col)
		if err != nil {
			t.Errorf("Error converting column %s: %v", test.col, err)
		}
		if result != test.expected {
			t.Errorf("ExcelColumnToNumber(%s): expected %d, got %d", test.col, test.expected, result)
		}
	}
}

func TestNumberToExcelColumn(t *testing.T) {
	tests := []struct {
		col      int
		expected string
	}{
		{0, "A"},
		{1, "B"},
		{25, "Z"},
		{26, "AA"},
		{27, "AB"},
		{51, "AZ"},
		{52, "BA"},
		{701, "ZZ"},
	}

	for _, test := range tests {
		result := numbers.NumberToExcelColumn(test.col)
		if result != test.expected {
			t.Errorf("NumberToExcelColumn(%d): expected %s, got %s", test.col, test.expected, result)
		}
	}
}

func TestExcelAddressToCoords(t *testing.T) {
	tests := []struct {
		addr        string
		expectedRow int
		expectedCol int
	}{
		{"A1", 0, 0},
		{"B2", 1, 1},
		{"Z1", 0, 25},
		{"AA1", 0, 26},
		{"A10", 9, 0},
		{"B100", 99, 1},
	}

	for _, test := range tests {
		row, col, err := numbers.ExcelAddressToCoords(test.addr)
		if err != nil {
			t.Errorf("Error parsing address %s: %v", test.addr, err)
		}
		if row != test.expectedRow || col != test.expectedCol {
			t.Errorf("ExcelAddressToCoords(%s): expected (%d,%d), got (%d,%d)",
				test.addr, test.expectedRow, test.expectedCol, row, col)
		}
	}
}

func TestCoordsToExcelAddress(t *testing.T) {
	tests := []struct {
		row      int
		col      int
		expected string
	}{
		{0, 0, "A1"},
		{1, 1, "B2"},
		{0, 25, "Z1"},
		{0, 26, "AA1"},
		{9, 0, "A10"},
		{99, 1, "B100"},
	}

	for _, test := range tests {
		result := numbers.CoordsToExcelAddress(test.row, test.col)
		if result != test.expected {
			t.Errorf("CoordsToExcelAddress(%d,%d): expected %s, got %s",
				test.row, test.col, test.expected, result)
		}
	}
}

func TestIsValidExcelAddress(t *testing.T) {
	validAddresses := []string{"A1", "B2", "Z1", "AA1", "ABC123", "XFD1048576"}
	invalidAddresses := []string{"", "A", "1", "A0", "1A", "A1B", "a1", "A-1"}

	for _, addr := range validAddresses {
		if !numbers.IsValidExcelAddress(addr) {
			t.Errorf("Address %s should be valid", addr)
		}
	}

	for _, addr := range invalidAddresses {
		if numbers.IsValidExcelAddress(addr) {
			t.Errorf("Address %s should be invalid", addr)
		}
	}
}

func TestParseRange(t *testing.T) {
	tests := []struct {
		rangeStr   string
		startRow   int
		startCol   int
		endRow     int
		endCol     int
		shouldFail bool
	}{
		{"A1:C3", 0, 0, 2, 2, false},
		{"B2:D4", 1, 1, 3, 3, false},
		{"A1:A1", 0, 0, 0, 0, false},
		{"C3:A1", 0, 0, 2, 2, false}, // Should swap to make start <= end
		{"A1", 0, 0, 0, 0, true},     // Invalid format
		{"A1:B2:C3", 0, 0, 0, 0, true}, // Invalid format
	}

	for _, test := range tests {
		startRow, startCol, endRow, endCol, err := numbers.ParseRange(test.rangeStr)
		if test.shouldFail {
			if err == nil {
				t.Errorf("ParseRange(%s): expected error but got none", test.rangeStr)
			}
		} else {
			if err != nil {
				t.Errorf("ParseRange(%s): unexpected error: %v", test.rangeStr, err)
			}
			if startRow != test.startRow || startCol != test.startCol ||
				endRow != test.endRow || endCol != test.endCol {
				t.Errorf("ParseRange(%s): expected (%d,%d:%d,%d), got (%d,%d:%d,%d)",
					test.rangeStr, test.startRow, test.startCol, test.endRow, test.endCol,
					startRow, startCol, endRow, endCol)
			}
		}
	}
}

func TestFormatCellValue(t *testing.T) {
	tests := []struct {
		value    interface{}
		expected string
	}{
		{nil, ""},
		{"Hello", "Hello"},
		{42, "42"},
		{3.14, "3.14"},
		{true, "TRUE"},
		{false, "FALSE"},
	}

	for _, test := range tests {
		result := numbers.FormatCellValue(test.value)
		if result != test.expected {
			t.Errorf("FormatCellValue(%v): expected %s, got %s", test.value, test.expected, result)
		}
	}
}

func TestGetCellRange(t *testing.T) {
	doc := numbers.NewDocument()
	defer doc.Close()

	table := doc.DefaultTable()

	// Fill a range with test data
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			table.Write(row, col, row*10+col)
		}
	}

	// Get range A1:C3
	cells, err := table.GetCellRange("A1:C3")
	if err != nil {
		t.Errorf("Error getting cell range: %v", err)
	}

	if len(cells) != 9 {
		t.Errorf("Expected 9 cells in range, got %d", len(cells))
	}

	// Check a few values
	expectedValues := []int{0, 1, 2, 10, 11, 12, 20, 21, 22}
	for i, cell := range cells {
		if cell.Value() != expectedValues[i] {
			t.Errorf("Cell %d: expected %d, got %v", i, expectedValues[i], cell.Value())
		}
	}
}

func TestSetCellRange(t *testing.T) {
	doc := numbers.NewDocument()
	defer doc.Close()

	table := doc.DefaultTable()

	// Set values for a range
	values := [][]interface{}{
		{"A", "B", "C"},
		{1, 2, 3},
		{true, false, true},
	}

	err := table.SetCellRange("A1:C3", values)
	if err != nil {
		t.Errorf("Error setting cell range: %v", err)
	}

	// Verify the values were set correctly
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			cell, _ := table.Cell(row, col)
			expected := values[row][col]
			actual := cell.Value()
			if actual != expected {
				t.Errorf("Cell (%d,%d): expected %v, got %v", row, col, expected, actual)
			}
		}
	}
}