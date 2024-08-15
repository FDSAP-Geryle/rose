package uploadmerchant

import (
	"fmt"
	"log"
	"time"

	models "rosei/pkg/models/merchant"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func UploadMerchant(c *fiber.Ctx, isWIP bool) error {
	name, path, err := Filepath(c)
	if err != nil {
		log.Println("Failed to get file path:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get file path",
		})
	}

	ExcelData, err := ExcelParser(path)
	if err != nil {
		log.Println("Failed to parse file:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse file",
		})
	}

	result, validData := ValidateMap(ExcelData)

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

	if isWIP {
		record := models.Recordwip{
			Date:        time.Now(),
			FilePath:    fmt.Sprintf("./assets/received_data/success/%s", name),
			FilePathErr: fmt.Sprintf("./assets/received_data/failed/%s", name),
			Status:      "DONE",
			User:        username,
			Notes:       result.Notes,
		}
		// create a excelfile
		InsertRecord(&record)
	} else {
		record := models.Record{
			Date:        time.Now(),
			FilePath:    fmt.Sprintf("./assets/received_data/success/%s", name),
			FilePathErr: fmt.Sprintf("./assets/received_data/failed/%s", name),
			Status:      "DONE",
			User:        username,
			Notes:       result.Notes,
		}
		// create a excelfile
		InsertRecord(&record)
	}
	InsertReceiveRecord(validData)

	return c.JSON(result)
}
