package tests

import (
	"testing"
	"time"

	"github.com/numbers-parser-go/pkg/numbers"
)

func TestNewDocument(t *testing.T) {
	doc := numbers.NewDocument()
	defer doc.Close()

	if len(doc.Sheets) == 0 {
		t.Error("NewDocument should create at least one sheet")
	}

	if doc.Sheets[0].Name != "Sheet 1" {
		t.Errorf("Expected first sheet name to be 'Sheet 1', got '%s'", doc.Sheets[0].Name)
	}

	if len(doc.Sheets[0].Tables) == 0 {
		t.Error("Default sheet should have at least one table")
	}

	table := doc.Sheets[0].Tables[0]
	if table.Name != "Table 1" {
		t.Errorf("Expected first table name to be 'Table 1', got '%s'", table.Name)
	}

	if table.NumRows != 12 || table.NumCols != 8 {
		t.Errorf("Expected table dimensions 12x8, got %dx%d", table.NumRows, table.NumCols)
	}
}

func TestAddSheet(t *testing.T) {
	doc := numbers.NewDocument()
	defer doc.Close()

	initialSheetCount := len(doc.Sheets)

	sheet := doc.AddSheet("Test Sheet")
	if sheet == nil {
		t.Error("AddSheet returned nil")
	}

	if len(doc.Sheets) != initialSheetCount+1 {
		t.Errorf("Expected %d sheets, got %d", initialSheetCount+1, len(doc.Sheets))
	}

	if sheet.Name != "Test Sheet" {
		t.Errorf("Expected sheet name 'Test Sheet', got '%s'", sheet.Name)
	}
}

func TestAddSheetWithEmptyName(t *testing.T) {
	doc := numbers.NewDocument()
	defer doc.Close()

	sheet := doc.AddSheet("")
	if sheet.Name != "Sheet 2" {
		t.Errorf("Expected auto-generated name 'Sheet 2', got '%s'", sheet.Name)
	}
}

func TestGetSheet(t *testing.T) {
	doc := numbers.NewDocument()
	defer doc.Close()

	doc.AddSheet("Test Sheet")

	// Test get by index
	sheet := doc.GetSheet(0)
	if sheet == nil {
		t.Error("GetSheet by index returned nil")
	}
	if sheet.Name != "Sheet 1" {
		t.Errorf("Expected 'Sheet 1', got '%s'", sheet.Name)
	}

	// Test get by name
	sheet = doc.GetSheet("Test Sheet")
	if sheet == nil {
		t.Error("GetSheet by name returned nil")
	}
	if sheet.Name != "Test Sheet" {
		t.Errorf("Expected 'Test Sheet', got '%s'", sheet.Name)
	}

	// Test invalid index
	sheet = doc.GetSheet(999)
	if sheet != nil {
		t.Error("GetSheet with invalid index should return nil")
	}

	// Test invalid name
	sheet = doc.GetSheet("Nonexistent Sheet")
	if sheet != nil {
		t.Error("GetSheet with invalid name should return nil")
	}
}

func TestDefaultTable(t *testing.T) {
	doc := numbers.NewDocument()
	defer doc.Close()

	table := doc.DefaultTable()
	if table == nil {
		t.Error("DefaultTable returned nil")
	}

	if table.Name != "Table 1" {
		t.Errorf("Expected default table name 'Table 1', got '%s'", table.Name)
	}
}

func TestAddStyle(t *testing.T) {
	doc := numbers.NewDocument()
	defer doc.Close()

	style := numbers.NewStyle()
	style.Name = "Test Style"
	style.Bold = true
	style.FontSize = 14.0

	doc.AddStyle(style)

	retrievedStyle, exists := doc.GetStyle("Test Style")
	if !exists {
		t.Error("Style was not added successfully")
	}

	if !retrievedStyle.Bold {
		t.Error("Style properties were not preserved")
	}

	if retrievedStyle.FontSize != 14.0 {
		t.Errorf("Expected font size 14.0, got %f", retrievedStyle.FontSize)
	}
}

func TestAddStyleWithoutName(t *testing.T) {
	doc := numbers.NewDocument()
	defer doc.Close()

	style := numbers.NewStyle()
	// Don't set name

	doc.AddStyle(style)

	if style.Name == "" {
		t.Error("Style should have been assigned an auto-generated name")
	}

	if style.Name != "Custom Style 1" {
		t.Errorf("Expected auto-generated name 'Custom Style 1', got '%s'", style.Name)
	}
}

func TestDocumentWriteAndRead(t *testing.T) {
	doc := numbers.NewDocument()
	defer doc.Close()

	table := doc.DefaultTable()

	// Write test data
	testData := [][]interface{}{
		{"Name", "Age", "Active"},
		{"John", 30, true},
		{"Jane", 25, false},
		{"Bob", 35, true},
	}

	for row, rowData := range testData {
		for col, value := range rowData {
			err := table.Write(row, col, value)
			if err != nil {
				t.Errorf("Error writing to cell (%d,%d): %v", row, col, err)
			}
		}
	}

	// Read and verify data
	rows, err := table.Rows()
	if err != nil {
		t.Errorf("Error reading rows: %v", err)
	}

	if len(rows) < len(testData) {
		t.Errorf("Expected at least %d rows, got %d", len(testData), len(rows))
	}

	for row := 0; row < len(testData); row++ {
		for col := 0; col < len(testData[row]); col++ {
			cell := rows[row][col]
			expected := testData[row][col]
			actual := cell.Value()

			if actual != expected {
				t.Errorf("Cell (%d,%d): expected %v, got %v", row, col, expected, actual)
			}
		}
	}
}

func TestDocumentWithDifferentTypes(t *testing.T) {
	doc := numbers.NewDocument()
	defer doc.Close()

	table := doc.DefaultTable()

	// Test different data types
	now := time.Now()
	duration := 2 * time.Hour

	err := table.Write(0, 0, "Text")
	if err != nil {
		t.Errorf("Error writing text: %v", err)
	}

	err = table.Write(0, 1, 42)
	if err != nil {
		t.Errorf("Error writing int: %v", err)
	}

	err = table.Write(0, 2, 3.14)
	if err != nil {
		t.Errorf("Error writing float: %v", err)
	}

	err = table.Write(0, 3, true)
	if err != nil {
		t.Errorf("Error writing bool: %v", err)
	}

	err = table.Write(0, 4, now)
	if err != nil {
		t.Errorf("Error writing time: %v", err)
	}

	err = table.Write(0, 5, duration)
	if err != nil {
		t.Errorf("Error writing duration: %v", err)
	}

	// Verify types
	cell, _ := table.Cell(0, 0)
	if !cell.IsText() {
		t.Error("Expected text cell")
	}

	cell, _ = table.Cell(0, 1)
	if !cell.IsNumber() {
		t.Error("Expected number cell")
	}

	cell, _ = table.Cell(0, 2)
	if !cell.IsNumber() {
		t.Error("Expected number cell for float")
	}

	cell, _ = table.Cell(0, 3)
	if !cell.IsBool() {
		t.Error("Expected bool cell")
	}

	cell, _ = table.Cell(0, 4)
	if !cell.IsDate() {
		t.Error("Expected date cell")
	}

	cell, _ = table.Cell(0, 5)
	if !cell.IsDuration() {
		t.Error("Expected duration cell")
	}
}