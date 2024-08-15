package main

import (
	"fmt"
	"log"
	"os"
	"rosei/pkg/config"
	routers "rosei/pkg/routers"
	middleware "rosei/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {

	LoadEnv()
	config.CreateConnection()

	app := fiber.New(fiber.Config{
		UnescapePath: true,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Use(logger.New())

	routers.SetupPublicRoutes(app)
	routers.SetupPrivateRoutes(app)

	PORT := middleware.GetEnv("PORT")
	fmt.Println("Port: ", PORT)
	if middleware.GetEnv("SSL") == "enabled" {
		log.Fatal(app.ListenTLS(
			fmt.Sprintf(":%s", PORT),
			middleware.GetEnv("SSL_CERTIFICATE"),
			middleware.GetEnv("SSL_KEY"),
		))
	} else {
		err := app.Listen(fmt.Sprintf(":%s", PORT))
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Failed to Load Environment File .env")
	}

	envi := os.Getenv("ENVIRONMENT")
	envFilePath := fmt.Sprintf("env/.env-%v", envi)
	fmt.Println("Environment:", envi)

	err = godotenv.Load(envFilePath)
	if err != nil {
		log.Fatal("Failed to Load Environment File ", envFilePath)
	}
}
