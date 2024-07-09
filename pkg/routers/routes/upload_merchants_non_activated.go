package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
)

func GetUploadMerchantNonActivated(c *fiber.Ctx) error {
	return ParseExcel(c)
}

func ParseExcel(c *fiber.Ctx) error {
	f, err := excelize.OpenFile("/Users/19-sdev-cmdi-rand/Downloads/downloadTemplate.xlsx")
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to open file",
			"error":   err.Error(),
		})
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Get all sheet names
	sheetNames := f.GetSheetList()
	if len(sheetNames) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "No sheets found in the Excel file",
		})
	}

	// Use the first sheet name
	sheetName := sheetNames[0]

	// Get all the rows in the first sheet
	rows, err := f.GetRows(sheetName)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get rows",
			"error":   err.Error(),
		})
	}

	if len(rows) < 2 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "No data found in the sheet",
		})
	}

	// Initialize a slice to hold the data
	var data []map[string]string
	headers := rows[0]

	// Iterate over rows starting from the second row (index 1)
	for _, row := range rows[1:] {
		rowData := make(map[string]string)
		for i, cell := range row {
			if i < len(headers) {
				header := headers[i]
				rowData[header] = cell
			}
		}
		data = append(data, rowData)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    data,
	})
}
