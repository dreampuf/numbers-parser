# numbers-parser-go

A Go library for parsing Apple Numbers `.numbers` files. This is a Go port of the Python [numbers-parser](https://github.com/masaccio/numbers-parser) library.

## Features

- Read and parse Apple Numbers spreadsheet files
- Extract data from sheets and tables
- Support for different cell types (text, numbers, dates, etc.)
- Cell formatting and styling
- Write support for creating new Numbers documents

## Installation

```bash
go get github.com/numbers-parser-go
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/numbers-parser-go/pkg/numbers"
)

func main() {
    // Open a Numbers document
    doc, err := numbers.OpenDocument("mydoc.numbers")
    if err != nil {
        log.Fatal(err)
    }
    defer doc.Close()
    
    // Get the first sheet and table
    sheet := doc.Sheets[0]
    table := sheet.Tables[0]
    
    // Read cell data
    rows, err := table.Rows()
    if err != nil {
        log.Fatal(err)
    }
    
    // Print first row
    for _, cell := range rows[0] {
        fmt.Printf("%v\t", cell.Value())
    }
}
```

## API Documentation

See the [docs](docs/) directory for detailed API documentation.

## License

MIT License - see LICENSE file for details.
