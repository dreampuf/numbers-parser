# numbers-parser-go API Documentation

## Overview

The `numbers-parser-go` library provides a Go interface for reading and writing Apple Numbers spreadsheet files. This document describes the main types and functions available in the library.

## Package Structure

```
github.com/numbers-parser-go/pkg/numbers
```

## Core Types

### Document

The `Document` type represents a complete Numbers spreadsheet file.

```go
type Document struct {
    Sheets []*Sheet
    // ... other fields
}
```

#### Constructor Functions

- `NewDocument() *Document` - Creates a new empty document
- `OpenDocument(filename string) (*Document, error)` - Opens an existing Numbers file

#### Methods

- `Close() error` - Closes the document and releases resources
- `Save(filename string) error` - Saves the document to a file
- `AddSheet(sheetName string) *Sheet` - Adds a new sheet
- `AddStyle(style *Style)` - Adds a style to the document
- `GetStyle(name string) (*Style, bool)` - Retrieves a style by name
- `GetSheet(identifier interface{}) *Sheet` - Gets a sheet by name or index
- `DefaultTable() *Table` - Returns the first table of the first sheet

### Sheet

The `Sheet` type represents a worksheet within a document.

```go
type Sheet struct {
    Name   string
    Tables []*Table
}
```

#### Methods

- `AddTable(tableName string, numRows, numCols int) *Table` - Adds a new table
- `GetTable(identifier interface{}) *Table` - Gets a table by name or index
- `GetTableByName(name string) *Table` - Gets a table by name
- `RemoveTable(identifier interface{}) bool` - Removes a table

### Table

The `Table` type represents a table within a sheet.

```go
type Table struct {
    Name          string
    NumRows       int
    NumCols       int
    NumHeaderRows int
    NumHeaderCols int
}
```

#### Methods

- `Rows() ([][]*Cell, error)` - Returns all rows as a 2D slice of cells
- `RowsValues() ([][]interface{}, error)` - Returns all cell values
- `Cell(args ...interface{}) (*Cell, error)` - Gets a cell by coordinates or Excel notation
- `Write(args ...interface{}) error` - Writes a value to a cell
- `AddRow(numRows int, startRow ...int) error` - Adds rows to the table
- `AddColumn(numCols int, startCol ...int) error` - Adds columns to the table
- `GetCellRange(rangeStr string) ([]*Cell, error)` - Gets cells in a range
- `SetCellRange(rangeStr string, values [][]interface{}) error` - Sets values for a range

### Cell

The `Cell` type represents an individual cell in a table.

```go
type Cell struct {
    Row  int
    Col  int
    Type CellType
}
```

#### Methods

- `Value() interface{}` - Returns the cell's value
- `FormattedValue() string` - Returns the formatted value as a string
- `SetValue(value interface{})` - Sets the cell's value
- `SetStyle(style *Style)` - Sets the cell's style
- `GetStyle() *Style` - Gets the cell's style
- `SetBorder(border *CellBorder)` - Sets the cell's border
- `GetBorder() *CellBorder` - Gets the cell's border

#### Type Check Methods

- `IsEmpty() bool` - Returns true if the cell is empty
- `IsText() bool` - Returns true if the cell contains text
- `IsNumber() bool` - Returns true if the cell contains a number
- `IsBool() bool` - Returns true if the cell contains a boolean
- `IsDate() bool` - Returns true if the cell contains a date
- `IsDuration() bool` - Returns true if the cell contains a duration
- `IsError() bool` - Returns true if the cell contains an error
- `IsMerged() bool` - Returns true if the cell is merged

### CellType

```go
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
```

## Styling Types

### Style

The `Style` type represents cell formatting and styling.

```go
type Style struct {
    Name            string
    Alignment       *Alignment
    BackgroundImage *BackgroundImage
    BackgroundColor *RGB
    FontColor       RGB
    FontSize        float64
    FontName        string
    Bold            bool
    Italic          bool
    Strikethrough   bool
    Underline       bool
    // ... other fields
}
```

#### Constructor

- `NewStyle() *Style` - Creates a new style with default values

### RGB

```go
type RGB struct {
    R, G, B uint8
}
```

#### Constructor

- `NewRGB(r, g, b uint8) RGB` - Creates a new RGB color

### Alignment

