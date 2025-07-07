package numbers

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// NumbersModel handles the low-level Numbers file format operations
type NumbersModel struct {
	zipReader   *zip.ReadCloser
	zipWriter   *zip.Writer
	filename    string
	isPackage   bool
	metadata    map[string]interface{}
	sheetData   map[string]*SheetData
	tableData   map[string]*TableData
}

// SheetData represents the internal data for a sheet
type SheetData struct {
	ID     string
	Name   string
	Tables []string
}

// TableData represents the internal data for a table
type TableData struct {
	ID      string
	Name    string
	NumRows int
	NumCols int
	Cells   map[string]*CellData
}

// CellData represents the internal data for a cell
type CellData struct {
	Row       int
	Col       int
	Type      CellType
	Value     interface{}
	Formatted string
}

// NewNumbersModel creates a new model for an existing Numbers file
func NewNumbersModel(filename string) (*NumbersModel, error) {
	// Check if it's a package (directory) or single file
	info, err := os.Stat(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	model := &NumbersModel{
		filename:  filename,
		isPackage: info.IsDir(),
		metadata:  make(map[string]interface{}),
		sheetData: make(map[string]*SheetData),
		tableData: make(map[string]*TableData),
	}

	if model.isPackage {
		// Handle package format (directory)
		return model, model.loadPackage()
	} else {
		// Handle single file format
		return model, model.loadArchive()
	}
}

// NewEmptyNumbersModel creates a new empty model
func NewEmptyNumbersModel() *NumbersModel {
	return &NumbersModel{
		metadata:  make(map[string]interface{}),
		sheetData: make(map[string]*SheetData),
		tableData: make(map[string]*TableData),
	}
}

// Close closes the model and releases resources
func (m *NumbersModel) Close() error {
	if m.zipReader != nil {
		return m.zipReader.Close()
	}
	return nil
}

// loadArchive loads a Numbers file in archive format
func (m *NumbersModel) loadArchive() error {
	var err error
	m.zipReader, err = zip.OpenReader(m.filename)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}

	// Parse the Numbers format
	for _, file := range m.zipReader.File {
		if err := m.processZipFile(file); err != nil {
			return fmt.Errorf("failed to process file %s: %w", file.Name, err)
		}
	}

	return nil
}

// loadPackage loads a Numbers file in package format
func (m *NumbersModel) loadPackage() error {
	// Walk through the package directory
	return filepath.Walk(m.filename, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			relPath, _ := filepath.Rel(m.filename, path)
			return m.processPackageFile(relPath, path)
		}
		return nil
	})
}

// processZipFile processes a file from the zip archive
func (m *NumbersModel) processZipFile(file *zip.File) error {
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return err
	}

	return m.processFileData(file.Name, data)
}

// processPackageFile processes a file from the package
func (m *NumbersModel) processPackageFile(relPath, fullPath string) error {
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return err
	}

	return m.processFileData(relPath, data)
}

// processFileData processes file data regardless of source
func (m *NumbersModel) processFileData(filename string, data []byte) error {
	// This is a simplified implementation
	// In a real implementation, you would parse the protobuf data
	// and extract sheet/table/cell information
	
	if strings.Contains(filename, "Index") {
		// Process index files to get document structure
		return m.processIndexFile(data)
	}
	
	if strings.Contains(filename, ".iwa") {
		// Process iWork Archive files
		return m.processIWAFile(filename, data)
	}
	
	return nil
}

// processIndexFile processes index files to extract document structure
func (m *NumbersModel) processIndexFile(data []byte) error {
	// Simplified implementation - in reality this would parse protobuf
	// For now, create some default structure
	
	// Create default sheet if none exists
	if len(m.sheetData) == 0 {
		sheetID := "sheet-1"
		m.sheetData[sheetID] = &SheetData{
			ID:     sheetID,
			Name:   "Sheet 1",
			Tables: []string{"table-1"},
		}
		
		tableID := "table-1"
		m.tableData[tableID] = &TableData{
			ID:      tableID,
			Name:    "Table 1",
			NumRows: 12,
			NumCols: 8,
			Cells:   make(map[string]*CellData),
		}
	}
	
	return nil
}

// processIWAFile processes iWork Archive files
func (m *NumbersModel) processIWAFile(filename string, data []byte) error {
	// Simplified implementation - would parse protobuf in reality
	return nil
}

// GetSheetIDs returns all sheet IDs
func (m *NumbersModel) GetSheetIDs() ([]string, error) {
	ids := make([]string, 0, len(m.sheetData))
	for id := range m.sheetData {
		ids = append(ids, id)
	}
	return ids, nil
}

