package controllers

import (
	"strings"
	"unicode"

	"github.com/gofiber/fiber/v2"
)

// PalindromeRequest represents the JSON request structure for palindrome check.
type PalindromeRequest struct {
	PalindromeWord string `json:"palindromeWord"`
}

// PalindromeResponse represents the JSON response structure for palindrome check.
type PalindromeResponse struct {
	PalindromeWord string `json:"palindromeWord"`
	IsPalindrome   bool   `json:"isPalindrome"`
}

// isPalindrome checks if the given string is a palindrome
func isPalindrome(s string) bool {
	s = strings.ToLower(s)

	isAlphanumeric := func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsNumber(r)
	}

	filtered := strings.Map(func(r rune) rune {
		if isAlphanumeric(r) {
			return r
		}
		return -1
	}, s)

	return isStringPalindrome(filtered)
}

// isStringPalindrome checks if a given string is a palindrome
func isStringPalindrome(s string) bool {
	length := len(s)
	for i := 0; i < length/2; i++ {
		if s[i] != s[length-1-i] {
			return false
		}
	}
	return true
}

// PalindromeHandler handles POST requests to check if a string is a palindrome
func PalindromeHandler(c *fiber.Ctx) error {
	var palindromeWord PalindromeRequest
	if err := c.BodyParser(&palindromeWord); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	isPalindrome := isPalindrome(palindromeWord.PalindromeWord)

	res := PalindromeResponse{
		PalindromeWord: palindromeWord.PalindromeWord,
		IsPalindrome:   isPalindrome,
	}

	return c.JSON(res)
}
