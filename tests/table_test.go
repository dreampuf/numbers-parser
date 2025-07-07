package tests

import (
	"testing"

	"github.com/numbers-parser-go/pkg/numbers"
)

func TestCellAccess(t *testing.T) {
	doc := numbers.NewDocument()
	defer doc.Close()

	table := doc.DefaultTable()

	// Write test data
	table.Write(1, 1, "Test Value")

	// Test access by coordinates
	cell, err := table.Cell(1, 1)
	if err != nil {
		t.Errorf("Error accessing cell by coordinates: %v", err)
	}
	if cell.Value() != "Test Value" {
		t.Errorf("Expected 'Test Value', got %v", cell.Value())
	}

	// Test access by Excel notation
	cell, err = table.Cell("B2")
	if err != nil {
		t.Errorf("Error accessing cell by Excel notation: %v", err)
	}
	if cell.Value() != "Test Value" {
		t.Errorf("Expected 'Test Value', got %v", cell.Value())
	}
}

func TestCellAccessOutOfBounds(t *testing.T) {
	doc := numbers.NewDocument()
	defer doc.Close()

	table := doc.DefaultTable()

	// Test out of bounds access
	_, err := table.Cell(1000, 1000)
	if err == nil {
		t.Error("Expected error for out of bounds access")
	}
}

func TestWriteExcelNotation(t *testing.T) {
	doc := numbers.NewDocument()
	defer doc.Close()

	table := doc.DefaultTable()

	// Write using Excel notation
	err := table.Write("A1", "Header")
	if err != nil {
		t.Errorf("Error writing with Excel notation: %v", err)
	}

	// Verify the value
	cell, _ := table.Cell(0, 0)
	if cell.Value() != "Header" {
		t.Errorf("Expected 'Header', got %v", cell.Value())
	}
}

func TestWriteInvalidArgs(t *testing.T) {
	doc := numbers.NewDocument()
	defer doc.Close()

	table := doc.DefaultTable()

	// Test insufficient arguments
	err := table.Write(0)
	if err == nil {
		t.Error("Expected error for insufficient arguments")
	}

	err = table.Write(0, 1)
	if err == nil {
		t.Error("Expected error for insufficient arguments")
	}
}

func TestTableExpansion(t *testing.T) {
	doc := numbers.NewDocument()
	defer doc.Close()

	table := doc.DefaultTable()
	initialRows := table.NumRows
	initialCols := table.NumCols

	// Write beyond current table bounds
	err := table.Write(20, 15, "Expansion Test")
	if err != nil {
		t.Errorf("Error writing to expanded position: %v", err)
	}

	// Check that table expanded
	if table.NumRows <= initialRows {
		t.Error("Table should have expanded rows")
	}

	if table.NumCols <= initialCols {
		t.Error("Table should have expanded columns")
	}

	// Verify the value is accessible
	cell, err := table.Cell(20, 15)
	if err != nil {
		t.Errorf("Error accessing expanded cell: %v", err)
	}
	if cell.Value() != "Expansion Test" {
		t.Errorf("Expected 'Expansion Test', got %v", cell.Value())
	}
}

func TestAddRow(t *testing.T) {
	doc := numbers.NewDocument()
	defer doc.Close()

	table := doc.DefaultTable()
	initialRows := table.NumRows

	// Add one row
	err := table.AddRow(1)
	if err != nil {
		t.Errorf("Error adding row: %v", err)
	}

	if table.NumRows != initialRows+1 {
		t.Errorf("Expected %d rows, got %d", initialRows+1, table.NumRows)
	}

	// Add multiple rows
	err = table.AddRow(5)
	if err != nil {
		t.Errorf("Error adding multiple rows: %v", err)
	}

	if table.NumRows != initialRows+6 {
		t.Errorf("Expected %d rows, got %d", initialRows+6, table.NumRows)
	}
}

func TestAddColumn(t *testing.T) {
	doc := numbers.NewDocument()
	defer doc.Close()

	table := doc.DefaultTable()
	initialCols := table.NumCols

	// Add one column
	err := table.AddColumn(1)
	if err != nil {
		t.Errorf("Error adding column: %v", err)
	}

	if table.NumCols != initialCols+1 {
		t.Errorf("Expected %d columns, got %d", initialCols+1, table.NumCols)
	}

	// Add multiple columns
	err = table.AddColumn(3)
	if err != nil {
		t.Errorf("Error adding multiple columns: %v", err)
	}

	if table.NumCols != initialCols+4 {
		t.Errorf("Expected %d columns, got %d", initialCols+4, table.NumCols)
	}
}

func TestAddRowAtPosition(t *testing.T) {
	doc := numbers.NewDocument()
	defer doc.Close()

	table := doc.DefaultTable()

	// Fill some cells with test data
	table.Write(0, 0, "Row 0")
	table.Write(1, 0, "Row 1")
	table.Write(2, 0, "Row 2")

	// Insert row at position 1
	err := table.AddRow(1, 1)
	if err != nil {
		t.Errorf("Error adding row at position: %v", err)
	}

	// Check that values shifted correctly
	cell, _ := table.Cell(0, 0)
	if cell.Value() != "Row 0" {
		t.Error("Row 0 should remain in place")
	}

	cell, _ = table.Cell(2, 0)
	if cell.Value() != "Row 1" {
		t.Error("Row 1 should have shifted to position 2")
	}

	cell, _ = table.Cell(3, 0)
	if cell.Value() != "Row 2" {
		t.Error("Row 2 should have shifted to position 3")
	}

	// Check that inserted row is empty
	cell, _ = table.Cell(1, 0)
	if !cell.IsEmpty() {
		t.Error("Inserted row should be empty")
	}
}

func TestRowsValues(t *testing.T) {
	doc := numbers.NewDocument()
	defer doc.Close()

	table := doc.DefaultTable()

	// Write test data
	testValues := [][]interface{}{
		{"A", 1, true},
		{"B", 2, false},
		{"C", 3, true},
	}

	for row, rowData := range testValues {
		for col, value := range rowData {
			table.Write(row, col, value)
		}
	}

	// Get values only
	values, err := table.RowsValues()
	if err != nil {
		t.Errorf("Error getting row values: %v", err)
	}

	// Check first few rows match our test data
	for row := 0; row < len(testValues); row++ {
		for col := 0; col < len(testValues[row]); col++ {
			if values[row][col] != testValues[row][col] {
				t.Errorf("Value mismatch at (%d,%d): expected %v, got %v",
					row, col, testValues[row][col], values[row][col])
			}
		}
	}
}

func TestCellTypes(t *testing.T) {
	doc := numbers.NewDocument()
	defer doc.Close()

	table := doc.DefaultTable()

	// Test empty cell
	cell, _ := table.Cell(0, 0)
	if !cell.IsEmpty() {
		t.Error("New cell should be empty")
	}

	// Test text cell
	table.Write(0, 0, "Hello")
	cell, _ = table.Cell(0, 0)
	if !cell.IsText() {
		t.Error("Cell should be text type")
	}

	// Test number cell
	table.Write(0, 1, 42)
	cell, _ = table.Cell(0, 1)
	if !cell.IsNumber() {
		t.Error("Cell should be number type")
	}

	// Test boolean cell
	table.Write(0, 2, true)
	cell, _ = table.Cell(0, 2)
	if !cell.IsBool() {
		t.Error("Cell should be boolean type")
	}
}