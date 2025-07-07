package numbers

import (
	"fmt"
)

// Document represents a Numbers spreadsheet document
type Document struct {
	filename string
	Sheets   []*Sheet
	model    *NumbersModel
	styles   map[string]*Style
}

// OpenDocument opens an existing Numbers document from the specified file path
func OpenDocument(filename string) (*Document, error) {
	model, err := NewNumbersModel(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open document: %w", err)
	}

	doc := &Document{
		filename: filename,
		model:    model,
		styles:   make(map[string]*Style),
	}

	if err := doc.loadSheets(); err != nil {
		return nil, fmt.Errorf("failed to load sheets: %w", err)
	}

	return doc, nil
}

// NewDocument creates a new empty Numbers document
func NewDocument() *Document {
	model := NewEmptyNumbersModel()
	doc := &Document{
		model:  model,
		styles: make(map[string]*Style),
		Sheets: []*Sheet{},
	}

	// Create default sheet and table
	sheet := &Sheet{
		Name:   "Sheet 1",
		Tables: []*Table{},
	}

	table := &Table{
		Name:           "Table 1",
		NumRows:        12,
		NumCols:        8,
		NumHeaderRows:  1,
		NumHeaderCols:  1,
		cells:          make([][]*Cell, 12),
	}

	// Initialize cells
	for i := range table.cells {
		table.cells[i] = make([]*Cell, 8)
		for j := range table.cells[i] {
			table.cells[i][j] = NewEmptyCell(i, j)
		}
	}

	sheet.Tables = append(sheet.Tables, table)
	doc.Sheets = append(doc.Sheets, sheet)

	return doc
}

// Close closes the document and releases any resources
func (d *Document) Close() error {
	if d.model != nil {
		return d.model.Close()
	}
	return nil
}

// Save saves the document to the specified filename
func (d *Document) Save(filename string) error {
	if d.model == nil {
		return fmt.Errorf("document model is nil")
	}

	// Update model with current data
	for _, sheet := range d.Sheets {
		for _, table := range sheet.Tables {
			if err := d.model.UpdateTableData(table); err != nil {
				return fmt.Errorf("failed to update table data: %w", err)
			}
		}
	}

	return d.model.Save(filename)
}

// AddSheet adds a new sheet to the document
func (d *Document) AddSheet(sheetName string) *Sheet {
	if sheetName == "" {
		sheetName = fmt.Sprintf("Sheet %d", len(d.Sheets)+1)
	}

	sheet := &Sheet{
		Name:   sheetName,
		Tables: []*Table{},
	}

	// Add default table
	table := &Table{
		Name:           "Table 1",
		NumRows:        12,
		NumCols:        8,
		NumHeaderRows:  1,
		NumHeaderCols:  1,
		cells:          make([][]*Cell, 12),
	}

	// Initialize cells
	for i := range table.cells {
		table.cells[i] = make([]*Cell, 8)
		for j := range table.cells[i] {
			table.cells[i][j] = NewEmptyCell(i, j)
		}
	}

	sheet.Tables = append(sheet.Tables, table)
	d.Sheets = append(d.Sheets, sheet)

	return sheet
}

// AddStyle adds a new style to the document
func (d *Document) AddStyle(style *Style) {
	if style.Name == "" {
		style.Name = fmt.Sprintf("Custom Style %d", len(d.styles)+1)
	}
	d.styles[style.Name] = style
}

// GetStyle returns a style by name
func (d *Document) GetStyle(name string) (*Style, bool) {
	style, exists := d.styles[name]
	return style, exists
}

// loadSheets loads sheets from the underlying model
func (d *Document) loadSheets() error {
	sheetIDs, err := d.model.GetSheetIDs()
	if err != nil {
		return err
	}

	d.Sheets = make([]*Sheet, 0, len(sheetIDs))
	for _, sheetID := range sheetIDs {
		sheet, err := d.loadSheet(sheetID)
		if err != nil {
			return fmt.Errorf("failed to load sheet %s: %w", sheetID, err)
		}
		d.Sheets = append(d.Sheets, sheet)
	}

	return nil
}

// loadSheet loads a single sheet from the model
func (d *Document) loadSheet(sheetID string) (*Sheet, error) {
	sheetName, err := d.model.GetSheetName(sheetID)
	if err != nil {
		return nil, err
	}

	sheet := &Sheet{
		Name:   sheetName,
		Tables: []*Table{},
		id:     sheetID,
	}

	tableIDs, err := d.model.GetTableIDs(sheetID)
	if err != nil {
		return nil, err
	}

	for _, tableID := range tableIDs {
		table, err := d.loadTable(tableID)
		if err != nil {
			return nil, fmt.Errorf("failed to load table %s: %w", tableID, err)
		}
		sheet.Tables = append(sheet.Tables, table)
	}

	return sheet, nil
}

// loadTable loads a single table from the model
func (d *Document) loadTable(tableID string) (*Table, error) {
	tableName, err := d.model.GetTableName(tableID)
	if err != nil {
		return nil, err
	}

	numRows, numCols, err := d.model.GetTableDimensions(tableID)
	if err != nil {
		return nil, err
	}

	table := &Table{
		Name:    tableName,
		NumRows: numRows,
		NumCols: numCols,
		id:      tableID,
		cells:   make([][]*Cell, numRows),
	}

	// Load all cells
	for row := 0; row < numRows; row++ {
		table.cells[row] = make([]*Cell, numCols)
		for col := 0; col < numCols; col++ {
			cell, err := d.model.GetCell(tableID, row, col)
			if err != nil {
				// Create empty cell if there's an error
				table.cells[row][col] = NewEmptyCell(row, col)
			} else {
				table.cells[row][col] = cell
			}
		}
	}

	return table, nil
}

// DefaultTable returns the first table of the first sheet
func (d *Document) DefaultTable() *Table {
	if len(d.Sheets) == 0 || len(d.Sheets[0].Tables) == 0 {
		return nil
	}
	return d.Sheets[0].Tables[0]
}

// GetSheet returns a sheet by name or index
func (d *Document) GetSheet(identifier interface{}) *Sheet {
	switch v := identifier.(type) {
	case int:
		if v >= 0 && v < len(d.Sheets) {
			return d.Sheets[v]
		}
	case string:
		for _, sheet := range d.Sheets {
			if sheet.Name == v {
				return sheet
			}
		}
	}
	return nil
}