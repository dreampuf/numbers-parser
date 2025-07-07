package main

import (
	"fmt"
	"log"
	"time"

	"github.com/dreampuf/numbers-parser-go/pkg/numbers"
)

func main() {
	// Example 1: Create a new document
	fmt.Println("Creating a new Numbers document...")
	doc := numbers.NewDocument()
	
	// Get the default table
	table := doc.DefaultTable()
	if table == nil {
		log.Fatal("Failed to get default table")
	}

	// Write some data
	table.Write(0, 0, "Name")
	table.Write(0, 1, "Age")
	table.Write(0, 2, "City")
	table.Write(0, 3, "Date")

	table.Write(1, 0, "John Doe")
	table.Write(1, 1, 30)
	table.Write(1, 2, "New York")
	table.Write(1, 3, time.Now())

	table.Write(2, 0, "Jane Smith")
	table.Write(2, 1, 25)
	table.Write(2, 2, "San Francisco")
	table.Write(2, 3, time.Date(2023, 12, 15, 0, 0, 0, 0, time.UTC))

	// Write using Excel notation
	table.Write("A4", "Bob Johnson")
	table.Write("B4", 35)
	table.Write("C4", "Chicago")
	table.Write("D4", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))

	// Save the document
	err := doc.Save("example_output.numbers")
	if err != nil {
		log.Printf("Error saving document: %v", err)
	} else {
		fmt.Println("Document saved as 'example_output.numbers'")
	}

	// Example 2: Read data from the document we just created
	fmt.Println("\nReading data from the document...")
	rows, err := table.Rows()
	if err != nil {
		log.Fatal("Error reading rows:", err)
	}

	for i, row := range rows {
		fmt.Printf("Row %d: ", i)
		for j, cell := range row {
			if j > 0 {
				fmt.Print("\t")
			}
			if cell.IsEmpty() {
				fmt.Print("(empty)")
			} else {
				fmt.Printf("%v", cell.Value())
			}
		}
		fmt.Println()
	}

	// Example 3: Access specific cells
	fmt.Println("\nAccessing specific cells...")
	
	// Access by coordinates
	cell, err := table.Cell(1, 0)
	if err != nil {
		log.Printf("Error accessing cell (1,0): %v", err)
	} else {
		fmt.Printf("Cell (1,0): %v\n", cell.Value())
	}

	// Access by Excel notation
	cell, err = table.Cell("B2")
	if err != nil {
		log.Printf("Error accessing cell B2: %v", err)
	} else {
		fmt.Printf("Cell B2: %v\n", cell.Value())
	}

	// Example 4: Create styled content
	fmt.Println("\nCreating styled content...")
	
	// Create a new sheet for styled content
	styledSheet := doc.AddSheet("Styled Data")
	styledTable := styledSheet.AddTable("Styled Table", 5, 3)

	// Create a style
	headerStyle := numbers.NewStyle()
	headerStyle.Name = "Header Style"
	headerStyle.Bold = true
	headerStyle.FontSize = 14.0
	headerStyle.FontColor = numbers.NewRGB(255, 255, 255)
	headerStyle.BackgroundColor = &numbers.RGB{R: 0, G: 123, B: 255}

	doc.AddStyle(headerStyle)

	// Write styled headers
	styledTable.Write(0, 0, "Product")
	styledTable.Write(0, 1, "Price")
	styledTable.Write(0, 2, "Stock")

	// Apply style to header cells
	for col := 0; col < 3; col++ {
		cell, _ := styledTable.Cell(0, col)
		cell.SetStyle(headerStyle)
	}

	// Add some data
	styledTable.Write(1, 0, "Widget A")
	styledTable.Write(1, 1, 19.99)
	styledTable.Write(1, 2, 100)

	styledTable.Write(2, 0, "Widget B")
	styledTable.Write(2, 1, 29.99)
	styledTable.Write(2, 2, 50)

	// Example 5: Working with different cell types
	fmt.Println("\nDemonstrating different cell types...")
	
	typeTable := styledSheet.AddTable("Cell Types", 8, 2)
	
	typeTable.Write(0, 0, "Type")
	typeTable.Write(0, 1, "Value")
	
	typeTable.Write(1, 0, "Text")
	typeTable.Write(1, 1, "Hello, World!")
	
	typeTable.Write(2, 0, "Number")
	typeTable.Write(2, 1, 3.14159)
	
	typeTable.Write(3, 0, "Boolean")
	typeTable.Write(3, 1, true)
	
	typeTable.Write(4, 0, "Date")
	typeTable.Write(4, 1, time.Now())
	
	typeTable.Write(5, 0, "Duration")
	typeTable.Write(5, 1, 2*time.Hour+30*time.Minute)
	
	typeTable.Write(6, 0, "Empty")
	// Don't write anything to (6, 1) to demonstrate empty cell

	// Check cell types
	for row := 1; row < 7; row++ {
		cell, _ := typeTable.Cell(row, 1)
		fmt.Printf("Row %d cell type: ", row)
		switch {
		case cell.IsEmpty():
			fmt.Println("Empty")
		case cell.IsText():
			fmt.Println("Text")
		case cell.IsNumber():
			fmt.Println("Number")
		case cell.IsBool():
			fmt.Println("Boolean")
		case cell.IsDate():
			fmt.Println("Date")
		case cell.IsDuration():
			fmt.Println("Duration")
		default:
			fmt.Println("Unknown")
		}
	}

	// Save the final document
	err = doc.Save("complete_example.numbers")
	if err != nil {
		log.Printf("Error saving final document: %v", err)
	} else {
		fmt.Println("Complete example saved as 'complete_example.numbers'")
	}

	// Clean up
	doc.Close()
	fmt.Println("\nExample completed successfully!")
}