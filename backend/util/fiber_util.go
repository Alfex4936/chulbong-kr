package util

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"

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

// JsonBodyParserStrict parses the JSON body into the provided struct and ensures no extra fields are present.
func JsonBodyParserStrict(c *fiber.Ctx, out interface{}) error {
	// Retrieve and validate Content-Type
	contentType := c.Get("Content-Type")
	if contentType == "" || !strings.HasPrefix(strings.ToLower(contentType), "application/json") {
		return fiber.ErrUnprocessableEntity
	}

	// Read the body
	body := c.Body()

	// Unmarshal into a map to check for unknown fields
	var tempMap map[string]interface{}
	if err := sonic.ConfigFastest.Unmarshal(body, &tempMap); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}

	// Get the expected fields from the struct tags
	expectedFields, err := getJSONTags(out)
	if err != nil {
		return fmt.Errorf("failed to get JSON tags: %w", err)
	}

	// Check for unknown fields
	for key := range tempMap {
		if _, ok := expectedFields[key]; !ok {
			return fmt.Errorf("unknown field: %s", key)
		}
	}

	// Unmarshal into the actual struct
	if err := sonic.ConfigFastest.Unmarshal(body, out); err != nil {
		return fmt.Errorf("failed to unmarshal JSON into struct: %w", err)
	}

	return nil
}

// getJSONTags extracts the JSON tags from the struct fields.
func getJSONTags(obj interface{}) (map[string]struct{}, error) {
	t := reflect.TypeOf(obj)
	if t.Kind() != reflect.Ptr {
		return nil, errors.New("expected pointer to struct")
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return nil, errors.New("expected pointer to struct")
	}

	tags := make(map[string]struct{})
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("json")
		if tag == "-" || tag == "" {
			continue
		}
		// Handle omitempty and other tag options
		tagName := strings.Split(tag, ",")[0]
		tags[tagName] = struct{}{}
	}
	return tags, nil
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
