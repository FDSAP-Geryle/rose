package uploadmerchant

import (
	"log"
	"regexp"
	models "rosei/pkg/models/merchant"
	"rosei/pkg/utils/go-utils/database"
	"strings"
)

func ValidateMap(ExcelData []map[string]interface{}) (result ValidationResult, NewExcelData []map[string]interface{}) {
	totalSuccess, totalFailed := 0, 0

	for _, row := range ExcelData {
		isValid := true

		for k, v := range row {
			switch k {
			case "Phone Number":
				if !isValidPhoneNumber(v) {
					isValid = false
				}
			case "email":
				if !isValidEmail(v) {
					isValid = false
				}
			}
		}

		if isValid {
			totalSuccess++
		} else {
			totalFailed++
		}
	}

	result = ValidationResult{
		TotalUpload:  len(ExcelData),
		TotalSuccess: totalSuccess,
		TotalFailed:  totalFailed,
		Notes:        "Validation completed",
	}

	return result, ExcelData
}

func isValidPhoneNumber(phoneNumber interface{}) bool {
	phoneNumberStr, ok := phoneNumber.(string)
	if !ok || phoneNumberStr == "" {
		log.Println("Invalid or empty phone number")
		return false
	}

	var record models.ReceiveUploadActivated
	err := database.DBConn.First(&record, "phone_number = ?", phoneNumberStr).Error
	if err == nil {
		log.Println("Phone number already exists in the database")
		return false
	} else if !strings.Contains(err.Error(), "record not found") {
		log.Printf("Database error: %v", err)
		return false
	}

	phoneNumberRegex := `^\+?[639]\d{9,12}$`
	return regexp.MustCompile(phoneNumberRegex).MatchString(phoneNumberStr)
}

func isValidEmail(email interface{}) bool {
	emailStr, ok := email.(string)
	if !ok || emailStr == "" {
		log.Println("Invalid or empty email")
		return false
	}

	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(emailRegex).MatchString(emailStr)
}
