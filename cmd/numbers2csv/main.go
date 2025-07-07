package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dreampuf/numbers-parser-go/pkg/numbers"
)

func main() {
	var (
		output    = flag.String("o", "", "Output CSV file (default: replace .numbers with .csv)")
		sheetName = flag.String("sheet", "", "Name of sheet to export (default: first sheet)")
		tableName = flag.String("table", "", "Name of table to export (default: first table)")
		version   = flag.Bool("version", false, "Show version and exit")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] input.numbers\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nConvert Apple Numbers documents to CSV format\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *version {
		fmt.Println("numbers2csv-go 1.0.0")
		return
	}

	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Error: Exactly one input file required\n")
		flag.Usage()
		os.Exit(1)
	}

	inputFile := args[0]

	// Determine output filename
	outputFile := *output
	if outputFile == "" {
		if strings.HasSuffix(strings.ToLower(inputFile), ".numbers") {
			outputFile = inputFile[:len(inputFile)-8] + ".csv"
		} else {
			outputFile = inputFile + ".csv"
		}
	}

	if err := convertToCSV(inputFile, outputFile, *sheetName, *tableName); err != nil {
		log.Fatalf("Error converting file: %v", err)
	}

	fmt.Printf("Successfully converted %s to %s\n", filepath.Base(inputFile), filepath.Base(outputFile))
}

func convertToCSV(inputFile, outputFile, sheetName, tableName string) error {
	// Open Numbers document
	doc, err := numbers.OpenDocument(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open Numbers document: %w", err)
	}
	defer doc.Close()

	// Find the sheet to export
	var sheet *numbers.Sheet
	if sheetName != "" {
		sheet = doc.GetSheet(sheetName)
		if sheet == nil {
			return fmt.Errorf("sheet '%s' not found", sheetName)
		}
	} else {
		if len(doc.Sheets) == 0 {
			return fmt.Errorf("no sheets found in document")
		}
		sheet = doc.Sheets[0]
	}

	// Find the table to export
	var table *numbers.Table
	if tableName != "" {
		table = sheet.GetTable(tableName)
		if table == nil {
			return fmt.Errorf("table '%s' not found in sheet '%s'", tableName, sheet.Name)
		}
	} else {
		if len(sheet.Tables) == 0 {
			return fmt.Errorf("no tables found in sheet '%s'", sheet.Name)
		}
		table = sheet.Tables[0]
	}

	// Create output file
	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Create CSV writer
	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	// Get table data
	rows, err := table.Rows()
	if err != nil {
		return fmt.Errorf("failed to read table data: %w", err)
	}

	// Write data to CSV
	for _, row := range rows {
		var record []string
		for _, cell := range row {
			if cell.IsEmpty() {
				record = append(record, "")
			} else {
				record = append(record, numbers.FormatCellValue(cell.Value()))
			}
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	return nil
}