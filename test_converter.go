package main

import (
	"fmt"
	"os"
	
	"github.com/lukaszbudnik/migrator/converter"
)

func main() {
	// Create a sample CSV file for testing
	csvContent := `name,age,city
John Doe,30,New York
Jane Smith,25,Los Angeles
Bob Johnson,35,Chicago`

	err := os.WriteFile("sample.csv", []byte(csvContent), 0644)
	if err != nil {
		fmt.Printf("Error creating sample CSV: %v\n", err)
		return
	}

	// Test converting CSV to JSON
	fmt.Println("Converting CSV to JSON...")
	jsonData, err := converter.ConvertFile("sample.csv", "json")
	if err != nil {
		fmt.Printf("Error converting CSV to JSON: %v\n", err)
		return
	}
	
	fmt.Printf("JSON Output:\n%s\n\n", jsonData)

	// Test converting CSV to SQL
	fmt.Println("Converting CSV to SQL...")
	sqlData, err := converter.ConvertFile("sample.csv", "sql")
	if err != nil {
		fmt.Printf("Error converting CSV to SQL: %v\n", err)
		return
	}
	
	fmt.Printf("SQL Output:\n%s\n\n", sqlData)

	// Clean up
	os.Remove("sample.csv")
}