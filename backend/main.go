package main

import (
	"chulbong-kr/database"
	"chulbong-kr/handlers"
	"chulbong-kr/middlewares"
	"chulbong-kr/services"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize global variables
	setTokenExpirationTime()
	services.AWS_REGION = os.Getenv("AWS_REGION")
	services.S3_BUCKET_NAME = os.Getenv("AWS_BUCKET_NAME")
	middlewares.TOKEN_COOKIE = os.Getenv("TOKEN_COOKIE")

	// Initialize database connection
	if err := database.Connect(); err != nil {
		panic(err)
	}

	// OAuth2 Configuration
	conf := &oauth2.Config{
		ClientID:     os.Getenv("G_CLIENT_ID"),
		ClientSecret: os.Getenv("G_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("G_REDIRECT"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		Prefork:       true, // Enable prefork mode for high-concurrency
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "",
		BodyLimit:     10 * 1024 * 1024, // limit to 10 MB
		IdleTimeout:   120 * time.Second,
		ReadTimeout:   10 * time.Second,
		WriteTimeout:  10 * time.Second,
	})

	app.Server().MaxConnsPerIP = 100

	// Enable CORS for all routes
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",                           // Allows all origins
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS", // Explicitly list allowed methods
		AllowHeaders:     "*",                           // TODO: Allow specific headers
		ExposeHeaders:    "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers",
		AllowCredentials: true,
	}))

	app.Use(logger.New())

	// Setup routes
	api := app.Group("/api/v1")

	api.Get("/google", handlers.GetGoogleAuthHandler(conf))

	// Authentication routes
	authGroup := api.Group("/auth")
	{
		authGroup.Post("/signup", handlers.SignUpHandler)
		authGroup.Post("/login", handlers.LoginHandler)
		authGroup.Get("/google/callback", handlers.GetGoogleCallbackHandler(conf))
	}

	// User routes
	userGroup := api.Group("/users")
	{
		userGroup.Use(middlewares.AuthMiddleware)
		userGroup.Delete("/me", handlers.DeleteUserHandler)
		userGroup.Delete("/s3/objects", handlers.DeleteObjectFromS3Handler)
	}

	// Marker routes
	api.Get("/markers", handlers.GetAllMarkersHandler)
	markerGroup := api.Group("/markers")
	{
		markerGroup.Use(middlewares.AuthMiddleware)
		markerGroup.Post("/new", handlers.CreateMarkerWithPhotosHandler)
		markerGroup.Get("/:id", handlers.GetMarker)
		markerGroup.Put("/:id", handlers.UpdateMarker)
		markerGroup.Post("/upload", handlers.UploadMarkerPhotoToS3Handler)
		markerGroup.Delete("/:markerID", handlers.DeleteMarkerHandler)
	}

	// Comment routes
	commentGroup := api.Group("/comments")
	{
		commentGroup.Use(middlewares.AuthMiddleware)
		commentGroup.Post("/", handlers.PostComment)
		commentGroup.Put("/:commentId", handlers.UpdateComment)
		commentGroup.Delete("/:commentId", handlers.DeleteComment)
	}

	app.Get("/example-get", handlers.GetExample)
	app.Put("/example-put", handlers.PutExample)
	app.Delete("/example-delete", handlers.DeleteExample)
	app.Post("/example-post", handlers.PostExample)
	app.Get("/example/:string/:id", handlers.DynamicRouteExample)
	app.Get("/example-optional/:param?", handlers.QueryParamsExample)

	// Cron jobs
	services.CronCleanUpToken()
	services.StartOrphanedPhotosCleanupCron()

	serverAddr := fmt.Sprintf("0.0.0.0:%s", os.Getenv("SERVER_PORT"))

	// Start the Fiber app
	if err := app.Listen(serverAddr); err != nil {
		panic(err)
	}
}

func setTokenExpirationTime() {
	// Get the token expiration interval from the environment variable
	durationStr := os.Getenv("TOKEN_EXPIRATION_INTERVAL")

	// Convert the duration from string to int
	durationInt, err := strconv.Atoi(durationStr)
	if err != nil {
		log.Fatalf("Error converting TOKEN_EXPIRATION_INTERVAL to int: %v", err)
	}

	// Assign the converted duration to the global variable
	services.TOKEN_DURATION = time.Duration(durationInt) * time.Hour
}
