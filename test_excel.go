package main

import (
	"fmt"
	"os"
	
	"github.com/xuri/excelize/v2"
	"github.com/lukaszbudnik/migrator/converter"
)

func main() {
	// Create a sample Excel file for testing
	f := excelize.NewFile()
	
	// Create some sample data
	f.SetCellValue("Sheet1", "A1", "name")
	f.SetCellValue("Sheet1", "B1", "age")
	f.SetCellValue("Sheet1", "C1", "city")
	f.SetCellValue("Sheet1", "A2", "John Doe")
	f.SetCellValue("Sheet1", "B2", "30")
	f.SetCellValue("Sheet1", "C2", "New York")
	f.SetCellValue("Sheet1", "A3", "Jane Smith")
	f.SetCellValue("Sheet1", "B3", "25")
	f.SetCellValue("Sheet1", "C3", "Los Angeles")
	f.SetCellValue("Sheet1", "A4", "Bob Johnson")
	f.SetCellValue("Sheet1", "B4", "35")
	f.SetCellValue("Sheet1", "C4", "Chicago")
	
	// Save the file
	if err := f.SaveAs("sample.xlsx"); err != nil {
		fmt.Printf("Error creating sample Excel: %v\n", err)
		return
	}
	
	// Close the file
	f.Close()

	// Test converting Excel to JSON
	fmt.Println("Converting Excel to JSON...")
	jsonData, err := converter.ConvertFile("sample.xlsx", "json")
	if err != nil {
		fmt.Printf("Error converting Excel to JSON: %v\n", err)
		return
	}
	
	fmt.Printf("JSON Output:\n%s\n\n", jsonData)

	// Test converting Excel to CSV
	fmt.Println("Converting Excel to CSV...")
	csvData, err := converter.ConvertFile("sample.xlsx", "csv")
	if err != nil {
		fmt.Printf("Error converting Excel to CSV: %v\n", err)
		return
	}
	
	fmt.Printf("CSV Output:\n%s\n\n", csvData)

	// Clean up
	os.Remove("sample.xlsx")
}