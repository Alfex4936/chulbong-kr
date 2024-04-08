package utils

import (
	"fmt"
	"math/rand"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var adjectives = []string{
	"귀여운",     // Cute
	"멋진",      // Cool
	"착한",      // Kind
	"용감한",     // Brave
	"영리한",     // Clever
	"재미있는",    // Fun
	"행복한",     // Happy
	"사랑스러운",   // Lovely
	"기운찬",     // Energetic
	"빛나는",     // Shining
	"평화로운",    // Peaceful
	"신비로운",    // Mysterious
	"자유로운",    // Free
	"매력적인",    // Charming
	"섬세한",     // Delicate
	"우아한",     // Elegant
	"활발한",     // Lively
	"강인한",     // Strong
	"독특한",     // Unique
	"무서운",     // Scary
	"꿈꾸는",     // Dreamy
	"느긋한",     // Relaxed
	"열정적인",    // Passionate
	"소중한",     // Precious
	"신선한",     // Fresh
	"창의적인",    // Creative
	"우수한",     // Excellent
	"재치있는",    // Witty
	"감각적인",    // Sensual
	"흥미로운",    // Interesting
	"유명한",     // Famous
	"현명한",     // Wise
	"대담한",     // Bold
	"침착한",     // Calm
	"신속한",     // Swift
	"화려한",     // Gorgeous
	"정열적인",    // Passionate (Alternate)
	"끈기있는",    // Persistent
	"애정이 깊은",  // Affectionate
	"민첩한",     // Agile
	"빠른",      // Quick
	"조용한",     // Quiet
	"명랑한",     // Cheerful
	"정직한",     // Honest
	"용서하는",    // Forgiving
	"용기있는",    // Courageous
	"성실한",     // Sincere
	"호기심이 많은", // Curious
	"겸손한",     // Humble
	"관대한",     // Generous
}

// 9 names
var names = []string{
	"라이언", // Ryan
	"어피치", // Apeach
	"콘",   // Con
	"무지",  // Muzi
	"네오",  // Neo
	"프로도", // Frodo
	"제이지", // Jay-G
	"튜브",  // Tube
	"철봉",  // chulbong
}

// GenerateKoreanNickname generates random user nickname
func GenerateKoreanNickname() string {

	// Select a random
	adjective := adjectives[rand.Intn(len(adjectives))]

	name := names[rand.Intn(len(names))]

	// Generate a unique identifier
	uid := uuid.New().String()

	// Use the first 8 characters of the UUID to keep it short
	shortUID := uid[:8]

	// possibilities for conflict
	// highly unlikely.
	// 25 * 9 * 16^8 (UUID first 8 characters)
	// UUID can conflict by root(16*8) = 65,536
	return fmt.Sprintf("%s %s [%s]", adjective, name, shortUID)
}

func GetUserIP(c *fiber.Ctx) string {
	clientIP := c.Get("Fly-Client-IP")
	if clientIP == "" {
		clientIP = c.Get("Fly-Client-Ip")
	}
	// If Fly-Client-IP is not found, fall back to X-Forwarded-For
	if clientIP == "" {
		clientIP = c.Get("X-Forwarded-For")
	}
	if clientIP == "" {
		clientIP = c.Get("X-Real-IP")
	}

	// If X-Forwarded-For is also empty, use c.IP() as the last resort
	if clientIP == "" {
		clientIP = c.IP()
	}
	return clientIP
}
