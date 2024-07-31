package merchantUtils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/xuri/excelize/v2"
)

// ExcelParser parses an Excel file and returns a slice of maps with column names as keys
func ExcelParser(filepath string) ([]map[string]interface{}, error) {
	f, err := excelize.OpenFile(filepath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Println("Failed to close file:", err)
		}
	}()

	// Assuming the data is in the first sheet
	sheetName := f.GetSheetName(0)

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	if len(rows) < 1 {
		return nil, fmt.Errorf("no data found in the sheet")
	}

	// Get header from the first row
	header := rows[0]

	var result []map[string]interface{}

	// Iterate over rows and map data to the header
	for _, row := range rows[1:] {
		rowData := make(map[string]interface{})
		for colIndex, cellValue := range row {
			if colIndex < len(header) {
				rowData[header[colIndex]] = cellValue
			}
		}
		result = append(result, rowData)
	}

	return result, nil
}

func InsertData(data []map[string]interface{}, Column map[string]string, rowIndex int) []map[string]interface{} {
	if rowIndex < 0 || rowIndex >= len(data) {
		fmt.Println("Invalid row index")
		return data
	}

	row := data[rowIndex]
	for key, val := range Column {
		row[key] = val
	}

	return data
}

func FilterData(data []map[string]interface{}, columnNames []string) []map[string]interface{} {
	var filteredData []map[string]interface{}

	for _, row := range data {
		filteredRow := make(map[string]interface{})
		for _, col := range columnNames {
			if value, ok := row[col]; ok {
				filteredRow[col] = value
			}
		}
		filteredData = append(filteredData, filteredRow)
	}

	return filteredData
}

func FilePathAssets(c *fiber.Ctx) string {
	// Create local ./assets directory if it doesn't exist
	if _, err := os.Stat("./assets"); os.IsNotExist(err) {
		os.Mkdir("./assets", os.ModePerm)
		os.Mkdir("./assets/template", os.ModePerm)
		os.Mkdir("./assets/success", os.ModePerm)
		os.Mkdir("./assets/failed", os.ModePerm)
	}

	// Retrieve uploaded file and save it to ./assets
	file, err := c.FormFile("file")
	if err != nil {
		return ""
	}

	layout := "2006-01-02T15_04_05"

	// Save to assets
	filePath := fmt.Sprintf("./assets/%s.xlsx", time.Now().Format(layout))
	if err := c.SaveFile(file, filePath); err != nil {
		return ""
	}
	return filePath
}
