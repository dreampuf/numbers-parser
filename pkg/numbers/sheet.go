package numbers

import (
	"fmt"
)

// Sheet represents a sheet within a Numbers document
type Sheet struct {
	Name   string
	Tables []*Table
	id     string
}

// AddTable adds a new table to the sheet
func (s *Sheet) AddTable(tableName string, numRows, numCols int) *Table {
	if tableName == "" {
		tableName = fmt.Sprintf("Table %d", len(s.Tables)+1)
	}

	if numRows <= 0 {
		numRows = 12
	}
	if numCols <= 0 {
		numCols = 8
	}

	table := &Table{
		Name:          tableName,
		NumRows:       numRows,
		NumCols:       numCols,
		NumHeaderRows: 1,
		NumHeaderCols: 1,
		cells:         make([][]*Cell, numRows),
	}

	// Initialize cells
	for i := range table.cells {
		table.cells[i] = make([]*Cell, numCols)
		for j := range table.cells[i] {
			table.cells[i][j] = NewEmptyCell(i, j)
		}
	}

	s.Tables = append(s.Tables, table)
	return table
}

// GetTable returns a table by name or index
func (s *Sheet) GetTable(identifier interface{}) *Table {
	switch v := identifier.(type) {
	case int:
		if v >= 0 && v < len(s.Tables) {
			return s.Tables[v]
		}
	case string:
		for _, table := range s.Tables {
			if table.Name == v {
				return table
			}
		}
	}
	return nil
}

// GetTableByName returns a table by name
func (s *Sheet) GetTableByName(name string) *Table {
	for _, table := range s.Tables {
		if table.Name == name {
			return table
		}
	}
	return nil
}

// RemoveTable removes a table from the sheet
func (s *Sheet) RemoveTable(identifier interface{}) bool {
	var index int = -1

	switch v := identifier.(type) {
	case int:
		if v >= 0 && v < len(s.Tables) {
			index = v
		}
	case string:
		for i, table := range s.Tables {
			if table.Name == v {
				index = i
				break
			}
		}
	}

	if index >= 0 {
		s.Tables = append(s.Tables[:index], s.Tables[index+1:]...)
		return true
	}
	return false
}