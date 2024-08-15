package routers

import (
	"rosei/pkg/controllers/healthchecks"
	"rosei/pkg/controllers/login"
	"rosei/pkg/controllers/register"

	"github.com/gofiber/fiber/v2"
)

func SetupPublicRoutes(app *fiber.App) {

	// Endpoints
	apiEndpoint := app.Group("/api")
	publicEndpoint := apiEndpoint.Group("/public")
	v1Endpoint := publicEndpoint.Group("/v1")

	// Service health check
	v1Endpoint.Get("/", healthchecks.CheckServiceHealth)
	v1Endpoint.Get("/login", login.Login)
	v1Endpoint.Get("/register", register.Register)
}
