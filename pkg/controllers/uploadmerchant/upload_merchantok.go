package uploadmerchant

import (
	"fmt"
	"log"
	"time"

	models "rosei/pkg/models/merchant"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func UploadMerchantok(c *fiber.Ctx) error {
	// Handle file upload and path retrieval
	name, path, err := Filepath(c)
	if err != nil {
		log.Println("Failed to get file path:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get file path",
		})
	}

	// Parse the Excel file
	data, err := ExcelParser(path)
	if err != nil {
		log.Println("Failed to parse file:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse file",
		})
	}

	// Validate the parsed data
	result, validData := ValidateMap(data)

	// Extract user from JWT claims
	user := c.Locals("user").(*jwt.Token)
	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("Invalid JWT claims")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid JWT claims",
		})
	}

	username, ok := claims["username"].(string)
	if !ok {
		log.Println("Username claim is missing or not a string")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Username claim is missing or not a string",
		})
	}

	// Prepare the record for saving
	record := models.Recordok{
		Date:         time.Now(),
		FilePath:     fmt.Sprintf("./assets/received_data/success/%s", name),
		FilePathErr:  fmt.Sprintf("./assets/received_data/failed/%s", name),
		Status:       "DONE",
		User:         username,
		Notes:        result.Notes,
		TotalUpload:  result.TotalUpload,
		TotalSuccess: result.TotalSuccess,
		TotalError:   result.TotalFailed,
	}

	// Save the record to the database
	InsertRecord(&record) // Pass a pointer to the record

	InsertReceiveRecord(validData)

	// Return the validation result as JSON
	return c.JSON(result)
}
