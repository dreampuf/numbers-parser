package numbers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Table represents a table within a sheet
type Table struct {
	Name          string
	NumRows       int
	NumCols       int
	NumHeaderRows int
	NumHeaderCols int
	cells         [][]*Cell
	id            string
}

// Rows returns all rows in the table
func (t *Table) Rows() ([][]*Cell, error) {
	return t.cells, nil
}

// RowsValues returns all cell values in the table (values only)
func (t *Table) RowsValues() ([][]interface{}, error) {
	result := make([][]interface{}, len(t.cells))
	for i, row := range t.cells {
		result[i] = make([]interface{}, len(row))
		for j, cell := range row {
			result[i][j] = cell.Value()
		}
	}
	return result, nil
}

// Cell returns a cell at the specified position
// Supports both (row, col) integers and Excel notation like "A1"
func (t *Table) Cell(args ...interface{}) (*Cell, error) {
	var row, col int
	var err error

	if len(args) == 1 {
		// Excel notation like "A1"
		if addr, ok := args[0].(string); ok {
			row, col, err = parseExcelAddress(addr)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("single argument must be string (Excel notation)")
		}
	} else if len(args) == 2 {
		// Row, col integers
		if r, ok := args[0].(int); ok {
			row = r
		} else {
			return nil, fmt.Errorf("first argument must be int (row)")
		}
		if c, ok := args[1].(int); ok {
			col = c
		} else {
			return nil, fmt.Errorf("second argument must be int (col)")
		}
	} else {
		return nil, fmt.Errorf("expected 1 or 2 arguments")
	}

	if row < 0 || row >= t.NumRows || col < 0 || col >= t.NumCols {
		return nil, fmt.Errorf("cell position (%d, %d) out of bounds", row, col)
	}

	return t.cells[row][col], nil
}

// Write writes a value to a cell at the specified position
func (t *Table) Write(args ...interface{}) error {
	if len(args) < 2 {
		return fmt.Errorf("write requires at least 2 arguments")
	}

	var row, col int
	var value interface{}
	var err error

	// Parse position arguments
	if isExcelAddress(args[0]) {
		// Excel notation: address, value, [style]
		if len(args) < 2 {
			return fmt.Errorf("Excel notation requires address and value")
		}
		if addr, ok := args[0].(string); ok {
			row, col, err = parseExcelAddress(addr)
			if err != nil {
				return err
			}
		}
		value = args[1]
	} else {
		// Row, col notation: row, col, value, [style]
		if len(args) < 3 {
			return fmt.Errorf("coordinate notation requires row, col, and value")
		}
		if r, ok := args[0].(int); ok {
			row = r
		} else {
			return fmt.Errorf("first argument must be int (row) or string (Excel notation)")
		}
		if c, ok := args[1].(int); ok {
			col = c
		} else {
			return fmt.Errorf("second argument must be int (col)")
		}
		value = args[2]
	}

	// Expand table if necessary
	if err := t.expandToFit(row, col); err != nil {
		return err
	}

	// Create cell based on value type
	cell, err := createCellFromValue(row, col, value)
	if err != nil {
		return err
	}

	t.cells[row][col] = cell
	return nil
}

// AddRow adds rows to the table
func (t *Table) AddRow(numRows int, startRow ...int) error {
	if numRows <= 0 {
		numRows = 1
	}

	insertAt := t.NumRows
	if len(startRow) > 0 && startRow[0] >= 0 && startRow[0] <= t.NumRows {
		insertAt = startRow[0]
	}

	// Create new rows
	newRows := make([][]*Cell, numRows)
	for i := range newRows {
		newRows[i] = make([]*Cell, t.NumCols)
		for j := range newRows[i] {
			newRows[i][j] = NewEmptyCell(insertAt+i, j)
		}
	}

	// Insert new rows
	if insertAt == t.NumRows {
		// Append at end
		t.cells = append(t.cells, newRows...)
	} else {
		// Insert in middle
		t.cells = append(t.cells[:insertAt], append(newRows, t.cells[insertAt:]...)...)
	}

	t.NumRows += numRows

	// Update row indices for cells after insertion point
	for i := insertAt + numRows; i < t.NumRows; i++ {
		for j := 0; j < t.NumCols; j++ {
			t.cells[i][j].Row = i
		}
	}

	return nil
}

