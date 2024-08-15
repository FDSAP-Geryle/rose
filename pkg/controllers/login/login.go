package login

import (
	"log"
	"strings"
	"time"

	models "rosei/pkg/models/user"
	"rosei/pkg/utils/go-utils/database"
	utils "rosei/pkg/utils/go-utils/fiber"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *fiber.Ctx) error {
	// Log request without sensitive details
	log.Println("Login attempt")

	var input models.User
	if err := c.BodyParser(&input); err != nil {
		log.Println("Error parsing body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	var user models.User
	result := database.DBConn.Where("user_name = ?", input.UserName).First(&user)
	if result.Error != nil {
		log.Println("Login failed for username:", input.UserName)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid username or password"})
	}

	if user.IsLock == 1 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Account is locked"})
	}

	// Check password and update failed login attempts if necessary
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		err := HandleFailedLogin(user.UserName)
		if err != nil {
			log.Println("Error updating user login details:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update login details"})
		}
		log.Println("Password mismatch for username:", input.UserName)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid username or password"})
	}

	if IsPasswordExpired(user) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Password has expired. Please reset your password."})
	}

	// Update only specific fields upon successful login
	err := database.DBConn.Model(&user).Updates(models.User{
		LastLoginDate: time.Now(),
		LastLoginFrom: GetDeviceInfo(c),
		Logged:        1,
	}).Error
	if err != nil {
		log.Println("Error updating user login details:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update login details"})
	}

	// Generate token and handle potential errors
	tokenString, err := utils.GenerateJWTToken(user.UserName, user.ID)
	if err != nil {
		log.Println("Error generating token for username:", input.UserName)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.JSON(fiber.Map{"message": "success", "token": tokenString})
}

func IsPasswordExpired(user models.User) bool {
	return time.Now().After(user.PwdExpiredDate)
}

func GetDeviceInfo(c *fiber.Ctx) string {
	userAgent := c.Get("User-Agent")
	if strings.Contains(userAgent, "Mobile") {
		return "Mobile Device"
	} else if strings.Contains(userAgent, "Windows") {
		return "Windows PC"
	} else if strings.Contains(userAgent, "Macintosh") {
		return "Macintosh"
	} else {
		return "Unknown Device"
	}
}

func LockUserAccount(userID uint) error {
	var user models.User
	err := database.DBConn.First(&user, userID).Error
	if err != nil {
		return err
	}

	user.IsLock = 1
	return database.DBConn.Save(&user).Error
}

func UnlockUserAccount(c *fiber.Ctx) error {
	// Extract user ID from URL parameter
	userID := c.Params("id")

	var userRecord models.User
	if err := database.DBConn.First(&userRecord, userID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Unlock user account
	userRecord.IsLock = 0
	userRecord.NumberOfFailedLogin = 0
	if err := database.DBConn.Save(&userRecord).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to unlock user account"})
	}

	return c.JSON(fiber.Map{"message": "User account unlocked successfully"})
}

func HandleFailedLogin(username string) error {
	var user models.User
	err := database.DBConn.Where("user_name = ?", username).First(&user).Error
	if err != nil {
		return err
	}

	user.NumberOfFailedLogin++
	if user.NumberOfFailedLogin >= 3 {
		user.IsLock = 1
	}
	return database.DBConn.Save(&user).Error
}

// func detectSuspiciousActivity(user models.User, location string) error {
// 	if isSuspiciousLocation(location) {
// 		user.IsLocked = true
// 		return database.DBConn.Save(&user).Error
// 	}
// 	return nil
// }
