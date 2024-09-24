package excelwriter

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/xuri/excelize/v2"
)

// func TestWriteToExcelFile(t *testing.T) {
// 	headers := []string{"Name", "Age", "City"}
// 	data := []map[string]interface{}{
// 		{"Name": "Alice", "Age": 30, "City": "New York"},
// 		{"Name": "Bob", "Age": 25, "City": "Los Angeles"},
// 		{"Name": "Charlie", "Age": 35, "City": "Chicago"},
// 	}

// 	var buf bytes.Buffer
// 	err := WriteToExcelFile(headers, data, &buf)
// 	if err != nil {
// 		t.Fatalf("expected no error, got %v", err)
// 	}

// }

func TestWriteToExcelFile(t *testing.T) {
	headers := []string{"Name", "Age", "City"}
	data := []map[string]interface{}{
		{"Name": "Alice", "Age": 30, "City": "New York"},
		{"Name": "Bob", "Age": 25, "City": "Los Angeles"},
		{"Name": "Charlie", "Age": 35, "City": "Chicago"},
	}

	var buf bytes.Buffer
	err := WriteToExcelFile(headers, data, &buf)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	f, err := excelize.OpenReader(&buf)
	if err != nil {
		t.Fatalf("expected no error opening excel file, got %v", err)
	}

	sheetName := "Sheet1"
	// for i, header := range headers {
	// 	cell := fmt.Sprintf("%c1", 'A'+i)
	// 	value, err := f.GetCellValue(sheetName, cell)
	// 	if err != nil {
	// 		t.Fatalf("expected no error getting cell value, got %v", err)
	// 	}
	// 	if value != header {
	// 		t.Errorf("expected header %s, got %s", header, value)
	// 	}
	// }

	for rowIndex, rowData := range data {
		for colIndex, header := range headers {
			cell, err := excelize.CoordinatesToCellName(colIndex+1, rowIndex+2)
			if err != nil {
				t.Fatalf("expected no error converting coordinates to cell name, got %v", err)
			}
			value, err := f.GetCellValue(sheetName, cell)
			if err != nil {
				t.Fatalf("expected no error getting cell value, got %v", err)
			}
			if value != fmt.Sprintf("%v", rowData[header]) {
				t.Errorf("expected cell value %v, got %s", rowData[header], value)
			}
		}
	}
}
