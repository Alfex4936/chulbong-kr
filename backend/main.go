package main

import (
	"chulbong-kr/database"
	"chulbong-kr/handlers"
	"chulbong-kr/middlewares"
	"chulbong-kr/services"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database connection
	if err := database.Connect(); err != nil {
		panic(err)
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		Prefork:       true, // Enable prefork mode for high-concurrency
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "",
		BodyLimit:     10 * 1024 * 1024, // limit to 10 MB
	})

	// Enable CORS for all routes
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Use(logger.New())

	// Setup routes
	api := app.Group("/api/v1")

	// Authentication routes
	authGroup := api.Group("/auth")
	{
		authGroup.Post("/signup", handlers.SignUpHandler)
		authGroup.Post("/login", handlers.LoginHandler)
	}

	// Marker routes
	markerGroup := api.Group("/markers")
	{
		markerGroup.Use(middlewares.AuthMiddleware)
		markerGroup.Post("/", handlers.CreateMarkerHandler)
		markerGroup.Get("/:id", handlers.GetMarker)
		markerGroup.Put("/:id", handlers.UpdateMarker)
		markerGroup.Post("/upload", handlers.UploadMarkerPhotoToS3Handler)
	}

	// Group routes under /api
	apiGroup := api.Group("/api")
	{
		apiGroup.Use(middlewares.AuthMiddleware)
		apiGroup.Get("/", func(c *fiber.Ctx) error {
			email := c.Locals("email").(string)

			return c.JSON(email)
		})
	}

	app.Get("/example-get", handlers.GetExample)
	app.Put("/example-put", handlers.PutExample)
	app.Delete("/example-delete", handlers.DeleteExample)
	app.Post("/example-post", handlers.PostExample)
	app.Get("/example/:string/:id", handlers.DynamicRouteExample)
	app.Get("/example-optional/:param?", handlers.QueryParamsExample)

	// Cron jobs
	services.CronCleanUpToken()

	// Start the Fiber app
	if err := app.Listen("0.0.0.0:9452"); err != nil {
		panic(err)
	}
}
