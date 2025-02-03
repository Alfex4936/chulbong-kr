// Package util provides encoding and decoding functions for Blurhash.
package util

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/rwcarlsen/goexif/exif"
)

const characters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz#$%*+,-.:;=?@[]^_{|}~"

var (
	sRGBToLinearTable  [256]float32
	linearToSRGBTable  [4096]int
	decodeCharacterMap [128]int
)

func init() {
	// Initialize sRGB to Linear lookup table
	for i := 0; i < 256; i++ {
		v := float32(i) / 255
		if v <= 0.04045 {
			sRGBToLinearTable[i] = v / 12.92
		} else {
			sRGBToLinearTable[i] = float32(math.Pow(float64((v+0.055)/1.055), 2.4))
		}
	}

	// Initialize Linear to sRGB lookup table
	for i := 0; i < 4096; i++ {
		v := float64(i) / 4095
		if v <= 0.0031308 {
			linearToSRGBTable[i] = int(v*12.92*255 + 0.5)
		} else {
			linearToSRGBTable[i] = int((1.055*math.Pow(v, 1/2.4)-0.055)*255 + 0.5)
		}
	}

	// Initialize decode character map
	for i := 0; i < 128; i++ {
		decodeCharacterMap[i] = -1
	}
	for i, c := range characters {
		decodeCharacterMap[c] = i
	}
}

// EncodeBlurHashImage encodes an image.Image into a Blurhash string.
func EncodeBlurHashImage(img image.Image, xComponents, yComponents int) string {
	// Attempt to directly use *image.RGBA if possible.
	rgba, ok := img.(*image.RGBA)
	if !ok {
		// Convert image to RGBA if not already
		bounds := img.Bounds()
		w := bounds.Dx()
		h := bounds.Dy()
		tmp := image.NewRGBA(bounds)
		dstPix := tmp.Pix
		stride := tmp.Stride

		// Convert to RGBA in one pass
		for y := 0; y < h; y++ {
			off := y * stride
			for x := 0; x < w; x++ {
				c := img.At(x+bounds.Min.X, y+bounds.Min.Y)
				r, g, b, a := c.RGBA()
				// Convert from 16-bit to 8-bit
				dstPix[off+0] = uint8(r >> 8)
				dstPix[off+1] = uint8(g >> 8)
				dstPix[off+2] = uint8(b >> 8)
				dstPix[off+3] = uint8(a >> 8)
				off += 4
			}
		}
		rgba = tmp
	}

	width := rgba.Bounds().Dx()
	height := rgba.Bounds().Dy()
	stride := rgba.Stride
	pix := rgba.Pix

	// Extract RGB data into a continuous slice to pass to Encode.
	// single allocation
	pixels := make([]uint8, width*height*3)
	idx := 0
	for y := 0; y < height; y++ {
		rowStart := y * stride
		for x := 0; x < width; x++ {
			p := rowStart + x*4
			pixels[idx+0] = pix[p]
			pixels[idx+1] = pix[p+1]
			pixels[idx+2] = pix[p+2]
			idx += 3
		}
	}

	bytesPerRow := width * 3
	return EncodeBlurHash(xComponents, yComponents, width, height, pixels, bytesPerRow)
}

func EncodeBlurHashImageWithMeta(img image.Image, xComponents, yComponents int, extension string, orientation int) string {
	blurhash := EncodeBlurHashImage(img, xComponents, yComponents)
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	// Format: width|height|orientation|ext|hash
	return fmt.Sprintf("%d!%d!%d!%s!%s", w, h, orientation, extension, blurhash)
}

func DecodeBlurHashWithMeta(extendedHash string, punch int) ([]uint8, int, int, int, string, error) {
	parts := strings.SplitN(extendedHash, "!", 5)
	if len(parts) < 5 {
		return nil, 0, 0, 1, "", errors.New("invalid extended blurhash format")
	}

	width, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, 0, 0, 1, "", err
	}
	height, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, 0, 0, 1, "", err
	}
	orientation, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, 0, 0, 1, "", err
	}
	extension := parts[3]
	hash := parts[4]

	pixels, err := DecodeBlurHash(hash, width, height, punch)
	if err != nil {
		return nil, 0, 0, 1, "", err
	}

	return pixels, width, height, orientation, extension, nil
}