// AddColumn adds columns to the table
func (t *Table) AddColumn(numCols int, startCol ...int) error {
	if numCols <= 0 {
		numCols = 1
	}

	insertAt := t.NumCols
	if len(startCol) > 0 && startCol[0] >= 0 && startCol[0] <= t.NumCols {
		insertAt = startCol[0]
	}

	// Expand each row
	for i := 0; i < t.NumRows; i++ {
		// Create new cells for this row
		newCells := make([]*Cell, numCols)
		for j := range newCells {
			newCells[j] = NewEmptyCell(i, insertAt+j)
		}

		// Insert new cells
		if insertAt == t.NumCols {
			// Append at end
			t.cells[i] = append(t.cells[i], newCells...)
		} else {
			// Insert in middle
			t.cells[i] = append(t.cells[i][:insertAt], append(newCells, t.cells[i][insertAt:]...)...)
		}

		// Update column indices for cells after insertion point
		for j := insertAt + numCols; j < t.NumCols+numCols; j++ {
			t.cells[i][j].Col = j
		}
	}

	t.NumCols += numCols
	return nil
}

// expandToFit expands the table to accommodate the specified row and column
func (t *Table) expandToFit(row, col int) error {
	// Expand rows if needed
	if row >= t.NumRows {
		rowsToAdd := row - t.NumRows + 1
		if err := t.AddRow(rowsToAdd); err != nil {
			return err
		}
	}

	// Expand columns if needed
	if col >= t.NumCols {
		colsToAdd := col - t.NumCols + 1
		if err := t.AddColumn(colsToAdd); err != nil {
			return err
		}
	}

	return nil
}

// parseExcelAddress parses Excel cell address like "A1" into row, col coordinates
func parseExcelAddress(addr string) (int, int, error) {
	re := regexp.MustCompile(`^([A-Z]+)(\d+)$`)
	matches := re.FindStringSubmatch(strings.ToUpper(addr))
	if len(matches) != 3 {
		return 0, 0, fmt.Errorf("invalid Excel address: %s", addr)
	}

	// Parse column (letters to number)
	col := 0
	for _, char := range matches[1] {
		col = col*26 + int(char-'A') + 1
	}
	col-- // Convert to 0-based

	// Parse row (1-based to 0-based)
	row, err := strconv.Atoi(matches[2])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid row number in address %s: %w", addr, err)
	}
	row-- // Convert to 0-based

	return row, col, nil
}

// isExcelAddress checks if a value is an Excel cell address
func isExcelAddress(val interface{}) bool {
	if addr, ok := val.(string); ok {
		re := regexp.MustCompile(`^[A-Z]+\d+$`)
		return re.MatchString(strings.ToUpper(addr))
	}
	return false
}

// createCellFromValue creates the appropriate cell type based on the value
func createCellFromValue(row, col int, value interface{}) (*Cell, error) {
	switch v := value.(type) {
	case nil:
		return NewEmptyCell(row, col), nil
	case string:
		return NewTextCell(row, col, v), nil
	case int:
		return NewNumberCell(row, col, float64(v)), nil
	case int64:
		return NewNumberCell(row, col, float64(v)), nil
	case float32:
		return NewNumberCell(row, col, float64(v)), nil
	case float64:
		return NewNumberCell(row, col, v), nil
	case bool:
		return NewBoolCell(row, col, v), nil
	case time.Time:
		return NewDateCell(row, col, v), nil
	case time.Duration:
		return NewDurationCell(row, col, v), nil
	default:
		// Convert to string as fallback
		return NewTextCell(row, col, fmt.Sprintf("%v", v)), nil
	}
}