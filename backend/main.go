package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	configfx "github.com/Alfex4936/chulbong-kr/configfx"
	servicefx "github.com/Alfex4936/chulbong-kr/di"
	"github.com/Alfex4936/chulbong-kr/handler"
	"github.com/Alfex4936/chulbong-kr/middleware"
	"github.com/Alfex4936/chulbong-kr/service"
	"github.com/Alfex4936/chulbong-kr/util"
	"github.com/Alfex4936/tzf"
	"github.com/blevesearch/bleve/v2"
	_ "github.com/blevesearch/bleve/v2/analysis/analyzer/custom"
	_ "github.com/blevesearch/bleve/v2/analysis/char/html"
	_ "github.com/blevesearch/bleve/v2/analysis/lang/cjk"
	_ "github.com/blevesearch/bleve/v2/analysis/token/edgengram"
	_ "github.com/blevesearch/bleve/v2/analysis/token/ngram"
	_ "github.com/blevesearch/bleve/v2/analysis/token/unicodenorm"
	_ "github.com/blevesearch/bleve/v2/analysis/tokenizer/unicode"
	_ "github.com/blevesearch/bleve/v2/index/upsidedown/store/boltdb"
	"github.com/dgraph-io/ristretto"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"

	_ "github.com/go-sql-driver/mysql"
	"github.com/goccy/go-json"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/rueidis"

	ristretto_store "github.com/eko/gocache/store/ristretto/v4"

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
	"github.com/gofiber/template/django/v3"

	amqp "github.com/rabbitmq/amqp091-go"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ඞ Fiber app constructor
func NewFiberApp(
	logger *zap.Logger,
	chatUtil *util.ChatUtil,
	wsConfig websocket.Config,
	markerHandler *handler.MarkerHandler,
	userHandler *handler.UserHandler,
	searchHandler *handler.SearchHandler,
	adminHandler *handler.AdminHandler,
	authHandler *handler.AuthHandler,
	chatHandler *handler.ChatHandler,
	commentHandler *handler.CommentHandler,
	notificatinHandler *handler.NotificationHandler,
	authMiddleware *middleware.AuthMiddleware,
	zapMiddleware *middleware.LogMiddleware) *fiber.App {
	app := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "nginx",
		BodyLimit:     30 * 1024 * 1024, // limit to 30 MB
		IdleTimeout:   120 * time.Second,
		ReadTimeout:   10 * time.Second,
		WriteTimeout:  10 * time.Second,
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
		AppName:       "chulbong-kr",
		Concurrency:   512 * 1024,
		Views:         django.New("./view", ".django"),
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

	// Middlewares
	app.Use(zapMiddleware.ZapLogMiddleware(logger))
	app.Use(healthcheck.New(healthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			// log.Printf("---- %s", chatUtil.CreateAnonymousID(c))
			return true
		},
		LivenessEndpoint: "/",
	}))

	app.Use(etag.New(etag.Config{
		Weak: true,
	}))

	app.Use(encryptcookie.New(encryptcookie.Config{
		Key:    os.Getenv("ENCRYPTION_KEY"),
		Except: []string{csrf.ConfigDefault.CookieName, "Etag"}, // exclude CSRF cookie
	}))

	app.Use("/debug/pprof", authMiddleware.CheckAdmin, pprof.New())
	app.Use("/debug/fgprof", authMiddleware.CheckAdmin, fgprof.New())

	app.Use(compress.New(compress.Config{
		// Next: func(c *fiber.Ctx) bool {
		// 	// Compress only for /api/v1/markers; return false to apply compression
		// 	return c.Path() != "/api/v1/markers"
		// },
		Level: compress.LevelBestSpeed,
	}))

	// ContentSecurityPolicy: "default-src 'self'; script-src 'self' 'unsafe-inline'; object-src 'none'"
	app.Use(helmet.New(helmet.Config{XSSProtection: "1; mode=block"}))
	app.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			// Skip rate limiting for /users/logout and /users/me
			path := c.Path()
			if path == "/api/v1/auth/logout" || path == "/api/v1/users/me" || path == "/api/v1/search/marker" {
				return true // Returning true skips the limiter
			}
			return false // Apply the limiter for all other paths
		},

		KeyGenerator: func(c *fiber.Ctx) string {
			return chatUtil.GetUserIP(c)
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
		SkipFailedRequests: false,
	}))

	// loginLimiter := limiter.New(limiter.Config{
	// 	KeyGenerator: func(c *fiber.Ctx) string {
	// 		return chatUtil.GetUserIP(c)
	// 	},
	// 	Max:               5, // Stricter limit for login
	// 	Expiration:        30 & time.Second,
	// 	LimiterMiddleware: limiter.SlidingWindow{},
	// 	LimitReached: func(c *fiber.Ctx) error {
	// 		c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
	// 		return c.Status(429).SendString("Too many login attempts, please try again later.")
	// 	},
	// 	SkipFailedRequests: false,
	// })

	// TODO: v3 not yet
	app.Get("/metrics", authMiddleware.CheckAdmin, monitor.New(monitor.Config{
		Title:   "chulbong-kr System Metrics",
		Refresh: time.Second * 30,
		Next:    nil,
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

	app.Use(func(c *fiber.Ctx) error {
		// List of paths to block
		blockedPaths := []string{
			"/mysql/scripts/setup.php",
			"/phpMyAdmin2/scripts/setup.php",
			"/phpma/scripts/setup.php",
			"/sqlweb/scripts/setup.php",
			"/dbadmin/scripts/setup.php",
		}

		// Check if the requested path is in the blocked paths list
		for _, path := range blockedPaths {
			if c.Path() == path {
				// Log the attempt for monitoring purposes
				log.Println("Blocked access attempt to:", c.Path())

				// You could return a 404 Not Found, or perhaps a 403 Forbidden
				return c.Status(fiber.StatusForbidden).SendString("Access forbidden, saving your information to server disk...: " + c.IP())
			}
		}

		// Proceed with the next middleware if the path is not blocked
		return c.Next()
	})

	app.Get("/swagger/*", authMiddleware.CheckAdmin, swagger.HandlerDefault)

	// Set up routes
	api := app.Group("/api/v1")
	handler.RegisterMarkerRoutes(api, markerHandler, authMiddleware)
	handler.RegisterReportRoutes(api, markerHandler, authMiddleware)
	handler.RegisterUserRoutes(api, userHandler, authMiddleware)
	handler.RegisterSearchRoutes(api, searchHandler)
	handler.RegisterAdminRoutes(api, adminHandler, authMiddleware)
	handler.RegisterAuthRoutes(api, authHandler, authMiddleware)
	handler.RegisterChatRoutes(app, wsConfig, chatHandler) // not /api/v1/
	handler.RegisterCommentRoutes(api, commentHandler, authMiddleware)
	handler.RegisterNotificationRoutes(app, wsConfig, notificatinHandler, authMiddleware) // not /api/v1/

	return app
}

// Provides a new logger instance.
func NewLogger() (*zap.Logger, error) {
	return zap.NewProduction()
}

func NewHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second, // Set a timeout to avoid hanging requests indefinitely
	}
}

