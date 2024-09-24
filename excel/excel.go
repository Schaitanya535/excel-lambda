package excelwriter

import (
	"fmt"
	"io"

	"github.com/xuri/excelize/v2"
)

// WriteToExcelFile function writes data to an excel file in batches
func WriteToExcelFile(headers []string, data []map[string]interface{}, writer io.Writer) error {

	f := excelize.NewFile()

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	sw, err := f.NewStreamWriter("Sheet1")

	if err != nil {
		return err
	}

	// Write data using sw
	for rowIndex, rowData := range data {
		row := make([]interface{}, len(headers))
		for colIndex, header := range headers {
			row[colIndex] = rowData[header]
		}
		cell, err := excelize.CoordinatesToCellName(1, rowIndex+2)
		if err != nil {
			return err
		}
		sw.SetRow(cell, row)
	}

	if err = sw.Flush(); err != nil {
		return err
	}

	// testFile, err := os.Create("test.xlsx")

	// if err != nil {
	// 	return err
	// }
	// f.Write(testFile)

	return f.Write(writer)

}
