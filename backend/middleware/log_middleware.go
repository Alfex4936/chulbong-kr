package middleware

import (
	"runtime"
	"time"

	"github.com/Alfex4936/chulbong-kr/util"
	"github.com/jmoiron/sqlx"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	processMsg         = "HTTP request processed"
	insertVisitorQuery = "INSERT IGNORE INTO Visitors (IPAddress, VisitDate) VALUES (?, ?)"
)

type LogMiddleware struct {
	ChatUtil *util.ChatUtil
	DB       *sqlx.DB
}

func NewLogMiddleware(chatUtil *util.ChatUtil, db *sqlx.DB, logger *zap.Logger) *LogMiddleware {

	return &LogMiddleware{
		ChatUtil: chatUtil,
		DB:       db,
	}
}

func (l *LogMiddleware) ZapLogMiddleware(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Start timer
		start := time.Now()

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Gather request/response details
		statusCode := c.Response().StatusCode()
		method := c.Method()
		path := c.OriginalURL()

		clientIP := l.ChatUtil.GetUserIP(c)
		//visitDate := time.Now().Format("2006-01-02")

		// Perform the database logging in a goroutine
		// go func() {
		// 	if clientIP != "" {
		// 		// Check if the clientIP is a well-formed IP address
		// 		if ip := net.ParseIP(clientIP); ip != nil {
		// 			check, err := l.ChatUtil.IsIPFromSouthKorea(clientIP)
		// 			if check || err != nil {
		// 				l.DB.Exec(insertVisitorQuery, clientIP, visitDate)
		// 			}
		// 		}
		// 	}
		// }()

		userAgent := c.Get(fiber.HeaderUserAgent)
		referer := c.Get(fiber.HeaderReferer)
		queryParams := c.OriginalURL()[len(c.Path()):]

		if duration.Seconds() > util.DELAY_THRESHOLD {
			go util.SendSlackNotification(duration, statusCode, clientIP, method, path, userAgent, queryParams, referer) // Send notification in a non-blocking way
		}

		// Choose the log level and construct the log message
		level := zap.InfoLevel
		if statusCode >= 500 {
			level = zap.ErrorLevel
		} else if statusCode >= 400 {
			level = zap.WarnLevel
		}

		// Include error details if an error occurred
		var errMsg string
		if err != nil {
			errMsg = err.Error()
		}

		// Construct the structured log
		logger.Check(level, processMsg).
			Write(
				zap.Int("status", statusCode),
				zap.String("method", method),
				zap.String("path", path),
				zap.String("client_ip", clientIP),
				zap.String("user_agent", userAgent),
				zap.String("referer", referer),
				zap.Duration("duration", duration),
				zap.String("error", errMsg), // Log the error message
				// zap.Stack("stacktrace"),
				// zap.Error(err), // Include the error if present
			)

		return err
	}
}

func logRuntimeMetrics(logger *zap.Logger) {
	var memStats runtime.MemStats

	for {
		// Pause before logging again
		time.Sleep(10 * time.Minute)

		// Capture current memory stats
		runtime.ReadMemStats(&memStats)

		// Log runtime statistics
		logger.Info("Runtime metrics",
			zap.Int("goroutines", runtime.NumGoroutine()),                     // Number of goroutines
			zap.Uint64("alloc", memStats.Alloc),                               // Allocated memory
			zap.Uint64("total_alloc", memStats.TotalAlloc),                    // Total allocated memory
			zap.Uint64("sys", memStats.Sys),                                   // System memory
			zap.Uint64("heap_alloc", memStats.HeapAlloc),                      // Heap memory allocated
			zap.Uint64("heap_sys", memStats.HeapSys),                          // Heap memory in use
			zap.Uint64("heap_idle", memStats.HeapIdle),                        // Heap memory idle
			zap.Uint64("heap_inuse", memStats.HeapInuse),                      // Heap memory in use
			zap.Uint64("heap_released", memStats.HeapReleased),                // Heap memory released
			zap.Uint64("heap_objects", memStats.HeapObjects),                  // Number of heap objects
			zap.Uint64("stack_inuse", memStats.StackInuse),                    // Stack memory in use
			zap.Uint64("stack_sys", memStats.StackSys),                        // Stack memory system
			zap.Uint64("gc_sys", memStats.GCSys),                              // GC system memory
			zap.Uint64("next_gc", memStats.NextGC),                            // Next GC will happen after this amount of heap allocation
			zap.Uint32("gc_cpu_fraction", uint32(memStats.GCCPUFraction*100)), // GC CPU fraction
			zap.Uint64("last_gc", memStats.LastGC),                            // Last GC time in nanoseconds
		)
	}
}
