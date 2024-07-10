package routers

import (
	controllers "rose/pkg/controllers"
	"rose/pkg/controllers/healthchecks"
	"rose/pkg/routers/routes"
	"rose/pkg/utils/go-utils/encryptDecrypt"

	"github.com/gofiber/fiber/v2"
)

func SetupPublicRoutes(app *fiber.App) {

	// Endpoints
	apiEndpoint := app.Group("/api")
	publicEndpoint := apiEndpoint.Group("/public")
	v1Endpoint := publicEndpoint.Group("/v1")

	// Service health check
	v1Endpoint.Get("/", healthchecks.CheckServiceHealth)

	//MyTest routes
	v1Endpoint.Get("/hello", func(c *fiber.Ctx) error { return c.SendString("hello mond") })
	v1Endpoint.Get("/merchant_upload", routes.GetUploadMerchant)
	v1Endpoint.Post("/merchant_upload_non_activated", routes.UploadMerchantNonActivated)

}

func SetupPublicRoutesB(app *fiber.App) {
	// Endpoints
	apiEndpoint := app.Group("/api")
	publicEndpoint := apiEndpoint.Group("/public")
	v1Endpoint := publicEndpoint.Group("/v1")

	// Service health check
	v1Endpoint.Get("/", healthchecks.CheckServiceHealthB)

	// Utility
	utilityEndpoint := v1Endpoint.Group("utility")
	utilityEndpoint.Post("/test/encrypt", encryptDecrypt.EncryptHandler)
	utilityEndpoint.Post("/test/decrypt", encryptDecrypt.DecryptHandler)

	// SMS
	smsEndpoint := v1Endpoint.Group("sms")
	smsEndpoint.Get("/test/get/all/sms/type", controllers.SMSBlastingContentType)

	// ACTIVITY
	activityEndpoint := v1Endpoint.Group("activity")
	activityEndpoint.Post("/isPalindrome", controllers.PalindromeHandler)

}
