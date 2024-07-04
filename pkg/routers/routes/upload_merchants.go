package routes

import (
	"log"
	"rose/pkg/utils/go-utils/database"

	"strconv"

	"github.com/gofiber/fiber/v2"
)

// GetUploadMerchant retrieves data from get_uploadmerchant PostgreSQL function with pagination
func GetUploadMerchant(c *fiber.Ctx) error {
	// Parse query parameters for pagination
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	perPage, err := strconv.Atoi(c.Query("perPage", "10"))
	if err != nil || perPage < 1 {
		perPage = 10
	}

	// Calculate offset for pagination
	offset := (page - 1) * perPage

	// Initialize a Gorm DB connection
	db := database.DBConn // Assuming DBConn is initialized properly with Gorm

	// Retrieve total count of records
	var totalCount int64
	if err := db.Raw("SELECT COUNT(*) FROM get_uploadmerchant()").Scan(&totalCount).Error; err != nil {
		log.Printf("Failed to fetch total count: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch total count",
			"error":   err.Error(),
		})
	}

	// Retrieve records with pagination
	var records []map[string]interface{} // Adjust the type based on your expected data structure

	if err := db.Raw("SELECT * FROM get_uploadmerchant() LIMIT ? OFFSET ?", perPage, offset).Scan(&records).Error; err != nil {
		log.Printf("Failed to fetch records: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch records",
			"error":   err.Error(),
		})
	}

	// Calculate total pages for pagination
	totalPages := (totalCount + int64(perPage) - 1) / int64(perPage)

	// Prepare JSON response structure
	response := fiber.Map{
		"records":     records,
		"currentPage": page,
		"totalPages":  totalPages,
		"totalCount":  totalCount,
	}

	// Return JSON response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success fetching records",
		"data":    response,
	})
}