func NewGeminiClient() (*genai.Client, error) {
	return genai.NewClient(context.Background(), option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
}

func NewLavinMqClient() (*amqp.Connection, error) {
	return amqp.Dial(os.Getenv("LAVINMQ_HOST"))
}

func NewGoCacheLocalStorage() (*ristretto_store.RistrettoStore, error) {
	estimatedMarkers := 10000
	approxMarkerSize := 100 // bytes
	maxCost := estimatedMarkers * approxMarkerSize

	ristrettoCache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 10,
		MaxCost:     int64(maxCost),
		BufferItems: 64,
	})
	if err != nil {
		return nil, err
	}
	ristrettoStore := ristretto_store.NewRistretto(ristrettoCache)

	return ristrettoStore, nil
}

// NewDatabase sets up the database connection
func NewDatabase() (*sqlx.DB, error) {
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUsername, dbPassword, dbHost, dbPort, dbName)
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}
	return db, nil
}

func NewRedis(lifecycle fx.Lifecycle) (*service.RedisClient, error) {
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

	if os.Getenv("DEPLOYMENT") == "production" {
		// Flush the Redis database to clear all keys
		err := rdb.Do(context.Background(), rdb.B().Flushall().Build()).Error()
		if err != nil {
			log.Fatalf("Error executing FLUSHALL SYNC: %v", err)
		}
	}

	safeClient := &service.RedisClient{Client: rdb}

	// Register lifecycle hooks for Redis
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				ticker := time.NewTicker(30 * time.Minute)
				defer ticker.Stop()

				for range ticker.C {
					safeClient.Mu.RLock()
					err := pingRedis(safeClient.Client)
					safeClient.Mu.RUnlock()

					if err != nil {
						log.Println("Redis ping failed, attempting to reconnect...")
						newClient, err := reconnectRedis()
						if err != nil {
							log.Fatalf("Failed to reconnect: %v", err)
						}
						safeClient.Reconnect(newClient)
					}
				}
			}()
			return nil
		},
		OnStop: func(context.Context) error {
			rdb.Close()
			return nil
		},
	})

	return safeClient, nil
}

