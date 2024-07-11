package routes

import (
	"fmt"
	"os"
	"rose/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

func DownloadMerchantTemplate(c *fiber.Ctx) error {
	return c.SendFile("./template/downloadTemplate.xlsx")
}
func DownloadMerchant(c *fiber.Ctx) error {
	filename := c.Params("template")
	filePath := fmt.Sprintf("./assets/%s", filename)
	return c.SendFile(filePath)
}
func UploadMerchant(c *fiber.Ctx) error {
	// Create Local ./upload directory
	if _, err := os.Stat("./assets"); os.IsNotExist(err) {
		os.Mkdir("./assets", os.ModePerm)
	}
	// File upload connection to frontend "file"
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	// Save the file to the server/Local
	filePath := fmt.Sprintf("./assets/%s", file.Filename)
	if err := c.SaveFile(file, filePath); err != nil {
		return err
	}
	// Call ParserExel
	return utils.ParseExcel(c, filePath)
}
