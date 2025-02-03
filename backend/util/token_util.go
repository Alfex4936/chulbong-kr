package util

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"io"
	mrand "math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Alfex4936/chulbong-kr/config"
	"github.com/gofiber/fiber/v2"
)

const (
	HexDigits      = "0123456789abcdef"
	letters        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	lettersLen     = len(letters)
	tokenBufSize   = 64 // Optimized buffer size for tokens
	randStrBufSize = 64 // Optimized buffer size for random strings
)

var (
	// Pool for base64/opaque tokens
	tokenBytePool = sync.Pool{
		New: func() interface{} {
			b := make([]byte, tokenBufSize)
			return &b
		},
	}
	// Pool for random strings.
	randomStringBytePool = sync.Pool{
		New: func() interface{} {
			b := make([]byte, randStrBufSize)
			return &b
		},
	}

	counter uint64
)

// Initialize the counter or other seed once
func init() {
	// read 8 random bytes into a seed
	var seed [8]byte
	_, _ = io.ReadFull(rand.Reader, seed[:])
	// Use it for a starting offset
	counter = binary.BigEndian.Uint64(seed[:])
}

type TokenUtil struct {
	Config *config.AppConfig
}

func NewTokenUtil(config *config.AppConfig) *TokenUtil {
	return &TokenUtil{Config: config}
}

// GenerateOpaqueToken generates a URL-safe base64 encoded random token.
func (t *TokenUtil) GenerateOpaqueToken(length int) (string, error) {
	if length <= 0 {
		return "", nil
	}

	// Calculate the needed byte length for the desired string length.
	// Since base64 encoding uses 4 characters for every 3 bytes, and we use RawURLEncoding which omits padding,
	// we can determine the byte length needed to get a specific encoded length.
	// We round up to ensure we always meet or exceed the desired length.
	rawLen := (length*3 + 3) / 4

	// Retrieve a byte slice from the pool.
	bufPtr := tokenBytePool.Get().(*[]byte)
	defer tokenBytePool.Put(bufPtr)

	// Resize buffer from the pool, only if necessary.
	if cap(*bufPtr) < rawLen {
		*bufPtr = make([]byte, rawLen)
	}
	buf := (*bufPtr)[:rawLen]

	// Fill with cryptographic random bytes.
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		return "", err
	}

	// Encode to a base64 string using the buffer directly.
	encodedLen := base64.RawURLEncoding.EncodedLen(len(buf))
	if cap(*bufPtr) < encodedLen {
		*bufPtr = make([]byte, encodedLen)
	}
	encodedBuf := (*bufPtr)[:encodedLen]

	base64.RawURLEncoding.Encode(encodedBuf, buf)
	return BytesToString(encodedBuf), nil
}

// GenerateRandomString creates a random string using optimized methods.
func (t *TokenUtil) GenerateRandomString(n int) string {
	if n <= 0 {
		return ""
	}

	// Retrieve a byte slice from the pool.
	bufPtr := randomStringBytePool.Get().(*[]byte)
	defer randomStringBytePool.Put(bufPtr)

	// Resize buffer only if necessary.
	if cap(*bufPtr) < n {
		*bufPtr = make([]byte, n)
	}
	buf := (*bufPtr)[:n]

	// Use batched random number generation.
	for i, cache, remain := n-1, mrand.Int63(), lettersLen; i >= 0; {
		if remain == 0 {
			cache, remain = mrand.Int63(), lettersLen
		}
		if idx := int(cache & int64(lettersLen-1)); idx < lettersLen {
			buf[i] = letters[idx]
			i--
		}
		cache >>= 6
		remain--
	}

	return BytesToString(buf)
}

func (t *TokenUtil) GenerateLoginCookie(value string) fiber.Cookie {
	cookie := fiber.Cookie{
		Name:     t.Config.LoginTokenCookie,
		Value:    value,                              // The token generated for the user
		Expires:  time.Now().Add(24 * 7 * time.Hour), // Set the cookie to expire in 7 days
		HTTPOnly: true,                               // Ensure the cookie is not accessible through client-side scripts
		Secure:   t.Config.IsProduction == "production",
		SameSite: fiber.CookieSameSiteLaxMode,
		Path:     "/", // Scope of the cookie
	}

	if t.Config.IsProduction == "production" {
		cookie.Domain = ".k-pullup.com" // Allow cookie to be shared across all subdomains
	}
	return cookie
}

func (t *TokenUtil) ClearLoginCookie() fiber.Cookie {
	cookie := fiber.Cookie{
		Name:     t.Config.LoginTokenCookie,
		Value:    "",              // The token generated for the user
		Expires:  time.Unix(0, 0), // Set the cookie to be expired in the past
		HTTPOnly: true,            // Ensure the cookie is not accessible through client-side scripts
		Secure:   t.Config.IsProduction == "production",
		SameSite: fiber.CookieSameSiteLaxMode,
		Path:     "/", // Scope of the cookie
	}

	if t.Config.IsProduction == "production" {
		cookie.Domain = ".k-pullup.com" // Allow cookie to be shared across all subdomains
	}
	return cookie
}

// FastLogID combines the current time and a unique counter into 16 bytes, then hex-encodes it
func FastLogID() string {
	// 1) read time (64 bits)
	now := uint64(time.Now().UnixNano())

	// 2) increment a shared counter (64 bits)
	c := atomic.AddUint64(&counter, 1)

	// 16 bytes total => 128 bits
	// [0..7]: time, [8..15]: counter
	var buf [16]byte
	binary.BigEndian.PutUint64(buf[:8], now)
	binary.BigEndian.PutUint64(buf[8:], c)

	// encode to 32 hex digits
	var dst [32]byte
	encodeHex(&dst, &buf)

	return BytesToString(dst[:])
}

// encodeHex is an unrolled loop converting each byte to 2 hex chars, in place.
func encodeHex(dst *[32]byte, src *[16]byte) {
	for i := 0; i < 16; i++ {
		b := src[i]
		dst[i*2] = HexDigits[b>>4]
		dst[i*2+1] = HexDigits[b&0x0F]
	}
}
