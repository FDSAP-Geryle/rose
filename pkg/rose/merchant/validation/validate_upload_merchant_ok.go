package utils

import (
	"fmt"
	"regexp"
	merchant "rose/pkg/rose/merchant/model"
)

func ValidateMap(data []map[string]interface{}) merchant.UploadMerchantOK {

	// Status (FAILED, DONE, ON PROGRESS)
	// Notes (Success, file empty / sheetName different, Incomplete fields)
	// User
	// TotalUpload
	// TotalSuccess
	// TotalError

	UploadMerchantOK := merchant.UploadMerchantOK{}

	for idx, row := range data {
		fmt.Printf("Validating row %d:\n", idx+1)
		for k, v := range row {
			if v == "" {
				UploadMerchantOK.Notes = "Incomplete fields"
			} else {
				UploadMerchantOK.Notes = "Success" // not sure
			}
			switch k {
			case "Phone Number":
				if PhoneNumber, ok := v.(string); !ok || isValidPhoneNumber(PhoneNumber) || PhoneNumber == "" {
					fmt.Printf(" - phone number must be more than 12 digits %v\n", v)
				} else {
					fmt.Printf(" - Valid phone number: %s\n", PhoneNumber)
				}
			case "Account Number":
				if AccountNumber, ok := v.(string); !ok || len(AccountNumber) < 16 {
					fmt.Printf(" - acount number must be 16 digits or more than 16 digits %v\n", v)
				} else {
					fmt.Printf(" - Valid Account Number: %s\n", AccountNumber)
				}
			case "email":
				if email, ok := v.(string); !ok || !isValidEmail(email) {
					fmt.Printf(" - Invalid email: %v\n", v)
				} else {
					fmt.Printf(" - Valid email: %s\n", email)
				}
			default:
				fmt.Printf(" - *****:%s\n", k)
			}
		}
		fmt.Println()
	}
	return UploadMerchantOK // return data must
}

func isValidEmail(email string) bool {
	// Basic email format check using regular expression
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(emailRegex).MatchString(email)
}
func isValidPhoneNumber(PhoneNumber string) bool {
	// Basic email format check using regular expression
	PhoneNumberRegex := `^\+?[639]\d{1,9}$`
	return regexp.MustCompile(PhoneNumberRegex).MatchString(PhoneNumber)
}
