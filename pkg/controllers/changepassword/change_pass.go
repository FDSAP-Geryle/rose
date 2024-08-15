package changepassword

import (
	"log"
	models "rosei/pkg/models/user"
	"rosei/pkg/utils/go-utils/database"
	"rosei/pkg/utils/go-utils/passwordHashing"

	"github.com/gofiber/fiber/v2"
)

func ChangePassword(c *fiber.Ctx) error {
	log.Println("Change password attempt")

	// Get user ID from URL parameters
	userid := c.Params("user_id")

	// Parse request body to get the old and new passwords
	var data struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Find the user by ID
	var user models.User
	if err := database.DBConn.Where("id = ?", userid).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Verify the old password
	if !passwordHashing.CheckPasswordHash(data.OldPassword, user.Password) { // Assuming you have a function `CheckPasswordHash`
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Incorrect old password"})
	}

	// Hash the new password
	hashedPassword, err := passwordHashing.HashPassword(data.NewPassword) // Assuming you have a function `HashPassword`
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	// Update user's password
	user.Password = hashedPassword
	if err := database.DBConn.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save new password"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Password changed successfully"})
}

// Password Expiration must chan

// login history
