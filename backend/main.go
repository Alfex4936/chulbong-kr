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

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	_ "chulbong-kr/docs"
)

// @title			chulbong-kr API
// @version		1.0
// @description	Pullup bar locations with KakaoMap API
// @contact.name	API Support
// @contact.email	chulbong.kr@gmail.com
// @license.name	MIT
// @license.url	https://github.com/Alfex4936/chulbong-kr/blob/main/LICENSE
// @host			localhost:9452
// @BasePath		/api/v1/
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

	// engine := html.New("./views", ".html")

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
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
		// Views:         engine,
	})
	app.Server().MaxConnsPerIP = 100

	app.Static("/toss/", "./public")

	// Middlewares
	app.Use(healthcheck.New(healthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		LivenessEndpoint: "/",
	}))
	logger, _ := zap.NewProduction()

	app.Use(middlewares.ZapLogMiddleware(logger))

	app.Use(helmet.New())
	app.Use(limiter.New(limiter.Config{
		Max:               100,
		Expiration:        30 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))
	app.Get("/metrics", middlewares.AdminOnly, monitor.New(monitor.Config{
		Title:   "chulbong-kr System Metrics",
		Refresh: time.Second * 60,
	}))
	app.Use(requestid.New())

	// Enable CORS for all routes
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173,https://chulbong-kr.vercel.app,https://developers.tosspayments.com", // List allowed origins
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",                                                              // Explicitly list allowed methods
		AllowHeaders:     "*",                                                                                        // TODO: Allow specific headers
		ExposeHeaders:    "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers",
		AllowCredentials: true,
	}))

	// app.Use(logger.New())
	app.Get("/swagger/*", middlewares.AdminOnly, swagger.HandlerDefault)

	// Setup routes
	api := app.Group("/api/v1")

	api.Get("/google", handlers.GetGoogleAuthHandler(conf))
	api.Get("/admin", middlewares.AdminOnly, func(c *fiber.Ctx) error { return c.JSON("good") })

	// Authentication routes
	authGroup := api.Group("/auth")
	{
		authGroup.Post("/signup", handlers.SignUpHandler)
		authGroup.Post("/login", handlers.LoginHandler)
		authGroup.Post("/logout", middlewares.AuthMiddleware, handlers.LogoutHandler)
		authGroup.Get("/google/callback", handlers.GetGoogleCallbackHandler(conf))
		authGroup.Post("/verify-email/send", handlers.SendVerificationEmailHandler)
		authGroup.Post("/verify-email/confirm", handlers.ValidateTokenHandler)

		// Finding password
		authGroup.Post("/request-password-reset", handlers.RequestResetPasswordHandler)
		authGroup.Post("/reset-password", handlers.ResetPasswordHandler)
	}

	// User routes
	userGroup := api.Group("/users")
	{
		userGroup.Use(middlewares.AuthMiddleware)
		userGroup.Get("/me", handlers.ProfileHandler)
		userGroup.Patch("/me", handlers.UpdateUserHandler)
		userGroup.Delete("/me", handlers.DeleteUserHandler)
		userGroup.Delete("/s3/objects", handlers.DeleteObjectFromS3Handler)
	}

	// Marker routes
	api.Get("/markers", handlers.GetAllMarkersHandler)

	markerGroup := api.Group("/markers")
	{
		markerGroup.Use(middlewares.AuthMiddleware)
		markerGroup.Get("/my", handlers.GetUserMarkersHandler)
		// markerGroup.Get("/:markerID", handlers.GetMarker)
		markerGroup.Get("/close", handlers.FindCloseMarkersHandler)
		markerGroup.Get("/:markerID/dislike-status", handlers.CheckDislikeStatus)
		markerGroup.Post("/new", handlers.CreateMarkerWithPhotosHandler)
		markerGroup.Post("/upload", handlers.UploadMarkerPhotoToS3Handler)
		markerGroup.Post("/:markerID/dislike", handlers.LeaveDislikeHandler)
		markerGroup.Put("/:markerID", handlers.UpdateMarker)
		markerGroup.Delete("/:markerID", handlers.DeleteMarkerHandler)
		markerGroup.Delete("/:markerID/dislike", handlers.UndoDislikeHandler)
	}

	// Comment routes
	commentGroup := api.Group("/comments")
	{
		commentGroup.Use(middlewares.AuthMiddleware)
		commentGroup.Post("/", handlers.PostComment)
		commentGroup.Put("/:commentId", handlers.UpdateComment)
		commentGroup.Delete("/:commentId", handlers.DeleteComment)
	}

	tossGroup := api.Group("/payments/toss")
	{
		tossGroup.Post("/confirm", handlers.ConfirmToss)
		// tossGroup.Get("/success", handlers.SuccessToss)
		// tossGroup.Get("/fail", handlers.FailToss)
	}

	// app.Get("/example-optional/:param?", handlers.QueryParamsExample)

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
