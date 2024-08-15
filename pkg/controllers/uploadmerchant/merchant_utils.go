package uploadmerchant

import (
	"fmt"
	"log"
	"os"
	models "rosei/pkg/models/merchant"
	"rosei/pkg/utils/go-utils/database"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
)

func Filepath(c *fiber.Ctx) (name string, path string, err error) {
	file, err := c.FormFile("file")
	if err != nil {
		return "", "", err
	}

	layout := "2006-01-02T15_04_05"
	name = fmt.Sprintf("%s.xlsx", time.Now().Format(layout))
	path = fmt.Sprintf("./assets/received_data/excelfile/%s", name)

	if err := c.SaveFile(file, path); err != nil {
		return name, "", err
	}
	return name, path, nil
}

type ValidationResult struct {
	TotalUpload  int
	TotalSuccess int
	TotalFailed  int
	Notes        string
}

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

func FilePathAssets(isSuccess bool, c *fiber.Ctx) string {
	// Check and create required directories if they don't exist
	if _, err := os.Stat("./assets"); os.IsNotExist(err) {
		os.Mkdir("./assets", os.ModePerm)
		os.Mkdir("./assets/template", os.ModePerm)
		os.Mkdir("./assets/original_files", os.ModePerm)
		os.Mkdir("./assets/received_data/success", os.ModePerm)
		os.Mkdir("./assets/received_data/failed", os.ModePerm)
	}

	// Retrieve the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return ""
	}

	// Save the original file in the "original_files" directory
	layout := "2006-01-02T15_04_05"
	originalFilePath := fmt.Sprintf("./assets/original_files/%s.xlsx", time.Now().Format(layout))
	if err := c.SaveFile(file, originalFilePath); err != nil {
		return ""
	}

	// Simulate processing the file to determine success or failure
	// Replace this with actual file processing logic
	var finalFilePath string

	if isSuccess {
		// If processing is successful, move file to "success" folder
		finalFilePath = fmt.Sprintf("./assets/received_data/success/%s.xlsx", time.Now().Format(layout))
	} else {
		// If processing fails, move file to "failed" folder
		finalFilePath = fmt.Sprintf("./assets/received_data/failed/%s.xlsx", time.Now().Format(layout))
	}

	// Move the file to the final destination
	if err := os.Rename(originalFilePath, finalFilePath); err != nil {
		return ""
	}

	return finalFilePath
}

func SnakeCase(data []map[string]interface{}) []map[string]interface{} {
	for _, row := range data {
		for key, value := range row {
			nospace := strings.ReplaceAll(key, " ", "_")
			lowernospace := strings.ToLower(nospace)
			delete(row, key)
			row[lowernospace] = value
		}
	}
	return data
}

func InsertRecord(record models.Recordable) {
	tableName := record.TableName()
	fmt.Printf("Inserting record into table: %s\n", tableName)
	if err := database.DBConn.Create(record).Error; err != nil {
		log.Println("Failed to save record:", err)
	}
}
