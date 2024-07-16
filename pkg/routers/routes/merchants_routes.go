package routes

import (
	"fmt"
	"log"
	"os"
	"rose/pkg/models/merchant"
	"rose/pkg/utils"
	"rose/pkg/utils/go-utils/database"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func DownloadMerchant(c *fiber.Ctx) error {
	filename := c.Params("template")
	filePath := fmt.Sprintf("./assets/%s", filename)
	return c.SendFile(filePath)
}

func DownloadMerchantTemplate(c *fiber.Ctx) error {
	return c.SendFile("./template/downloadTemplate.xlsx")
}

func UploadMerchant(c *fiber.Ctx) error {
	if err := ProcessUploads(c); err != nil {
		log.Println("Error processing uploads:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Upload failed",
			"retCode": "400",
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Upload successful",
		"retCode": "200",
	})
}

func GetAssetsPath(c *fiber.Ctx) (string, error) {
	// Create local ./assets directory if it doesn't exist
	if _, err := os.Stat("./assets"); os.IsNotExist(err) {
		os.Mkdir("./assets", os.ModePerm)
	}

	// Retrieve uploaded file and save it to ./assets
	file, err := c.FormFile("file")
	if err != nil {
		return "", err
	}
	filePath := fmt.Sprintf("./assets/%s", file.Filename)
	if err := c.SaveFile(file, filePath); err != nil {
		return "", err
	}
	return filePath, nil
}

func ProcessUploads(c *fiber.Ctx) error {
	// Get file path from request
	filePath, err := GetAssetsPath(c)
	if err != nil {
		log.Println("Error getting file path:", err)
		return err
	}

	// Parse Excel file using custom ExcelParser
	rows, err := utils.ExcelParser(filePath)
	if err != nil {
		log.Println("Error parsing Excel:", err)
		return err
	}

	// Iterate through rows and process data
	for i, row := range rows {
		if i == 0 {
			continue // Skip header row
		}

		// Populate receiveUpload struct
		receiveUpload := merchant.ReceiveUploadActivated{
			PhoneNumber:   getStringValue(row, "Phone Number"),
			OwnerName:     getStringValue(row, "Owner Name"),
			IDType:        getStringValue(row, "ID Type"),
			IDNumber:      getStringValue(row, "ID Number"),
			DOB:           time.Now(), // Adjust as per your business logic
			BIN:           getStringValue(row, "BIN"),
			Cid:           getStringValue(row, "CID"),
			AccountNumber: getStringValue(row, "Account Number"),
			AccountType:   getStringValue(row, "Account Type"),
			AccountName:   getStringValue(row, "Account Name"),
			UploadedAt:    time.Now(), // Adjust as per your business logic
			OwnerAddress:  getStringValue(row, "Owner Address"),
			BranchName:    getStringValue(row, "Branch Name"),
			BranchCode:    getStringValue(row, "Branch Code"),
			UnitName:      getStringValue(row, "Unit Name"),
			CenterName:    getStringValue(row, "Center Name"),
			CenterCode:    getStringValue(row, "Center Code"),
			Mid:           "0001211121121211" + strconv.Itoa(i),
			Mpan:          "00120112121211211" + strconv.Itoa(i),
		}

		// Populate tempMerchantOk struct
		tempMerchantOk := merchant.TempMerchantOk{
			BIN:        getStringValue(row, "BIN"),
			Cid:        getStringValue(row, "CID"),
			BranchName: getStringValue(row, "Branch Name"),
			BranchCode: getStringValue(row, "Branch Code"),
			UnitName:   getStringValue(row, "Unit Name"),
			UnitCode:   getStringValue(row, "Unit Code"),
			CenterName: getStringValue(row, "Center Name"),
			CenterCode: getStringValue(row, "Center Code"),
			AOName:     getStringValue(row, "AO Name"),
			AOCode:     getStringValue(row, "AO Code"),
		}

		// Start database transaction
		db := database.DBConn.Begin()
		if err := db.Create(&receiveUpload).Error; err != nil {
			db.Rollback()
			log.Println("Error creating receiveUpload record:", err)
			return err
		}
		if err := db.Create(&tempMerchantOk).Error; err != nil {
			db.Rollback()
			log.Println("Error creating tempMerchantOk record:", err)
			return err
		}
		db.Commit() // Commit transaction
	}

	return nil
}

func getStringValue(row map[string]interface{}, key string) string {
	if value, ok := row[key]; ok {
		if strValue, ok := value.(string); ok {
			return strValue
		}
	}
	return ""
}
