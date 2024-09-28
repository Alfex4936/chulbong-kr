package util

import (
	"bytes"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
)

// JsonBodyParser binds the request body to a json struct with optimizations for performance.
func JsonBodyParser(c *fiber.Ctx, out interface{}) error {
	// Retrieve the Content-Type header as a byte slice
	contentType := c.Request().Header.ContentType()

	// If Content-Type is empty, return an error
	if len(contentType) == 0 {
		return fiber.ErrUnprocessableEntity
	}

	// Convert Content-Type to lower-case in place to handle case-insensitivity
	ToLower(contentType)

	// Parse vendor-specific Content-Type (e.g., application/problem+json -> application/json)
	parsedCType := ParseVendorSpecificContentType(contentType)

	// Remove any parameters from Content-Type (e.g., application/json; charset=utf-8 -> application/json)
	if semiColonIndex := bytes.IndexByte(parsedCType, ';'); semiColonIndex != -1 {
		parsedCType = parsedCType[:semiColonIndex]
	}

	// Check if the Content-Type ends with "json"
	if !bytes.HasSuffix(parsedCType, []byte("json")) {
		return fiber.ErrUnprocessableEntity
	}

	return c.App().Config().JSONDecoder(c.Body(), out)
}

// JsonBodyParserFast binds the request body to a json struct with optimizations for performance.
// using sonic.ConfigFastest
func JsonBodyParserFast(c *fiber.Ctx, out interface{}) error {
	// Retrieve the Content-Type header as a byte slice
	contentType := c.Request().Header.ContentType()

	// If Content-Type is empty, return an error
	if len(contentType) == 0 {
		return fiber.ErrUnprocessableEntity
	}

	// Convert Content-Type to lower-case in place to handle case-insensitivity
	ToLower(contentType)

	// Parse vendor-specific Content-Type (e.g., application/problem+json -> application/json)
	parsedCType := ParseVendorSpecificContentType(contentType)

	// Remove any parameters from Content-Type (e.g., application/json; charset=utf-8 -> application/json)
	if semiColonIndex := bytes.IndexByte(parsedCType, ';'); semiColonIndex != -1 {
		parsedCType = parsedCType[:semiColonIndex]
	}

	// Check if the Content-Type ends with "json"
	if !bytes.HasSuffix(parsedCType, []byte("json")) {
		return fiber.ErrUnprocessableEntity
	}
	return sonic.ConfigFastest.Unmarshal(c.Body(), out)
	// return c.App().Config().JSONDecoder(c.Body(), out)
}

// ToLower converts an ASCII byte slice to lower-case in place.
func ToLower(b []byte) {
	for i := 0; i < len(b); i++ {
		if 'A' <= b[i] && b[i] <= 'Z' {
			b[i] += 'a' - 'A'
		}
	}
}

// ParseVendorSpecificContentType efficiently parses vendor-specific content types.
// It transforms types like "application/problem+json" to "application/json".
func ParseVendorSpecificContentType(cType []byte) []byte {
	plusIndex := bytes.IndexByte(cType, '+')
	if plusIndex == -1 {
		return cType
	}

	semiColonIndex := bytes.IndexByte(cType, ';')
	var parsableType []byte

	if semiColonIndex == -1 {
		parsableType = cType[plusIndex+1:]
	} else if plusIndex < semiColonIndex {
		parsableType = cType[plusIndex+1 : semiColonIndex]
	} else {
		return cType[:semiColonIndex]
	}

	slashIndex := bytes.IndexByte(cType, '/')
	if slashIndex == -1 {
		return cType
	}

	// Create a new slice to hold "application/json"
	parsed := make([]byte, slashIndex+1+len(parsableType))
	copy(parsed, cType[:slashIndex+1])
	copy(parsed[slashIndex+1:], parsableType)
	return parsed
}