// DecodeBlurHash generates raw RGB pixel data from a Blurhash string.
func DecodeBlurHash(hash string, width, height, punch int) ([]uint8, error) {
	if !IsValidBlurhash(hash) {
		return nil, errors.New("invalid blurhash")
	}

	sizeFlag, err := decodeInt(hash[:1])
	if err != nil {
		return nil, err
	}
	numY := (sizeFlag / 9) + 1
	numX := (sizeFlag % 9) + 1

	expectedLength := 4 + 2*numX*numY
	if len(hash) != expectedLength {
		return nil, errors.New("invalid blurhash length")
	}

	quantizedMaxValue, err := decodeInt(hash[1:2])
	if err != nil {
		return nil, err
	}
	maxValue := (float64(quantizedMaxValue) + 1) / 166.0

	colors := make([][3]float64, numX*numY)
	// Decode DC & AC components
	dcValue, err := decodeInt(hash[2:6])
	if err != nil {
		return nil, err
	}
	colors[0][0], colors[0][1], colors[0][2] = decodeDC(dcValue)

	for i := 1; i < numX*numY; i++ {
		start := 4 + i*2
		val, err := decodeInt(hash[start : start+2])
		if err != nil {
			return nil, err
		}
		colors[i][0], colors[i][1], colors[i][2] = decodeAC(val, maxValue*float64(punch))
	}

	// Pre-split colors into separate arrays for faster index-based access
	rColors := make([]float64, numX*numY)
	gColors := make([]float64, numX*numY)
	bColors := make([]float64, numX*numY)
	for i := 0; i < numX*numY; i++ {
		rColors[i] = colors[i][0]
		gColors[i] = colors[i][1]
		bColors[i] = colors[i][2]
	}

	// Precompute cosines
	cosX := precomputeCosinesFloat64(width, numX)
	cosY := precomputeCosinesFloat64(height, numY)

	pixels := make([]uint8, width*height*3)

	// Decode loop with optimizations
	for y := 0; y < height; y++ {
		cosYComponent := cosY[y]
		yOffset := y * width * 3
		for x := 0; x < width; x++ {
			cosXComponent := cosX[x]

			var r, g, b float64
			// Combine loops efficiently
			for j := 0; j < numY; j++ {
				cy := cosYComponent[j]
				jBase := j * numX
				for i := 0; i < numX; i++ {
					c := cosXComponent[i] * cy
					idx := jBase + i
					r += rColors[idx] * c
					g += gColors[idx] * c
					b += bColors[idx] * c
				}
			}

			// Inline linearToSRGB logic
			if r < 0 {
				r = 0
			} else if r > 1 {
				r = 1
			}
			ri := linearToSRGBTable[int(r*4095)]

			if g < 0 {
				g = 0
			} else if g > 1 {
				g = 1
			}
			gi := linearToSRGBTable[int(g*4095)]

			if b < 0 {
				b = 0
			} else if b > 1 {
				b = 1
			}
			bi := linearToSRGBTable[int(b*4095)]

			pixelIndex := yOffset + x*3
			pixels[pixelIndex] = uint8(ri)
			pixels[pixelIndex+1] = uint8(gi)
			pixels[pixelIndex+2] = uint8(bi)
		}
	}

	return pixels, nil
}

// EncodeBlurHash generates a Blurhash string for the given pixel data.
func EncodeBlurHash(xComponents, yComponents, width, height int, pixels []uint8, bytesPerRow int) string {
	if xComponents < 1 || xComponents > 9 || yComponents < 1 || yComponents > 9 {
		return ""
	}

	// Precompute cosine values using float32
	cosX := precomputeCosines(width, xComponents)
	cosY := precomputeCosines(height, yComponents)

	// Initialize factors as a flat slice
	factorsCount := xComponents * yComponents
	// Separate arrays for R, G, B factors
	factorsR := make([]float32, factorsCount)
	factorsG := make([]float32, factorsCount)
	factorsB := make([]float32, factorsCount)

	// Convert all sRGB pixels to linear float32 upfront
	totalPixels := width * height * 3
	linearPixels := make([]float32, totalPixels)
	for i := 0; i < totalPixels; i++ {
		linearPixels[i] = sRGBToLinearTable[pixels[i]]
	}

	// Compute the factors
	for y := 0; y < height; y++ {
		cosYComponents := cosY[y]
		// Row start in linearPixels
		rowStart := y * bytesPerRow // bytesPerRow = width*3
		for x := 0; x < width; x++ {
			pIndex := rowStart + x*3
			rLinear := linearPixels[pIndex]
			gLinear := linearPixels[pIndex+1]
			bLinear := linearPixels[pIndex+2]

			cosXComponents := cosX[x]

			// Instead of doing `idx` increments, we compute offsets directly:
			// off = (j*xComponents + i)
			// just do two nested loops and calculate off once.
			// contribute contributions for each component
			for j := 0; j < yComponents; j++ {
				cy := cosYComponents[j]
				rowOffset := j * xComponents
				for i := 0; i < xComponents; i++ {
					c := cosXComponents[i] * cy
					off := rowOffset + i
					factorsR[off] += c * rLinear
					factorsG[off] += c * gLinear
					factorsB[off] += c * bLinear
				}
			}
		}
	}

	// Normalize the factors
	invCount := 1.0 / float32(width*height)
	for i := 0; i < factorsCount; i++ {
		factorsR[i] *= invCount
		factorsG[i] *= invCount
		factorsB[i] *= invCount
	}

	// Find max AC component
	var maxVal float32
	for i := 1; i < factorsCount; i++ {
		r := factorsR[i]
		if r < 0 {
			r = -r
		}
		if r > maxVal {
			maxVal = r
		}

		g := factorsG[i]
		if g < 0 {
			g = -g
		}
		if g > maxVal {
			maxVal = g
		}

		b := factorsB[i]
		if b < 0 {
			b = -b
		}
		if b > maxVal {
			maxVal = b
		}
	}

	quantMax := 0
	if maxVal > 0 {
		q := float64(maxVal)*166.0 - 1.0
		if q < 0 {
			q = 0
		} else if q > 82 {
			q = 82
		}
		quantMax = int(q)
	}

	maxAc := (float64(quantMax) + 1) / 166.0

	// Encode DC
	dcValue := encodeDC(float64(factorsR[0]), float64(factorsG[0]), float64(factorsB[0]))

	// Encode AC
	acValues := make([]int, factorsCount-1)
	for i := 1; i < factorsCount; i++ {
		acValues[i-1] = encodeAC(float64(factorsR[i])/maxAc, float64(factorsG[i])/maxAc, float64(factorsB[i])/maxAc)
	}

	// Build Blurhash string
	sizeFlag := (yComponents-1)*9 + (xComponents - 1)
	var result []byte
	result = append(result, encode83(sizeFlag, 1)...)
	result = append(result, encode83(quantMax, 1)...)
	result = append(result, encode83(dcValue, 4)...)
	for _, ac := range acValues {
		result = append(result, encode83(ac, 2)...)
	}

	return string(result)
}

