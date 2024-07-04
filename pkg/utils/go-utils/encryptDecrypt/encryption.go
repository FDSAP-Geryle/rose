package encryptDecrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"os"
	"rose/pkg/models/response"
	models "rose/pkg/models/struct"

	"github.com/gofiber/fiber/v2"
)

var iv = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 13, 05}

// const secretKey string = "abc&1*~#^2^#s0^=)^^7%b34"
// Encrypt method is to encrypt or hide any classified text
func EncryptHandler(c *fiber.Ctx) error {
	// Parse the request body into an EncryptDecryptRequest struct
	encryptRequest := new(models.EncryptDecryptRequest)
	if err := c.BodyParser(encryptRequest); err != nil {
		return err
	}

	// Get the secret key from the environment variable
	secretKey := os.Getenv("SECRET_KEY")

	// Encrypt all fields using the provided secret key
	encryptedHost, err := Encrypt(encryptRequest.Host, secretKey)
	if err != nil {
		return err
	}
	encryptedDbName, err := Encrypt(encryptRequest.DbName, secretKey)
	if err != nil {
		return err
	}
	encryptedUsername, err := Encrypt(encryptRequest.Username, secretKey)
	if err != nil {
		return err
	}
	encryptedPassword, err := Encrypt(encryptRequest.Password, secretKey)
	if err != nil {
		return err
	}

	// Construct the response
	response := response.EncryptDecyrptResponseModel{
		RetCode: "200",
		Message: "Fields Encrypted Successfully",
		Data: []models.EncryptDecryptRequest{
			{
				Host:     encryptedHost,
				DbName:   encryptedDbName,
				Username: encryptedUsername,
				Password: encryptedPassword,
			},
		},
	}

	// Return the response
	return c.JSON(response)
}

func Encrypt(text, secretKey string) (string, error) {
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, iv)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return encodeBase64(cipherText), nil
}

func encodeBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// Decrypt method is to extract back the encrypted text
func DecryptHandler(c *fiber.Ctx) error {
	// Parse the request body into a User struct
	decryptRequest := new(models.EncryptDecryptRequest)
	if err := c.BodyParser(decryptRequest); err != nil {
		return err
	}

	secretKey := os.Getenv("SECRET_KEY")

	if err := c.BodyParser(&decryptRequest); err != nil {
		return err
	}

	// Decrypt all fields using the provided secret key
	decryptedHost, err := Decrypt(decryptRequest.Host, secretKey)
	if err != nil {
		return err
	}
	decryptedDbName, err := Decrypt(decryptRequest.DbName, secretKey)
	if err != nil {
		return err
	}
	decryptedUsername, err := Decrypt(decryptRequest.Username, secretKey)
	if err != nil {
		return err
	}
	decryptedPassword, err := Decrypt(decryptRequest.Password, secretKey)
	if err != nil {
		return err
	}

	// Construct the response
	// response := response.EncryptDecyrptResponseModel{
	// 	RetCode: "200",
	// 	Message: "Fields Decrypted Successfully",
	// 	Data: []models.EncryptDecryptRequest{
	// 		{
	// 			Host:     decryptedHost,
	// 			DbName:   decryptedDbName,
	// 			Username: decryptedUsername,
	// 			Password: decryptedPassword,
	// 		},
	// 	},
	// }

	// Return the response
	return c.Status(200).JSON(response.EncryptDecyrptResponseModel{
		RetCode: "200",
		Message: "Fields Decrypted Successfully",
		Data: []models.EncryptDecryptRequest{
			{
				Host:     decryptedHost,
				DbName:   decryptedDbName,
				Username: decryptedUsername,
				Password: decryptedPassword,
			},
		},
	})
}

func Decrypt(text, secretKey string) (string, error) {
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return "", err
	}
	cipherText := decodeBase64(text)
	cfb := cipher.NewCFBDecrypter(block, iv)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}

func decodeBase64(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}
