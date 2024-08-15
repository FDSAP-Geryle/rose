package routers

import (
	"rosei/pkg/controllers/changepassword"
	"rosei/pkg/controllers/download"
	merchants "rosei/pkg/controllers/getmerchant"
	"rosei/pkg/controllers/healthchecks"
	"rosei/pkg/controllers/login"
	"rosei/pkg/controllers/logout"
	"rosei/pkg/controllers/uploadmerchant"
	models "rosei/pkg/models/merchant"
	middleware "rosei/pkg/utils"
	"time"

	fiberUtils "rosei/pkg/utils/go-utils/fiber"

	"github.com/gofiber/fiber/v2"
)

func SetupPrivateRoutes(app *fiber.App) {

	app.Use(fiberUtils.AuthenticationMiddleware(fiberUtils.JWTConfig{
		Duration:     1 * time.Hour,
		CookieMaxAge: 15 * 60,
		SetCookies:   true,
		SecretKey:    []byte(middleware.GetEnv("SECRET_KEY")),
	}))

	apiEndpoint := app.Group("/api")
	publicEndpoint := apiEndpoint.Group("/private")
	v1Endpoint := publicEndpoint.Group("/v1")

	v1Endpoint.Get("/", healthchecks.CheckServiceHealth)
	v1Endpoint.Get("/logout", logout.Logout)
	v1Endpoint.Get("/hello", func(c *fiber.Ctx) error { return c.SendString("private") })
	v1Endpoint.Get("/unlock-user/:id", login.UnlockUserAccount)
	v1Endpoint.Get("/change-password", changepassword.ChangePassword)

	v1Endpoint.Get("/merchant", merchants.Get[models.Record])
	v1Endpoint.Get("/merchantwip", merchants.Get[models.Recordwip])
	v1Endpoint.Get("/merchantok", merchants.Get[models.Recordok])

	v1Endpoint.Get("/upload-merchant", func(c *fiber.Ctx) error {
		return uploadmerchant.UploadMerchant(c, false)
	})
	v1Endpoint.Get("/upload-merchantwip", func(c *fiber.Ctx) error {
		return uploadmerchant.UploadMerchant(c, true)
	})
	v1Endpoint.Get("/upload-merchantok", uploadmerchant.UploadMerchantok)

	v1Endpoint.Get("/download/:filename", download.DownloadErr)
	v1Endpoint.Get("/download/:filename", download.DownloadSucc)
	v1Endpoint.Get("/template", download.DownloadMerchantTemplate)

}
