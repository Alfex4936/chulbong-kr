package main

import (
	"chulbong-kr/database"
	"chulbong-kr/handlers"
	"chulbong-kr/middlewares"
	"chulbong-kr/services"
	"chulbong-kr/utils"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/Alfex4936/tzf"
	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/fgprof"
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
	"github.com/gofiber/swagger"
	"github.com/redis/rueidis"
	"go.uber.org/zap"

	// "github.com/gofiber/storage/redis/v3"

	"github.com/gofiber/template/django/v3"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"

	// amqp "github.com/rabbitmq/amqp091-go"

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
	if os.Getenv("DEPLOYMENT") != "production" {
		godotenv.Overload()
	}

	// Increase GOMAXPROCS
	runtime.GOMAXPROCS(runtime.NumCPU() * 2) // twice the number of CPUs

	setUpExternalConnections()
	setUpGlobals()

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
		Views:         django.New("./views", ".django"),
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
	// app.Server().MaxConnsPerIP = 10

	// Middlewares
	setUpMiddlewares(app)

	// API
	websocketConfig := websocket.Config{
		// Set the handshake timeout to a reasonable duration to prevent slowloris attacks.
		HandshakeTimeout: 5 * time.Second,

		Origins: []string{"https://test.k-pullup.com", "https://www.k-pullup.com"},

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
	}

	app.Get("/ws/:markerID", func(c *fiber.Ctx) error {
		// Extract markerID from the parameter
		markerID := c.Params("markerID")
		reqID := c.Query("request-id")

		// Use GetBanDetails to check if the user is banned and get the remaining ban time
		banned, remainingTime, err := services.WsRoomManager.GetBanDetails(markerID, reqID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
		}
		if banned {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":         "User is banned",
				"remainingTime": remainingTime.Seconds(), // Respond with remaining time in seconds
			})
		}

		// Proceed with WebSocket upgrade if not banned
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	}, websocket.New(func(c *websocket.Conn) {
		// Extract markerID from the parameter again if necessary
		markerID := c.Params("markerID")
		reqID := c.Query("request-id")

		// Now, the connection is already upgraded to WebSocket, and passed the ban check.
		handlers.HandleChatRoomHandler(c, markerID, reqID)
	}, websocketConfig))

	// HTML
	app.Get("/main", func(c *fiber.Ctx) error {
		return c.Render("login", fiber.Map{})
	})

	// Setup MY routes
	api := app.Group("/api/v1")
	handlers.RegisterAdminRoutes(api)
	handlers.RegisterAuthRoutes(api)
	handlers.RegisterUserRoutes(api)
	handlers.RegisterMarkerRoutes(api)
	handlers.RegisterCommentRoutes(api)
	handlers.RegisterTossPaymentRoutes(api)
	handlers.RegisterReportRoutes(api)
	handlers.RegisterNotificationRoutes(api, websocketConfig)

	// Cron jobs
	services.RunAllCrons()

	// Server settings
	serverAddr := fmt.Sprintf("0.0.0.0:%s", os.Getenv("SERVER_PORT"))

	// Check if the DEPLOYMENT is not local
	if os.Getenv("DEPLOYMENT") == "production" {
		// Send Slack notification
		go utils.SendDeploymentSuccessNotification(app.Config().AppName, "fly.io")

		// Random ranking
		go services.ResetAndRandomizeClickRanking()
	} else {
		log.Printf("There are %d APIs available in chulbong-kr", countAPIs(app))
	}

	// Start the Fiber app
	if err := app.Listen(serverAddr); err != nil {
		panic(err)
	}
}

func setUpMiddlewares(app *fiber.App) {
	logger, _ := zap.NewProduction()
	app.Use(middlewares.ZapLogMiddleware(logger))

	// Middlewares
	app.Use(healthcheck.New(healthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			log.Printf("---- %s", utils.CreateAnonymousID(c))
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
	app.Use(fgprof.New())

	app.Use(compress.New(compress.Config{
		// Next: func(c *fiber.Ctx) bool {
		// 	// Compress only for /api/v1/markers; return false to apply compression
		// 	return c.Path() != "/api/v1/markers"
		// },
		Level: compress.LevelBestSpeed,
	}))

	app.Use(helmet.New(helmet.Config{XSSProtection: "1; mode=block"}))
	app.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			// Skip rate limiting for /users/logout and /users/me
			path := c.Path()
			if path == "/api/v1/auth/logout" || path == "/api/v1/users/me" {
				return true // Returning true skips the limiter
			}
			return false // Apply the limiter for all other paths
		},

		KeyGenerator: func(c *fiber.Ctx) string {
			return utils.GetUserIP(c)
		},
		Max:               60,
		Expiration:        1 * time.Minute,
		LimiterMiddleware: limiter.SlidingWindow{},
		LimitReached: func(c *fiber.Ctx) error {
			// Custom response when rate limit is exceeded
			c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
			c.Status(429).SendString("Too many requests, please try again later.")
			return nil
		},
		SkipFailedRequests: true,
	}))
	app.Get("/metrics", middlewares.AdminOnly, monitor.New(monitor.Config{
		Title:   "chulbong-kr System Metrics",
		Refresh: time.Second * 30,
	}))
	app.Use(requestid.New())

	// Enable CORS for all routes
	app.Use(cors.New(cors.Config{
		// AllowOrigins: "http://localhost:5173,https://chulbong-kr.vercel.app,https://www.k-pullup.com", // List allowed origins
		AllowOriginsFunc: func(origin string) bool {
			// Check if the origin is a subdomain of k-pullup.com
			return strings.HasSuffix(origin, ".k-pullup.com") || origin == "https://www.k-pullup.com" || origin == "https://chulbong-kr.vercel.app" || origin == "http://localhost:5173"
		},
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS", // Explicitly list allowed methods
		AllowHeaders: "*",                           // TODO: Allow specific headers
		// ExposeHeaders:    "Accept",
		AllowCredentials: true,
	}))

	// app.Use(logger.New())
	app.Get("/swagger/*", middlewares.AdminOnly, swagger.HandlerDefault)
}

