package middleware

import (
	"strings"
	"time"

	"github.com/Alfex4936/chulbong-kr/util"
	"github.com/gofiber/contrib/fgprof"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func SetupMiddlewares(
	app *fiber.App,
	logger *zap.Logger,
	chatUtil *util.ChatUtil,
	authMiddleware *AuthMiddleware,
	zapMiddleware *LogMiddleware,
) {
	// Middlewares
	app.Use(zapMiddleware.ZapLogMiddleware(logger))
	app.Use(healthcheck.New(healthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			// logger.Debug("----", zap.String("user_id", chatUtil.CreateAnonymousID(c)))
			return true
		},
		LivenessEndpoint: "/",
	}))

	app.Use(etag.New(etag.Config{
		Weak: true,
	}))

	app.Use(encryptcookie.New(encryptcookie.Config{
		Key:    viper.GetString("ENCRYPTION_KEY"),
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
	app.Use(helmet.New(helmet.Config{XSSProtection: "1; mode=block", CrossOriginEmbedderPolicy: "credentialless"}))

	app.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			// Skip rate limiting for /users/logout and /users/me
			path := c.Path()
			// path == "/metrics" ||
			if path == "/api/v1/auth/logout" || path == "/api/v1/users/me" || path == "/api/v1/markers/save-offline" {
				return true // Returning true skips this limiter
			}
			return false // Apply the limiter for all other paths
		},

		KeyGenerator: func(c *fiber.Ctx) string {
			return chatUtil.GetUserIP(c)
		},
		Max:               60,
		Expiration:        1 * time.Minute,
		LimiterMiddleware: limiter.SlidingWindow{},
		// LimiterMiddleware: middleware.SlidingWindow{},
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
	// app.Get("/metrics", authMiddleware.CheckAdmin, monitor.New(monitor.Config{
	// 	Title:   "chulbong-kr System Metrics",
	// 	Refresh: time.Second * 1,
	// 	// ChartJsURL: "https://cdn.jsdelivr.net/npm/chart.js", // TODO: update to v4
	// 	Next: nil,
	// 	CustomHead: `
	// 	<style>
	// 		body {
	// 			margin: 0;
	// 			font: 16px / 1.6 'Roboto', sans-serif;
	// 			background-color: #f0f0f0;
	// 		}
	// 		.wrapper {
	// 			max-width: 900px;
	// 			margin: 0 auto;
	// 			padding: 30px 0;
	// 		}
	// 		.title {
	// 			text-align: center;
	// 			margin-bottom: 2em;
	// 		}
	// 		.title h1 {
	// 			font-size: 1.8em;
	// 			padding: 0;
	// 			margin: 0;
	// 		}
	// 		.row {
	// 			display: flex;
	// 			margin-bottom: 20px;
	// 			align-items: center;
	// 		}
	// 		.row .column:first-child { width: 35%; }
	// 		.row .column:last-child { width: 65%; }
	// 		.metric {
	// 			color: #777;
	// 			font-weight: 900;
	// 		}
	// 		h2 {
	// 			padding: 0;
	// 			margin: 0;
	// 			font-size: 2.2em;
	// 		}
	// 		h2 span {
	// 			font-size: 12px;
	// 			color: #777;
	// 		}
	// 		h2 span.ram_os { color: rgba(255, 150, 0, .8); }
	// 		h2 span.ram_total { color: rgba(0, 200, 0, .8); }
	// 		canvas {
	// 			width: 200px;
	// 			height: 180px;
	// 		}
	// 	</style>
	// 	<style type="text/css">
	// 		/* Chart.js */
	// 		@keyframes chartjs-render-animation{from{opacity:.99}to{opacity:1}}.chartjs-render-monitor{animation:chartjs-render-animation 1ms}
	// 		.chartjs-size-monitor,.chartjs-size-monitor-expand,.chartjs-size-monitor-shrink{
	// 			position:absolute;direction:ltr;left:0;top:0;right:0;bottom:0;overflow:hidden;pointer-events:none;visibility:hidden;z-index:-1
	// 		}
	// 		.chartjs-size-monitor-expand>div{position:absolute;width:1000000px;height:1000000px;left:0;top:0}
	// 		.chartjs-size-monitor-shrink>div{position:absolute;width:200%;height:200%;left:0;top:0}
	// 	</style>
	// `,
	// }))
	app.Use(requestid.New(requestid.Config{Header: "", Generator: util.FastLogID}))

	// Enable CORS for all routes
	app.Use(cors.New(cors.Config{
		// AllowOrigins: "http://localhost:5173,https://chulbong-kr.vercel.app,https://www.k-pullup.com", // List allowed origins
		AllowOriginsFunc: func(origin string) bool {
			// Check if the origin is a subdomain of k-pullup.com
			return strings.HasSuffix(origin, ".k-pullup.com") || origin == "https://www.k-pullup.com" || origin == "http://localhost:5173" || origin == "https://ligne-j-train.vercel.app"
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
				logger.Info("Blocked access attempt to:", zap.String("path", c.Path()))

				// or perhaps a 403 Forbidden
				return c.Status(fiber.StatusForbidden).SendString("Access forbidden, saving your information to server disk...: " + c.IP())
			}
		}

		// Proceed with the next middleware if the path is not blocked
		return c.Next()
	})

	app.Get("/swagger/*", authMiddleware.CheckAdmin, swagger.HandlerDefault)
}
