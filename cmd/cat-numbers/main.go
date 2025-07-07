package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/numbers-parser-go/pkg/numbers"
)

func main() {
	var (
		listTables = flag.Bool("T", false, "List the names of tables and exit")
		listSheets = flag.Bool("S", false, "List the names of sheets and exit")
		brief      = flag.Bool("b", false, "Don't prefix data rows with name of sheet/table")
		version    = flag.Bool("V", false, "Show version and exit")
		sheetName  = flag.String("s", "", "Name of sheet to include in export")
		tableName  = flag.String("t", "", "Name of table to include in export")
		output     = flag.String("o", "", "Output CSV file (default: stdout)")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] document...\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nExport data from Apple Numbers spreadsheet tables\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *version {
		fmt.Println("cat-numbers-go 1.0.0")
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: No document specified\n")
		flag.Usage()
		os.Exit(1)
	}

	for _, filename := range args {
		if err := processDocument(filename, *listTables, *listSheets, *brief, *sheetName, *tableName, *output); err != nil {
			log.Printf("Error processing %s: %v", filename, err)
		}
	}
}

func processDocument(filename string, listTables, listSheets, brief bool, sheetFilter, tableFilter, output string) error {
	doc, err := numbers.OpenDocument(filename)
	if err != nil {
		return fmt.Errorf("failed to open document: %w", err)
	}
	defer doc.Close()

	if listSheets {
		fmt.Printf("Sheets in %s:\n", filepath.Base(filename))
		for i, sheet := range doc.Sheets {
			fmt.Printf("  %d: %s\n", i, sheet.Name)
		}
		return nil
	}

	if listTables {
		fmt.Printf("Tables in %s:\n", filepath.Base(filename))
		for i, sheet := range doc.Sheets {
			fmt.Printf("  Sheet %d (%s):\n", i, sheet.Name)
			for j, table := range sheet.Tables {
				fmt.Printf("    %d: %s\n", j, table.Name)
			}
		}
		return nil
	}

	// Set up output writer
	var writer *csv.Writer
	var outputFile *os.File

	if output != "" {
		outputFile, err = os.Create(output)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer outputFile.Close()
		writer = csv.NewWriter(outputFile)
	} else {
		writer = csv.NewWriter(os.Stdout)
	}
	defer writer.Flush()

	// Export data
	for _, sheet := range doc.Sheets {
		// Filter by sheet name if specified
		if sheetFilter != "" && sheet.Name != sheetFilter {
			continue
		}

		for _, table := range sheet.Tables {
			// Filter by table name if specified
			if tableFilter != "" && table.Name != tableFilter {
				continue
			}

			if err := exportTable(writer, sheet, table, brief); err != nil {
				return fmt.Errorf("failed to export table %s: %w", table.Name, err)
			}
		}
	}

	return nil
}

func exportTable(writer *csv.Writer, sheet *numbers.Sheet, table *numbers.Table, brief bool) error {
	rows, err := table.Rows()
	if err != nil {
		return err
	}

	for rowIndex, row := range rows {
		var record []string

		// Add sheet/table prefix unless brief mode
		if !brief {
			if rowIndex == 0 {
				record = append(record, fmt.Sprintf("Sheet: %s, Table: %s", sheet.Name, table.Name))
				// Fill remaining columns for header alignment
				for len(record) < len(row)+1 {
					record = append(record, "")
				}
				if err := writer.Write(record); err != nil {
					return err
				}
				record = nil
			}
		}

		// Convert cell values to strings
		for _, cell := range row {
			if cell.IsEmpty() {
				record = append(record, "")
			} else {
				record = append(record, numbers.FormatCellValue(cell.Value()))
			}
		}

		if err := writer.Write(record); err != nil {
			return err
		}
	}

	// Add empty line between tables unless brief mode
	if !brief {
		if err := writer.Write([]string{""}); err != nil {
			return err
		}
	}

	return nil
}