```go
type Alignment struct {
    Horizontal HorizontalAlignment
    Vertical   VerticalAlignment
}
```

#### Constructor

- `NewAlignment(horizontal HorizontalAlignment, vertical VerticalAlignment) *Alignment`

### Border Types

```go
type Border struct {
    Width float64
    Color RGB
    Style BorderType
}

type CellBorder struct {
    Top    *Border
    Right  *Border
    Bottom *Border
    Left   *Border
}
```

## Utility Functions

### Excel Address Conversion

- `ExcelColumnToNumber(col string) (int, error)` - Converts Excel column letters to number
- `NumberToExcelColumn(col int) string` - Converts column number to Excel letters
- `ExcelAddressToCoords(addr string) (int, int, error)` - Converts "A1" to coordinates
- `CoordsToExcelAddress(row, col int) string` - Converts coordinates to "A1"
- `IsValidExcelAddress(addr string) bool` - Validates Excel address format
- `ParseRange(rangeStr string) (startRow, startCol, endRow, endCol int, err error)` - Parses range like "A1:C3"

### Cell Value Formatting

- `FormatCellValue(value interface{}) string` - Formats a cell value for display

## Constants

### Default Values

```go
const (
    DefaultRowCount    = 12
    DefaultColCount    = 8
    DefaultHeaderRows  = 1
    DefaultHeaderCols  = 1
    DefaultFontSize    = 12.0
    DefaultFontName    = "Helvetica"
    DefaultTextInset   = 1.0
    DefaultBorderWidth = 0.35
)
```

### Predefined Colors

```go
var (
    ColorBlack = RGB{R: 0, G: 0, B: 0}
    ColorWhite = RGB{R: 255, G: 255, B: 255}
    ColorRed   = RGB{R: 255, G: 0, B: 0}
    ColorGreen = RGB{R: 0, G: 255, B: 0}
    ColorBlue  = RGB{R: 0, G: 0, B: 255}
)
```

## Error Handling

All functions that can fail return an error as their last return value. Common error conditions include:

- Invalid file format
- File not found
- Invalid cell addresses
- Out of bounds access
- Insufficient arguments

## Examples

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/numbers-parser-go/pkg/numbers"
)

func main() {
    // Create a new document
    doc := numbers.NewDocument()
    defer doc.Close()
    
    // Get the default table
    table := doc.DefaultTable()
    
    // Write some data
    table.Write(0, 0, "Name")
    table.Write(0, 1, "Age")
    table.Write(1, 0, "John")
    table.Write(1, 1, 30)
    
    // Save the document
    err := doc.Save("example.numbers")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Document saved successfully!")
}
```

### Reading Data

```go
// Open an existing document
doc, err := numbers.OpenDocument("data.numbers")
if err != nil {
    log.Fatal(err)
}
defer doc.Close()

// Get the first table
table := doc.DefaultTable()

// Read all rows
rows, err := table.Rows()
if err != nil {
    log.Fatal(err)
}

// Print data
for i, row := range rows {
    fmt.Printf("Row %d: ", i)
    for j, cell := range row {
        if j > 0 {
            fmt.Print(", ")
        }
        fmt.Print(cell.Value())
    }
    fmt.Println()
}
```

### Using Styles

```go
// Create a style
headerStyle := numbers.NewStyle()
headerStyle.Name = "Header"
headerStyle.Bold = true
headerStyle.FontSize = 14.0
headerStyle.FontColor = numbers.NewRGB(255, 255, 255)
headerStyle.BackgroundColor = &numbers.RGB{R: 0, G: 123, B: 255}

// Add style to document
doc.AddStyle(headerStyle)

// Apply style to a cell
cell, _ := table.Cell(0, 0)
cell.SetStyle(headerStyle)
```

### Excel Notation

```go
// Write using Excel notation
table.Write("A1", "Header")
table.Write("B1", "Value")

// Read using Excel notation
cell, err := table.Cell("A1")
if err != nil {
    log.Fatal(err)
}
fmt.Println(cell.Value())
```

## Limitations

This is a simplified implementation focusing on core functionality. Some advanced Numbers features are not supported:

- Complex formulas
- Pivot tables
- Advanced formatting options
- Password-protected files
- Some chart types
- Advanced image handling

For production use, you may need to extend the implementation based on your specific requirements.