func setUpExternalConnections() {
	// Initialize database connection
	if err := database.Connect(); err != nil {
		panic(err)
	}

	// Initialize redis
	rdb, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress:       []string{os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")},
		Username:          os.Getenv("REDIS_USERNAME"),
		Password:          os.Getenv("REDIS_PASSWORD"),
		DisableCache:      true, // dragonfly doesn't support CACHING command
		SelectDB:          0,
		ForceSingleClient: true,
		// PoolSize:    10 * runtime.GOMAXPROCS(0),
		// MaxRetries:  5,
		// PipelineMultiplex: 2, // Default is typically sufficient
		// BlockingPoolSize:  5,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	})
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}

	ctx := context.Background()
	// Build the Ping command
	pingCmd := rdb.B().Ping().Build()

	// Execute the Ping command
	err = rdb.Do(ctx, pingCmd).Error()
	if err != nil {
		log.Fatalf("Error pinging to Redis: %v", err)
	}

	// Start the Redis health check routine
	go redisHealthCheck(rdb)

	if os.Getenv("DEPLOYMENT") == "production" {
		// Flush the Redis database to clear all keys
		err := rdb.Do(ctx, rdb.B().Flushall().Build()).Error()
		if err != nil {
			log.Fatalf("Error executing FLUSHALL SYNC: %v", err)
		}
	}

	services.RedisStore = rdb

	// geminiClient, err := genai.NewClient(context.Background(), option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// services.Gemini = geminiClient

	// log.Println(services.ChatGemini())

	// Message Broker
	// connection, err := amqp.Dial(os.Getenv("LAVINMQ_HOST"))
	// if err != nil {
	// 	log.Panicf("Failed to connect to LavinMQ")
	// }
	// services.LavinMQClient = connection

}

func setUpGlobals() {
	finder, err := tzf.NewDefaultFinder()
	if err != nil {
		log.Fatalf("Failed to initialize timezone finder: %v", err)
	}
	utils.TimeZoneFinder = finder

	if err := utils.LoadBadWords("badwords.txt"); err != nil {
		log.Fatalf("Failed to load bad words: %v", err)
	}

	// Initialize global variables
	setTokenExpirationTime()
	services.AWS_REGION = os.Getenv("AWS_REGION")
	services.S3_BUCKET_NAME = os.Getenv("AWS_BUCKET_NAME")
	utils.LOGIN_TOKEN_COOKIE = os.Getenv("TOKEN_COOKIE")
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

// countAPIs counts the number of APIs in a Fiber app
func countAPIs(app *fiber.App) int {
	numAPIs := 0
	for _, route := range app.GetRoutes(true) {
		// Check if the route is for an API (skip middleware routes)
		if route.Path[len(route.Path)-1] != '*' {
			numAPIs++
		}
	}
	return numAPIs
}

func redisHealthCheck(rdb rueidis.Client) {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for t := range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		err := rdb.Do(ctx, rdb.B().Ping().Build()).Error()
		cancel()

		if err != nil {
			log.Printf("Redis ping failed at %v: %v", t, err)
			// Attempt to reconnect
			reconnectRedis(rdb)
		} else {
			log.Printf("Redis ping success at %v", t)
		}
	}
}
func reconnectRedis(rdb rueidis.Client) {
	for i := 0; i < 3; i++ { // Try reconnecting 3 times
		time.Sleep(time.Duration(i+1) * time.Second) // Exponential back-off strategy
		rdb.Close()
		rdb, err := rueidis.NewClient(rueidis.ClientOption{
			InitAddress:  []string{os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")},
			Username:     os.Getenv("REDIS_USERNAME"),
			Password:     os.Getenv("REDIS_PASSWORD"),
			DisableCache: true,
			TLSConfig:    &tls.Config{InsecureSkipVerify: true},
		})
		if err == nil {
			services.RedisStore = rdb
			log.Println("Reconnected to Redis successfully")
			return
		}
		log.Printf("Failed to reconnect to Redis: %v", err)
	}

	log.Fatal("Failed to reconnect to Redis after several attempts")
}
