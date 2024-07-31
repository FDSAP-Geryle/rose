package controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func DownloadMerchant(c *fiber.Ctx) error {
	filename := c.Params("template")
	filePath := fmt.Sprintf("./assets/%s", filename)
	return c.SendFile(filePath)
}
