package main

import (
	"fmt"
	"log"
	"os"
	"rose/pkg/config"
	routers "rose/pkg/routers"
	middleware "rose/pkg/utils"
	"rose/pkg/utils/go-utils/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error Loading Env File 1: ", err)
		// Since .env is not critical, we proceed without it
	}

	envi := os.Getenv("ENVIRONMENT")
	if envi == "" {
		log.Fatal("ENVIRONMENT variable not set")
	}

	envFilePath := fmt.Sprintf("project_env_files/.env-%v", envi)
	fmt.Println("Env file path: ", envFilePath)
	err = godotenv.Load(envFilePath)
	if err != nil {
		log.Fatal("Error Loading Env File 2: ", err)
	}

	secretKey := os.Getenv("SECRET_KEY")
	if envi == "" {
		log.Fatal("SECRET_KEY variable not set", secretKey)
	}

	// Create database connection
	config.CreateConnection()

	// Migrate database schema
	if err := database.DBConn.AutoMigrate(); err != nil {
		log.Fatal(err)
	}

	// Initialize DB here
	//middleware.CreateConnection()

	// Declare & initialize fiber
	app := fiber.New(fiber.Config{
		UnescapePath: true,
		BodyLimit:    32 * 1024 * 1024,
	})

	// Configure application CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "POST, GET, OPTIONS, PUT, DELETE",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		// AllowOrigins: "http://localhost:49300",
		// AllowHeaders: "Origin, Content-Type, Accept",
	}))
	// For GoRoutine implementation
	// appb := fiber.New(fiber.Config{
	// 	UnescapePath: true,
	// })

	// For GoRoutine implementation
	// appb.Use(cors.New(cors.Config{
	// 	AllowOrigins: "*",
	// 	AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	// }))

	// Declare & initialize logger
	app.Use(logger.New(logger.Config{ //Modified Logs
		//Format:     "${cyan}${time} ${white}| ${green}${status} ${white}| ${ip} | ${host} | ${method} | ${magenta}${path} ${white} | ${red}${latency} ${white} | \n${yellow}${body} ${white} | ${responseData}\n",
		Format:     "${cyan}${time} ${white}| ${green}${status} ${white}| ${ip} | ${host} | ${method} | ${magenta}${path} ${white} | ${red}${latency} ${white}\n",
		TimeFormat: "01/02/2006 3:04 PM",
		TimeZone:   "Asia/Manila",
	}))

	// For GoRoutine implementation
	// appb.Use(logger.New())

	// Declare & initialize routes
	routers.SetupPublicRoutes(app)
	routers.SetupPrivateRoutes(app)

	// For GoRoutine implementation
	// routers.SetupPublicRoutesB(appb)
	// go func() {
	// 	err := appb.Listen(fmt.Sprintf(":8002"))
	// 	if err != nil {
	// 		log.Fatal(err.Error())
	// 	}
	// }()

	//fmt.Println("Port: ", middleware.GetEnv("PORT"))
	// Serve the application
	if middleware.GetEnv("SSL") == "enabled" {
		log.Fatal(app.ListenTLS(
			fmt.Sprintf(":%s", middleware.GetEnv("PORT")),
			middleware.GetEnv("SSL_CERTIFICATE"),
			middleware.GetEnv("SSL_KEY"),
		))
	} else {
		err := app.Listen(fmt.Sprintf(":%s", middleware.GetEnv("PORT")))
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}
