package register

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	models "rosei/pkg/models/user"
	"rosei/pkg/utils/go-utils/database"
	"rosei/pkg/utils/go-utils/passwordHashing"

	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {
	// Log request attempt (without sensitive data)
	log.Println("Register attempt")

	var input models.User
	if err := c.BodyParser(&input); err != nil {
		log.Println("Error parsing body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// VALIDATION

	// Username
	var validUsername = regexp.MustCompile(`^[a-zA-Z0-9_]{4,20}$`)
	if !validUsername.MatchString(input.UserName) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid username. It must be 4-20 characters long and can only contain letters, numbers, and underscores."})
	}

	// Password
	password := input.Password
	if !isValidPassword(password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid password. It must be at least 8 characters long, contain an uppercase letter, a lowercase letter, a number, and a special character."})
	}

	// Email
	var validEmail = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !validEmail.MatchString(input.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid email address."})
	}

	// Fullname
	var validFullname = regexp.MustCompile(`^[a-zA-Z]+(?:[ ][a-zA-Z]+)*$`)
	if !validFullname.MatchString(input.FullName) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid fullname. It must only contain letters and spaces."})
	}

	// Mobile number
	var validMobile = regexp.MustCompile(`^(?:\+639|09)[0-9]{9}$`)
	if !validMobile.MatchString(input.MobileNo) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid mobile number. It must be either +639 followed by 9 digits or 09 followed by 9 digits."})
	}

	// Hash & Salt
	saltedHashPassword, err := passwordHashing.HashPassword(password)
	if err != nil {
		log.Println("Error hashing password:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}
	input.Password = saltedHashPassword

	// Set password expiry date
	const passwordExpiryDuration = 90 * 24 * time.Hour // For testing, adjust as needed
	input.PwdExpiredDate = time.Now().Add(passwordExpiryDuration).UTC()

	// Create user
	if err := database.DBConn.Create(&input).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			constraintName := extractConstraintName(err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("%s already exists", constraintName)})
		}
		log.Println("Error creating user:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	// Create password history
	passwordHistory := models.PasswordHistory{
		UserId:    input.ID,
		Password:  input.Password,
		UpdatedAt: time.Now(),
	}
	if err := database.DBConn.Create(&passwordHistory).Error; err != nil {
		log.Println("Error creating password history:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to record password history"})
	}

	return c.JSON(fiber.Map{"message": "User created successfully"})
}

// Helper function to validate passwords
func isValidPassword(password string) bool {
	tests := []string{".{8,}", "[a-z]", "[A-Z]", "[0-9]", "[^\\d\\w]"}
	for _, test := range tests {
		t, _ := regexp.MatchString(test, password)
		if !t {
			return false
		}
	}
	return true
}

// Extract constraint name from error message
func extractConstraintName(errStr string) string {
	parts := strings.Split(errStr, `"`)
	if len(parts) < 3 {
		return ""
	}
	return parts[1]
}