// precomputeCosines generates a cosine table for given size and number of components using float32.
func precomputeCosines(size, components int) [][]float32 {
	cosines := make([][]float32, size)
	for i := 0; i < size; i++ {
		cosines[i] = make([]float32, components)
		for j := 0; j < components; j++ {
			angle := (math.Pi * float64(i) * float64(j)) / float64(size)
			cosines[i][j] = float32(math.Cos(angle))
		}
	}
	return cosines
}

// precomputeCosinesFloat64 generates a cosine table for decoding using float64.
func precomputeCosinesFloat64(size, components int) [][]float64 {
	cosines := make([][]float64, size)
	for i := 0; i < size; i++ {
		cosines[i] = make([]float64, components)
		for j := 0; j < components; j++ {
			cosines[i][j] = math.Cos(math.Pi * float64(i) * float64(j) / float64(size))
		}
	}
	return cosines
}

func linearToSRGB(value float64) int {
	if value <= 0 {
		return linearToSRGBTable[0]
	} else if value >= 1 {
		return linearToSRGBTable[4095]
	}
	index := int(value * 4095.0)
	return linearToSRGBTable[index]
}

func decodeInt(input string) (int, error) {
	result := 0
	for _, c := range input {
		if c > 127 {
			return 0, errors.New("invalid character in blurhash")
		}
		index := decodeCharacterMap[c]
		if index == -1 {
			return 0, errors.New("invalid character in blurhash")
		}
		result = result*83 + index
	}
	return result, nil
}

// IsValidBlurhash checks if the provided hash is a valid Blurhash.
func IsValidBlurhash(hash string) bool {
	if len(hash) < 6 {
		return false
	}

	sizeFlag, err := decodeInt(hash[:1])
	if err != nil {
		return false
	}

	numY := (sizeFlag / 9) + 1
	numX := (sizeFlag % 9) + 1
	expectedLength := 4 + 2*numX*numY
	return len(hash) == expectedLength
}

func IsValidExtendedBlurhash(hash string) bool {
	// Check if it follows the extended format: width!height!orientation!extension!blurhash
	parts := strings.Split(hash, "!")
	if len(parts) != 5 {
		return false // Invalid if not exactly 5 parts
	}

	// Validate width and height
	width, err := strconv.Atoi(parts[0])
	if err != nil || width <= 0 {
		return false
	}

	height, err := strconv.Atoi(parts[1])
	if err != nil || height <= 0 {
		return false
	}

	extension := parts[2]
	if len(extension) == 0 {
		return false // Extension must not be empty
	}

	// orientation := parts[3]
	// if len(orientation) == 0 {

	// }

	blurhash := parts[4]
	return IsValidBlurhash(blurhash)
}