func NewWsConfig() websocket.Config {
	return websocket.Config{
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
}

// Load timezone finder
func NewTimeZoneFinder() (tzf.F, error) {
	finder, err := tzf.NewDefaultFinder()
	if err != nil {
		log.Fatalf("Error loading timezone finder: %v", err)
		return &tzf.DefaultFinder{}, err
	}
	return finder, nil
}

func NewBleveIndex() (bleve.Index, error) {
	return bleve.Open("markers.bleve")
}

// MAIN Fx
func main() {
	// Load environment variables from a .env file if not in production
	if os.Getenv("DEPLOYMENT") != "production" {
		godotenv.Overload()
	}

	// Set GOMAXPROCS to twice the number of logical CPUs
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

	// Create an Fx application with provided dependencies and lifecycle hooks
	fx.New(
		servicefx.FxMarkerModule,
		servicefx.FxExternalModle,
		servicefx.FxChatModule,
		servicefx.FxUtilModule,
		servicefx.FxUserModule,
		servicefx.FxAPIModule,
		servicefx.FxFacadeModule,

		configfx.FxConfigModule,

		fx.Provide(
			NewHTTPClient,
			NewLogger,
			NewDatabase,
			NewRedis,
			NewWsConfig,
			NewTimeZoneFinder,
			NewBleveIndex,
			NewGoCacheLocalStorage,
			// NewGeminiClient,
			// NewLavinMqClient,

			middleware.NewAuthMiddleware,
			middleware.NewLogMiddleware,

			NewFiberApp,
		),
		fx.Invoke(
			registerHooks,
			util.RegisterBadWordUtilLifecycle,
			service.RegisterSchedulerLifecycle,
			util.RegisterPdfInitLifecycle,
			service.RegisterMarkerLifecycle,
		), // func(diGraph fx.DotGraph) {
		// 	log.Println("➡️", diGraph)
		// }

	).Run()
}

// registerHooks sets up lifecycle hooks for starting and stopping the Fiber server
func registerHooks(lc fx.Lifecycle,
	app *fiber.App, db *sqlx.DB, logger *zap.Logger,
	rankService *service.MarkerRankService,
	redisService *service.RedisService,
	safeClient *service.RedisClient,
) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			serverAddr := fmt.Sprintf("0.0.0.0:%s", os.Getenv("SERVER_PORT"))

			logger.Info("💖 Starting Fiber v2 server...")
			go func() {
				if os.Getenv("DEPLOYMENT") == "production" {
					// Send Slack notification
					go util.SendDeploymentSuccessNotification(app.Config().AppName, "fly.io")
					// Random ranking
					go rankService.ResetAndRandomizeClickRanking()
				} else {
					log.Printf("There are %d APIs available in chulbong-kr", countAPIs(app))
				}

				if err := app.Listen(serverAddr); err != nil {
					logger.Fatal("Failed to start Fiber v3", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(context.Context) error {
			logger.Info("=== Shutting down Fiber v2 server...")
			if err := db.Close(); err != nil {
				logger.Error("Failed to close database connection", zap.Error(err))
			}
			safeClient.Mu.Lock()
			defer safeClient.Mu.Unlock()
			if safeClient.Client != nil {
				safeClient.Client.Close()
			}

			return app.Shutdown()
		},
	})
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

func pingRedis(rdb rueidis.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return rdb.Do(ctx, rdb.B().Ping().Build()).Error()
}

func reconnectRedis() (rueidis.Client, error) {
	var newRdb rueidis.Client
	for i := 0; i < 3; i++ {
		time.Sleep(time.Duration(i+1) * time.Second)
		var err error
		newRdb, err = rueidis.NewClient(rueidis.ClientOption{
			InitAddress:  []string{os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")},
			Username:     os.Getenv("REDIS_USERNAME"),
			Password:     os.Getenv("REDIS_PASSWORD"),
			DisableCache: true,
			TLSConfig:    &tls.Config{InsecureSkipVerify: true},
		})
		if err == nil {
			return newRdb, nil
		}
		log.Printf("Attempt %d to reconnect failed with error: %v", i+1, err)
	}
	return nil, fmt.Errorf("failed to reconnect to Redis after several attempts")
}
