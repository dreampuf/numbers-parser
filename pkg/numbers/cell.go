package numbers

import (
	"fmt"
	"time"
)

// CellType represents the type of a cell
type CellType int

const (
	CellTypeEmpty CellType = iota
	CellTypeText
	CellTypeNumber
	CellTypeBool
	CellTypeDate
	CellTypeDuration
	CellTypeError
	CellTypeMerged
)

// Cell represents a cell in a table
type Cell struct {
	Row          int
	Col          int
	Type         CellType
	value        interface{}
	formattedVal string
	style        *Style
	border       *CellBorder
}

// Value returns the cell's value
func (c *Cell) Value() interface{} {
	return c.value
}

// FormattedValue returns the cell's formatted value as a string
func (c *Cell) FormattedValue() string {
	if c.formattedVal != "" {
		return c.formattedVal
	}
	
	if c.value == nil {
		return ""
	}
	
	return fmt.Sprintf("%v", c.value)
}

// SetValue sets the cell's value
func (c *Cell) SetValue(value interface{}) {
	c.value = value
	c.updateType()
}

// SetStyle sets the cell's style
func (c *Cell) SetStyle(style *Style) {
	c.style = style
}

// GetStyle returns the cell's style
func (c *Cell) GetStyle() *Style {
	return c.style
}

// SetBorder sets the cell's border
func (c *Cell) SetBorder(border *CellBorder) {
	c.border = border
}

// GetBorder returns the cell's border
func (c *Cell) GetBorder() *CellBorder {
	return c.border
}

// updateType updates the cell type based on its value
func (c *Cell) updateType() {
	switch c.value.(type) {
	case nil:
		c.Type = CellTypeEmpty
	case string:
		c.Type = CellTypeText
	case int, int64, float32, float64:
		c.Type = CellTypeNumber
	case bool:
		c.Type = CellTypeBool
	case time.Time:
		c.Type = CellTypeDate
	case time.Duration:
		c.Type = CellTypeDuration
	default:
		c.Type = CellTypeText
	}
}

// NewEmptyCell creates a new empty cell
func NewEmptyCell(row, col int) *Cell {
	return &Cell{
		Row:   row,
		Col:   col,
		Type:  CellTypeEmpty,
		value: nil,
	}
}

// NewTextCell creates a new text cell
func NewTextCell(row, col int, value string) *Cell {
	return &Cell{
		Row:   row,
		Col:   col,
		Type:  CellTypeText,
		value: value,
	}
}

// NewNumberCell creates a new number cell
func NewNumberCell(row, col int, value float64) *Cell {
	return &Cell{
		Row:   row,
		Col:   col,
		Type:  CellTypeNumber,
		value: value,
	}
}

// NewBoolCell creates a new boolean cell
func NewBoolCell(row, col int, value bool) *Cell {
	return &Cell{
		Row:   row,
		Col:   col,
		Type:  CellTypeBool,
		value: value,
	}
}

// NewDateCell creates a new date cell
func NewDateCell(row, col int, value time.Time) *Cell {
	return &Cell{
		Row:   row,
		Col:   col,
		Type:  CellTypeDate,
		value: value,
	}
}

// NewDurationCell creates a new duration cell
func NewDurationCell(row, col int, value time.Duration) *Cell {
	return &Cell{
		Row:   row,
		Col:   col,
		Type:  CellTypeDuration,
		value: value,
	}
}

// NewErrorCell creates a new error cell
func NewErrorCell(row, col int) *Cell {
	return &Cell{
		Row:  row,
		Col:  col,
		Type: CellTypeError,
	}
}

// NewMergedCell creates a new merged cell
func NewMergedCell(row, col int) *Cell {
	return &Cell{
		Row:  row,
		Col:  col,
		Type: CellTypeMerged,
	}
}

// IsEmpty returns true if the cell is empty
func (c *Cell) IsEmpty() bool {
	return c.Type == CellTypeEmpty || c.value == nil
}

// IsText returns true if the cell contains text
func (c *Cell) IsText() bool {
	return c.Type == CellTypeText
}

// IsNumber returns true if the cell contains a number
func (c *Cell) IsNumber() bool {
	return c.Type == CellTypeNumber
}

// IsBool returns true if the cell contains a boolean
func (c *Cell) IsBool() bool {
	return c.Type == CellTypeBool
}

// IsDate returns true if the cell contains a date
func (c *Cell) IsDate() bool {
	return c.Type == CellTypeDate
}

// IsDuration returns true if the cell contains a duration
func (c *Cell) IsDuration() bool {
	return c.Type == CellTypeDuration
}

// IsError returns true if the cell contains an error
func (c *Cell) IsError() bool {
	return c.Type == CellTypeError
}

// IsMerged returns true if the cell is merged
func (c *Cell) IsMerged() bool {
	return c.Type == CellTypeMerged
}

// String returns a string representation of the cell
func (c *Cell) String() string {
	return fmt.Sprintf("Cell(%d,%d): %v", c.Row, c.Col, c.value)
}