// Helper functions
func maxFloat64(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func minFloat64(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func decodeDC(value int) (float64, float64, float64) {
	r := sRGBToLinearTable[(value>>16)&255]
	g := sRGBToLinearTable[(value>>8)&255]
	b := sRGBToLinearTable[value&255]
	return float64(r), float64(g), float64(b)
}

func decodeAC(value int, maximumValue float64) (float64, float64, float64) {
	quantR := value / (19 * 19)
	quantG := (value / 19) % 19
	quantB := value % 19

	r := signPow((float64(quantR)-9)/9.0, 2.0) * maximumValue
	g := signPow((float64(quantG)-9)/9.0, 2.0) * maximumValue
	b := signPow((float64(quantB)-9)/9.0, 2.0) * maximumValue
	return r, g, b
}

func encodeDC(r, g, b float64) int {
	return (linearToSRGB(r) << 16) + (linearToSRGB(g) << 8) + linearToSRGB(b)
}

func encodeAC(r, g, b float64) int {
	quantR := int(maxFloat64(0, minFloat64(18, math.Floor(signPow(r, 0.5)*9+9.5))))
	quantG := int(maxFloat64(0, minFloat64(18, math.Floor(signPow(g, 0.5)*9+9.5))))
	quantB := int(maxFloat64(0, minFloat64(18, math.Floor(signPow(b, 0.5)*9+9.5))))
	return quantR*19*19 + quantG*19 + quantB
}

func signPow(value, exp float64) float64 {
	if value < 0 {
		return -math.Pow(-value, exp)
	}
	return math.Pow(value, exp)
}

func encode83(value, length int) string {
	result := make([]byte, length)
	for i := length - 1; i >= 0; i-- {
		digit := value % 83
		result[i] = characters[digit]
		value /= 83
	}
	return string(result)
}

// Remaining functions for image manipulation and EXIF orientation

// Rotate90 rotates the image 90 degrees clockwise.
func Rotate90(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	rotated := image.NewRGBA(image.Rect(0, 0, bounds.Dy(), bounds.Dx()))
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			rotated.Set(bounds.Max.Y-y-1, x, img.At(x, y))
		}
	}
	return rotated
}

// Rotate180 rotates the image 180 degrees.
func Rotate180(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	rotated := image.NewRGBA(bounds)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			rotated.Set(bounds.Max.X-x-1, bounds.Max.Y-y-1, img.At(x, y))
		}
	}
	return rotated
}

// Rotate270 rotates the image 270 degrees clockwise (or 90 degrees counterclockwise).
func Rotate270(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	rotated := image.NewRGBA(image.Rect(0, 0, bounds.Dy(), bounds.Dx()))
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			rotated.Set(y, bounds.Max.X-x-1, img.At(x, y))
		}
	}
	return rotated
}

// FixOrientation adjusts the image based on the EXIF orientation.
func FixOrientation(img image.Image, orientation int) image.Image {
	switch orientation {
	case 3:
		return Rotate180(img)
	case 6:
		return Rotate90(img)
	case 8:
		return Rotate270(img)
	default:
		return img
	}
}

// GetOrientation retrieves the EXIF orientation of the image.
func GetOrientation(file *os.File) int {
	file.Seek(0, 0) // Reset file pointer
	x, err := exif.Decode(file)
	if err != nil {
		log.Println("No EXIF data found or error decoding EXIF:", err)
		return 1 // Default orientation
	}

	orientationTag, err := x.Get(exif.Orientation)
	if err != nil {
		log.Println("No orientation tag found in EXIF:", err)
		return 1 // Default orientation
	}

	orientation, err := orientationTag.Int(0)
	if err != nil {
		log.Println("Error reading orientation value:", err)
		return 1 // Default orientation
	}

	return orientation
}

// GetOrientationByReader retrieves the EXIF orientation of the image.
func GetOrientationByReader(file io.ReadSeeker) int {
	// Reset file pointer to the beginning
	file.Seek(0, 0)

	// Decode EXIF data
	x, err := exif.Decode(file)
	if err != nil {
		log.Println("No EXIF data found or error decoding EXIF:", err)
		return 1 // Default orientation
	}

	// Retrieve the orientation tag
	orientationTag, err := x.Get(exif.Orientation)
	if err != nil {
		log.Println("No orientation tag found in EXIF:", err)
		return 1 // Default orientation
	}

	// Convert orientation to integer
	orientation, err := orientationTag.Int(0)
	if err != nil {
		log.Println("Error reading orientation value:", err)
		return 1 // Default orientation
	}

	return orientation
}

// PixelsToImage converts raw RGB pixel data into an *image.RGBA.
func PixelsToImage(pixels []uint8, width, height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	pixelIndex := 0

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r := pixels[pixelIndex]
			g := pixels[pixelIndex+1]
			b := pixels[pixelIndex+2]
			pixelIndex += 3

			// Set pixel in the image
			img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}

	return img
}
