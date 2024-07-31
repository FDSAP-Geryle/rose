package routers

import (
	"rose/pkg/rose/healthchecks"
	"rose/pkg/rose/merchant/controller"

	"github.com/gofiber/fiber/v2"
)

func SetupPublicRoutes(app *fiber.App) {

	// Endpoints
	apiEndpoint := app.Group("/api")
	publicEndpoint := apiEndpoint.Group("/public")
	v1Endpoint := publicEndpoint.Group("/v1")

	// Service health check
	v1Endpoint.Get("/", healthchecks.CheckServiceHealth)

	//rose routes
	v1Endpoint.Get("/merchant_upload", controller.GetUploadMerchant)
	v1Endpoint.Post("/merchants/upload", controller.UploadMerchantOK)
	v1Endpoint.Get("merchants/download/:filename", controller.DownloadMerchant)
	v1Endpoint.Get("merchants/template", controller.DownloadMerchantTemplate)

}
