package controller

import (
	"log"
	merchant "rose/pkg/rose/merchant/model"
	merchantUtils "rose/pkg/rose/merchant/utils"
	utils "rose/pkg/rose/merchant/validation"
	"rose/pkg/utils/go-utils/database"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func DownloadMerchantTemplate(c *fiber.Ctx) error {
	return c.SendFile("assets/template/downloadTemplate.xlsx")
}

func UploadMerchantOK(c *fiber.Ctx) error {
	filepath := merchantUtils.FilePathAssets(c)
	data, err := merchantUtils.ExcelParser(filepath)
	if err != nil {
		log.Println("Error: Failed to parse excel file.")
	}
	utils.ValidateMap(data)

	uploadMerchantOK := merchant.UploadMerchantOK{
		Date:         time.Now().Format("2006-01-02 15:04:05"),
		FilePath:     filepath,
		FilePathErr:  filepath,
		Status:       "",
		Notes:        "",
		User:         "AdminPogi",
		TotalUpload:  "",
		TotalSuccess: "",
		TotalError:   "",
	}
	db := database.DBConn
	db.Create(uploadMerchantOK)
	return nil
}

func Unkown(c *fiber.Ctx) error {
	filepath := merchantUtils.FilePathAssets(c)

	data, err := merchantUtils.ExcelParser(filepath)
	if err != nil {
		log.Println("Error: Failed to parse excel file.")
	}

	receiveUploadActivated := ReceiveUploadActivated(data)

	db := database.DBConn
	db.Model(&merchant.ReceiveUploadActivated{}).Create(SnakeCase(receiveUploadActivated))

	return c.JSON(SnakeCase(receiveUploadActivated))
}

func ReceiveUploadActivated(Data []map[string]interface{}) []map[string]interface{} {
	Columns := []string{
		// id, merchant_name, merchant_address, business_type, business_category, business_location, business_location_type, operational_hour, visiting_hour, barangay, municipality, province, region, status, owner_address, owner_province, owner_region, owner_municipality, owner_barangay, owner_postal_code, agent_id, agent_name,
		// "Uploaded At", "MID", MPAN | "AO Name", "AOCode", "No"
		"Phone Number",
		"Owner Name",
		"ID Type",
		"ID Number",
		"DOB",
		"BIN",
		"CID",
		"Account Number",
		"Account Type",
		"Account Name",
		"Owner Address",
		"Branch Name",
		"Branch Code",
		"Unit Name",
		"Unit Code",
		"Center Name",
		"Center Code",
	}
	UpdateData := merchantUtils.FilterData(Data, Columns)

	InsertColumns := map[string]string{}

	for i := range UpdateData {
		InsertColumns["uploaded_at"] = time.Now().Format("2006-01-02 15:04:05")
		InsertColumns["mid"] = strconv.Itoa(i)  // change
		InsertColumns["mpan"] = strconv.Itoa(i) //change
		merchantUtils.InsertData(UpdateData, InsertColumns, i)
	}

	return UpdateData
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
