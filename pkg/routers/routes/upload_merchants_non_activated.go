package routes

import (
	"fmt"
	"os"
	"rose/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

func UploadMerchantNonActivated(c *fiber.Ctx) error {
	// Create Local ./upload directory
	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
		os.Mkdir("./uploads", os.ModePerm)
	}
	// File upload connection to frontend "file"
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	// Save the file to the server/Local
	filePath := fmt.Sprintf("./uploads/%s", file.Filename)
	if err := c.SaveFile(file, filePath); err != nil {
		return err
	}
	// Call ParserExel
	return utils.ParseExcel(c, filePath)
}
