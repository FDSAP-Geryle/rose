package utils

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

func ExcelParser(filepath string) ([]map[string]interface{}, error) {
	// Open file
	f, err := excelize.OpenFile(filepath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Ensure the file is closed at the end of function execution
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err) // Handle closing error, if any
		}
	}()

	// Get all sheet names
	sheetNames := f.GetSheetList()
	if len(sheetNames) == 0 {
		return nil, fmt.Errorf("no sheets found in the Excel file")
	}

	// Use the first sheet name
	sheetName := sheetNames[0]

	// Get all the rows in the first sheet
	rows, err := f.GetRows(sheetName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// If there are no rows or only headers, return an error
	if len(rows) < 2 {
		return nil, fmt.Errorf("no data rows found in the Excel sheet")
	}

	// Prepare to store data
	var data []map[string]interface{}

	// Extract headers from the first row
	headers := rows[0]

	// Iterate over rows and create map for each row
	for _, row := range rows[1:] {
		rowData := make(map[string]interface{})
		for j, cell := range row {
			// Use header as key and cell value as interface{} in map
			rowData[headers[j]] = cell
		}
		data = append(data, rowData)
	}

	return data, nil
}
