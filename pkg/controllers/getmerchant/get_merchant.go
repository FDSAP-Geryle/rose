package merchants

import (
	"fmt"
	"rosei/pkg/utils/go-utils/database"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func Get[T any](c *fiber.Ctx) error {
	fmt.Println("Merchant records")

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid page number")
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid limit number")
	}

	offset := (page - 1) * limit

	var totalCount int64
	result := database.DBConn.Model(new(T)).Count(&totalCount)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Database query failed: " + result.Error.Error())
	}

	totalPages := int((totalCount + int64(limit) - 1) / int64(limit))

	var records []T
	result = database.DBConn.Limit(limit).Offset(offset).Find(&records)
	if result.Error != nil {
		fmt.Println("Database query error:", result.Error)
		return c.Status(fiber.StatusInternalServerError).SendString("Database query failed: " + result.Error.Error())
	}

	response := fiber.Map{
		"records":     records,
		"currentPage": page,
		"totalPages":  totalPages,
		"totalCount":  totalCount,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"retCode": "200",
		"message": "Success fetching records",
		"data":    response,
	})
}
