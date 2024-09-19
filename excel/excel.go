package excelwriter

import (
	"fmt"
	"io"

	"github.com/xuri/excelize/v2"
)

// WriteToExcelFile function writes data to an excel file in batches
func WriteToExcelFile(headers []string, data [][]map[string]string, batchSize int, writer io.Writer) error {
	f := excelize.NewFile()

	// Write headers
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue("Sheet1", cell, header)
	}

	// Write data in batches
	for i, batch := range data {
		startRow := i*batchSize + 2
		for rowIndex, rowData := range batch {
			for colIndex, header := range headers {
				cell := fmt.Sprintf("%c%d", 'A'+colIndex, startRow+rowIndex)
				f.SetCellValue("Sheet1", cell, rowData[header])
			}
		}
	}

	// Write the file to the provided writer (can be a buffer)
	return f.Write(writer)
}
