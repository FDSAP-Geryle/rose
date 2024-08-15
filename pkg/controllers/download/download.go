package download

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// download path err success

func DownloadErr(c *fiber.Ctx) error { //can add download err and download success folder
	filename := c.Params("filename")
	filePath := fmt.Sprintf("./assets/received_data/failed/%s", filename)
	return c.SendFile(filePath)
}

func DownloadSucc(c *fiber.Ctx) error { //can add download err and download success folder
	filename := c.Params("filename")
	filePath := fmt.Sprintf("./assets/received_data/success/%s", filename)
	return c.SendFile(filePath)
}

func DownloadMerchantTemplate(c *fiber.Ctx) error {
	return c.SendFile("assets/template/MerchantTemplate.xlsx")
}