// GetSheetName returns the name of a sheet
func (m *NumbersModel) GetSheetName(sheetID string) (string, error) {
	sheet, exists := m.sheetData[sheetID]
	if !exists {
		return "", fmt.Errorf("sheet %s not found", sheetID)
	}
	return sheet.Name, nil
}

// GetTableIDs returns all table IDs for a sheet
func (m *NumbersModel) GetTableIDs(sheetID string) ([]string, error) {
	sheet, exists := m.sheetData[sheetID]
	if !exists {
		return nil, fmt.Errorf("sheet %s not found", sheetID)
	}
	return sheet.Tables, nil
}

// GetTableName returns the name of a table
func (m *NumbersModel) GetTableName(tableID string) (string, error) {
	table, exists := m.tableData[tableID]
	if !exists {
		return "", fmt.Errorf("table %s not found", tableID)
	}
	return table.Name, nil
}

// GetTableDimensions returns the dimensions of a table
func (m *NumbersModel) GetTableDimensions(tableID string) (int, int, error) {
	table, exists := m.tableData[tableID]
	if !exists {
		return 0, 0, fmt.Errorf("table %s not found", tableID)
	}
	return table.NumRows, table.NumCols, nil
}

// GetCell returns a cell from a table
func (m *NumbersModel) GetCell(tableID string, row, col int) (*Cell, error) {
	table, exists := m.tableData[tableID]
	if !exists {
		return nil, fmt.Errorf("table %s not found", tableID)
	}

	cellKey := fmt.Sprintf("%d,%d", row, col)
	cellData, exists := table.Cells[cellKey]
	if !exists {
		// Return empty cell if not found
		return NewEmptyCell(row, col), nil
	}

	cell := &Cell{
		Row:          cellData.Row,
		Col:          cellData.Col,
		Type:         cellData.Type,
		value:        cellData.Value,
		formattedVal: cellData.Formatted,
	}

	return cell, nil
}

// UpdateTableData updates the model with table data
func (m *NumbersModel) UpdateTableData(table *Table) error {
	tableData, exists := m.tableData[table.id]
	if !exists {
		// Create new table data
		tableData = &TableData{
			ID:      table.id,
			Name:    table.Name,
			NumRows: table.NumRows,
			NumCols: table.NumCols,
			Cells:   make(map[string]*CellData),
		}
		m.tableData[table.id] = tableData
	}

	// Update cells
	for row := 0; row < table.NumRows; row++ {
		for col := 0; col < table.NumCols; col++ {
			cell := table.cells[row][col]
			cellKey := fmt.Sprintf("%d,%d", row, col)
			
			tableData.Cells[cellKey] = &CellData{
				Row:       cell.Row,
				Col:       cell.Col,
				Type:      cell.Type,
				Value:     cell.Value(),
				Formatted: cell.FormattedValue(),
			}
		}
	}

	return nil
}

// Save saves the model to a file
func (m *NumbersModel) Save(filename string) error {
	// Create output file
	outFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Create zip writer
	zipWriter := zip.NewWriter(outFile)
	defer zipWriter.Close()

	// Copy existing files if updating
	if m.zipReader != nil {
		for _, file := range m.zipReader.File {
			if err := m.copyZipFile(zipWriter, file); err != nil {
				return fmt.Errorf("failed to copy file %s: %w", file.Name, err)
			}
		}
	} else {
		// Create new Numbers file structure
		if err := m.createNewNumbersFile(zipWriter); err != nil {
			return fmt.Errorf("failed to create new Numbers file: %w", err)
		}
	}

	return nil
}

// copyZipFile copies a file from the input zip to the output zip
func (m *NumbersModel) copyZipFile(zipWriter *zip.Writer, file *zip.File) error {
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	w, err := zipWriter.Create(file.Name)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, rc)
	return err
}

// createNewNumbersFile creates a new Numbers file structure
func (m *NumbersModel) createNewNumbersFile(zipWriter *zip.Writer) error {
	// This is a simplified implementation
	// A real implementation would create the proper Numbers file structure
	// with protobuf-encoded data
	
	// Create basic file structure
	files := map[string][]byte{
		"index.xml":     []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?><root></root>"),
		"preview.jpg":   []byte{}, // Empty preview
		"buildVersionHistory.plist": []byte(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0"><dict></dict></plist>`),
	}

	for filename, data := range files {
		w, err := zipWriter.Create(filename)
		if err != nil {
			return err
		}
		if _, err := w.Write(data); err != nil {
			return err
		}
	}

	return nil
}