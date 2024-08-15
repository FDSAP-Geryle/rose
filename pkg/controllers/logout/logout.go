package logout

import (
	models "rosei/pkg/models/user"
	"rosei/pkg/utils/go-utils/database"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func Logout(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))

	var userRecord models.User
	if err := database.DBConn.First(&userRecord, userID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "User not found"})
	}

	userRecord.Logged = 0
	userRecord.LastLogoutDate = time.Now()
	database.DBConn.Save(&userRecord)

	return c.JSON(fiber.Map{"message": "Logged out successfully"})
}
