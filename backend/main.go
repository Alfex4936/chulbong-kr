package main

import (
	"chulbong-kr/database"
	"chulbong-kr/handlers"
	"chulbong-kr/middlewares"
	"chulbong-kr/services"
	"chulbong-kr/utils"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/storage/redis/v3"
	"github.com/gofiber/swagger"
	"github.com/gofiber/template/django/v3"
	_ "github.com/joho/godotenv/autoload"
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
	// godotenv.Overload()

	// Increase GOMAXPROCS
	runtime.GOMAXPROCS(runtime.NumCPU() * 2) // twice the number of CPUs

	redisPortStr := os.Getenv("REDIS_PORT")
	redisPort, err := strconv.Atoi(redisPortStr)
	if err != nil {
		log.Fatalf("Failed to convert REDIS_PORT to integer: %v", err)
	}

	// Initialize redis
	store := redis.New(redis.Config{
		Host:      os.Getenv("REDIS_HOST"),
		Port:      redisPort,
		Username:  os.Getenv("REDIS_USERNAME"),
		Password:  os.Getenv("REDIS_PASSWORD"),
		Database:  0,
		Reset:     true,
		TLSConfig: nil,
		PoolSize:  10 * runtime.GOMAXPROCS(0),
	})
	services.RedisStore = store

	if err := utils.LoadBadWords("badwords.txt"); err != nil {
		log.Fatalf("Failed to load bad words: %v", err)
	}

	// Initialize global variables
	setTokenExpirationTime()
	services.AWS_REGION = os.Getenv("AWS_REGION")
	services.S3_BUCKET_NAME = os.Getenv("AWS_BUCKET_NAME")
	utils.LOGIN_TOKEN_COOKIE = os.Getenv("TOKEN_COOKIE")

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
	engine := django.New("./views", ".django")

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		Prefork:       false, // Enable prefork mode for high-concurrency
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "",
		BodyLimit:     30 * 1024 * 1024, // limit to 30 MB
		IdleTimeout:   120 * time.Second,
		ReadTimeout:   10 * time.Second,
		WriteTimeout:  10 * time.Second,
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
		AppName:       "chulbong-kr",
		Concurrency:   512 * 1024,
		Views:         engine,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			// Initial status code defaults to 500
			code := fiber.StatusInternalServerError

			// Retrieve the custom status code if it's a *fiber.Error
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			// Define a user-friendly error response
			errorResponse := fiber.Map{
				"success": false,
				"message": "Something went wrong on our end. Please try again later.",
			}

			// Customize the message for known error codes
			switch code {
			case fiber.StatusNotFound: // 404
				errorResponse["message"] = "The requested resource could not be found."
			case fiber.StatusInternalServerError: // 500
				errorResponse["message"] = "An unexpected error occurred. We're working to fix the problem. Please try again later."
				// TODO: Optionally add a reference code
				// errorResponse["reference_code"] = "REF123456"
			}

			// Send a JSON response with the error message and status code
			return ctx.Status(code).JSON(errorResponse)
		},
	})
	app.Server().MaxConnsPerIP = 10

	go services.ProcessClickEventsBatch()

	logger, _ := zap.NewProduction()
	app.Use(middlewares.ZapLogMiddleware(logger))

	// Middlewares
	app.Use(healthcheck.New(healthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		LivenessEndpoint: "/",
	}))

	app.Use(encryptcookie.New(encryptcookie.Config{
		Key:    os.Getenv("ENCRYPTION_KEY"),
		Except: []string{csrf.ConfigDefault.CookieName, "Etag"}, // exclude CSRF cookie
	}))

	app.Use(etag.New(etag.Config{
		Weak: true,
	}))

	app.Use(pprof.New())

	app.Use(compress.New(compress.Config{
		Next: func(c *fiber.Ctx) bool {
			// Compress only for /api/v1/markers; return false to apply compression
			return c.Path() != "/api/v1/markers"
		},
		Level: compress.LevelBestSpeed,
	}))

	app.Use(helmet.New(helmet.Config{XSSProtection: "1; mode=block"}))
	app.Use(limiter.New(limiter.Config{
		Max:               50,
		Expiration:        45 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))
	app.Get("/metrics", middlewares.AdminOnly, monitor.New(monitor.Config{
		Title:   "chulbong-kr System Metrics",
		Refresh: time.Second * 30,
	}))
	app.Use(requestid.New())

	// Enable CORS for all routes
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173,https://chulbong-kr.vercel.app,https://www.k-pullup.com", // List allowed origins
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",                                                   // Explicitly list allowed methods
		AllowHeaders: "*",                                                                             // TODO: Allow specific headers
		// ExposeHeaders:    "Accept",
		AllowCredentials: true,
	}))

	// app.Use(logger.New())
	app.Get("/swagger/*", middlewares.AdminOnly, swagger.HandlerDefault)

	app.Use("/ws/:markerID", func(c *fiber.Ctx) error {
		// middlewares.AuthMiddleware(c)
		// if c.Locals("username") == nil {
		// 	return fiber.ErrUnauthorized
		// }

		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:markerID", websocket.New(func(c *websocket.Conn) {
		// Extract markerID from the parameter
		markerID := c.Params("markerID")
		reqID := c.Query("request-id")
		handlers.HandleChatRoomHandler(c, markerID, reqID)
	}, websocket.Config{
		// Set the handshake timeout to a reasonable duration to prevent slowloris attacks.
		HandshakeTimeout: 5 * time.Second,

		// Origins: []string{"http://localhost:8080", "https://chulbong-kr.vercel.app", "https://www.k-pullup.com"},

		EnableCompression: true,

		RecoverHandler: func(c *websocket.Conn) {
			// Custom recover logic. By default, it logs the error and stack trace.
			if r := recover(); r != nil {
				fmt.Fprintf(os.Stderr, "WebSocket panic: %v\n", r)
				debug.PrintStack()
				c.WriteMessage(websocket.CloseMessage, []byte{})
				c.Close()
			}
		},
	}))

	// HTML
	app.Get("/main", func(c *fiber.Ctx) error {
		return c.Render("login", fiber.Map{})
	})
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
		userGroup.Get("/favorites", handlers.GetFavoritesHandler)
		userGroup.Patch("/me", handlers.UpdateUserHandler)
		userGroup.Delete("/me", handlers.DeleteUserHandler)
		userGroup.Delete("/s3/objects", handlers.DeleteObjectFromS3Handler)
	}

	// Marker routes
	api.Get("/markers", handlers.GetAllMarkersHandler)

	// api.Get("/markers-addr", middlewares.AdminOnly, handlers.GetAllMarkersWithAddrHandler)
	// api.Post("/markers-addr", middlewares.AdminOnly, handlers.UpdateMarkersAddressesHandler)
	// api.Get("/markers-db", middlewares.AdminOnly, handlers.GetMarkersClosebyAdmin)

	api.Get("/markers/:markerId/details", middlewares.AuthSoftMiddleware, handlers.GetMarker)
	api.Get("/markers/close", handlers.FindCloseMarkersHandler)
	api.Get("/markers/ranking", handlers.GetMarkerRankingHandler)
	api.Get("/markers/area-ranking", handlers.GetCurrentAreaMarkerRankingHandler)
	api.Get("/markers/convert", handlers.ConvertWGS84ToWCONGNAMULHandler)
	api.Get("/markers/weather", handlers.GetWeatherByWGS84Handler)

	api.Post("/markers/upload", middlewares.AdminOnly, handlers.UploadMarkerPhotoToS3Handler)

	markerGroup := api.Group("/markers")
	{
		markerGroup.Use(middlewares.AuthMiddleware)

		markerGroup.Get("/my", handlers.GetUserMarkersHandler)
		markerGroup.Get("/:markerID/dislike-status", handlers.CheckDislikeStatus)
		markerGroup.Get("/:markerID/facilities", handlers.GetFacilitiesHandler)
		// markerGroup.Get("/:markerId", handlers.GetMarker)

		markerGroup.Post("/new", handlers.CreateMarkerWithPhotosHandler)
		markerGroup.Post("/facilities", handlers.SetMarkerFacilitiesHandler)
		markerGroup.Post("/:markerID/dislike", handlers.LeaveDislikeHandler)
		markerGroup.Post("/:markerID/favorites", handlers.AddFavoriteHandler)

		markerGroup.Put("/:markerID", handlers.UpdateMarker)

		markerGroup.Delete("/:markerID", handlers.DeleteMarkerHandler)
		markerGroup.Delete("/:markerID/dislike", handlers.UndoDislikeHandler)
		markerGroup.Delete("/:markerID/favorites", handlers.RemoveFavoriteHandler)
	}

	// Comment routes
	api.Get("/comments/:markerId/comments", handlers.LoadCommentsHandler) // no auth

	commentGroup := api.Group("/comments")
	{
		commentGroup.Use(middlewares.AuthMiddleware)
		commentGroup.Post("", handlers.PostCommentHandler)
		commentGroup.Patch("/:commentId", handlers.UpdateCommentHandler)
		commentGroup.Delete("/:commentId", handlers.RemoveCommentHandler)
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
	services.CronCleanUpPasswordTokens()
	services.CronResetClickRanking()
	services.StartOrphanedPhotosCleanupCron()

	serverAddr := fmt.Sprintf("0.0.0.0:%s", os.Getenv("SERVER_PORT"))

	// Check if the DEPLOYMENT is not local
	if os.Getenv("DEPLOYMENT") == "production" {
		// Send Slack notification
		go utils.SendDeploymentSuccessNotification("chulbong-kr", "fly.io")
	}

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
