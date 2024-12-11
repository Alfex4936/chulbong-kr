// https://github.com/woltapp/blurhash

package util

import (
	"errors"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"strings"

	"github.com/rwcarlsen/goexif/exif"
)

const characters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz#$%*+,-.:;=?@[]^_{|}~"

var (
	sRGBToLinearTable  [256]float64
	linearToSRGBTable  [4096]int
	decodeCharacterMap [128]int
)

func init() {
	// Initialize sRGB to Linear lookup table
	for i := 0; i < 256; i++ {
		v := float64(i) / 255
		if v <= 0.04045 {
			sRGBToLinearTable[i] = v / 12.92
		} else {
			sRGBToLinearTable[i] = math.Pow((v+0.055)/1.055, 2.4)
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

// EncodeImage encodes an image.Image into a Blurhash string.
func EncodeImage(img image.Image, xComponents, yComponents int) string {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	// Preallocate pixel slice
	pixels := make([]uint8, width*height*3)
	idx := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			pixels[idx] = uint8(r >> 8)
			pixels[idx+1] = uint8(g >> 8)
			pixels[idx+2] = uint8(b >> 8)
			idx += 3
		}
	}

	// Encode the pixel data
	bytesPerRow := width * 3
	return Encode(xComponents, yComponents, width, height, pixels, bytesPerRow)
}

// Decode generates an image from a Blurhash string.
func Decode(hash string, width, height, punch int) ([]uint8, error) {
	if !IsValidBlurhash(hash) {
		return nil, errors.New("invalid blurhash")
	}

	sizeFlag, _ := decodeInt(hash[:1])
	numY := (sizeFlag / 9) + 1
	numX := (sizeFlag % 9) + 1

	if len(hash) != 4+2*numX*numY {
		return nil, errors.New("invalid blurhash length")
	}

	quantizedMaxValue, _ := decodeInt(hash[1:2])
	maxValue := (float64(quantizedMaxValue) + 1) / 166

	colors := make([][3]float64, numX*numY)
	for i := 0; i < numX*numY; i++ {
		if i == 0 {
			value, _ := decodeInt(hash[2:6])
			colors[0][0], colors[0][1], colors[0][2] = decodeDC(value)
		} else {
			start := 4 + i*2
			value, _ := decodeInt(hash[start : start+2])
			colors[i][0], colors[i][1], colors[i][2] = decodeAC(value, maxValue*float64(punch))
		}
	}

	// Precompute cosines
	cosX := precomputeCosinesDecode(width, numX)
	cosY := precomputeCosinesDecode(height, numY)

	// Decode the image
	pixels := make([]uint8, width*height*3)
	for y := 0; y < height; y++ {
		cosYComponent := cosY[y]
		for x := 0; x < width; x++ {
			r, g, b := 0.0, 0.0, 0.0
			cosXComponent := cosX[x]
			idx := 0
			for j := 0; j < numY; j++ {
				for i := 0; i < numX; i++ {
					basis := cosXComponent[i] * cosYComponent[j]
					color := colors[idx]
					r += color[0] * basis
					g += color[1] * basis
					b += color[2] * basis
					idx++
				}
			}
			pixelIndex := (y*width + x) * 3
			pixels[pixelIndex+0] = uint8(linearToSRGB(r))
			pixels[pixelIndex+1] = uint8(linearToSRGB(g))
			pixels[pixelIndex+2] = uint8(linearToSRGB(b))
		}
	}

	return pixels, nil
}

// Encode generates a Blurhash string for the given pixel data.
func Encode(xComponents, yComponents, width, height int, pixels []uint8, bytesPerRow int) string {
	if xComponents < 1 || xComponents > 9 || yComponents < 1 || yComponents > 9 {
		return ""
	}

	// Precompute cosine values
	cosX := precomputeCosines(width, xComponents)
	cosY := precomputeCosines(height, yComponents)

	// Initialize factors
	factorsCount := xComponents * yComponents
	factors := make([][3]float64, factorsCount)

	// Compute the factors
	for y := 0; y < height; y++ {
		cosYComponents := cosY[y]
		for x := 0; x < width; x++ {
			r := sRGBToLinear(int(pixels[y*bytesPerRow+x*3+0]))
			g := sRGBToLinear(int(pixels[y*bytesPerRow+x*3+1]))
			b := sRGBToLinear(int(pixels[y*bytesPerRow+x*3+2]))

			cosXComponents := cosX[x]
			idx := 0
			for j := 0; j < yComponents; j++ {
				cosYComponent := cosYComponents[j]
				for i := 0; i < xComponents; i++ {
					basis := cosXComponents[i] * cosYComponent
					factors[idx][0] += basis * r
					factors[idx][1] += basis * g
					factors[idx][2] += basis * b
					idx++
				}
			}
		}
	}

	// Normalize the factors
	for i := range factors {
		normalization := 1.0
		if i == 0 {
			normalization = 1.0
		} else {
			normalization = 2.0
		}
		scale := normalization / float64(width*height)
		factors[i][0] *= scale
		factors[i][1] *= scale
		factors[i][2] *= scale
	}

	// Encode DC component
	dc := factors[0]
	dcValue := encodeDC(dc[0], dc[1], dc[2])

	// Encode AC components
	maximumValue := 0.0
	for i := 1; i < factorsCount; i++ {
		maximumValue = math.Max(maximumValue, math.Abs(factors[i][0]))
		maximumValue = math.Max(maximumValue, math.Abs(factors[i][1]))
		maximumValue = math.Max(maximumValue, math.Abs(factors[i][2]))
	}

	quantizedMaxValue := int(math.Floor(maximumValue*166 - 0.5))
	if quantizedMaxValue < 0 {
		quantizedMaxValue = 0
	} else if quantizedMaxValue > 82 {
		quantizedMaxValue = 82
	}
	maximumValue = (float64(quantizedMaxValue) + 1) / 166

	acValues := make([]int, factorsCount-1)
	for i := 1; i < factorsCount; i++ {
		acValues[i-1] = encodeAC(factors[i][0], factors[i][1], factors[i][2], maximumValue)
	}

	// Build the Blurhash string
	var builder strings.Builder
	builder.Grow(2 + 4 + 2*(factorsCount-1)) // Preallocate the required size
	builder.WriteString(encode83((xComponents-1)+(yComponents-1)*9, 1))
	builder.WriteString(encode83(quantizedMaxValue, 1))
	builder.WriteString(encode83(dcValue, 4))
	for _, value := range acValues {
		builder.WriteString(encode83(value, 2))
	}

	return builder.String()
}

// Helper functions

func precomputeCosines(size, components int) [][]float64 {
	cosines := make([][]float64, size)
	for i := 0; i < size; i++ {
		cosines[i] = make([]float64, components)
		for j := 0; j < components; j++ {
			cosines[i][j] = math.Cos(math.Pi * float64(j) * float64(i) / float64(size))
		}
	}
	return cosines
}

func precomputeCosinesDecode(size, components int) [][]float64 {
	cosines := make([][]float64, size)
	for i := 0; i < size; i++ {
		cosines[i] = make([]float64, components)
		for j := 0; j < components; j++ {
			cosines[i][j] = math.Cos(math.Pi * float64(i) * float64(j) / float64(size))
		}
	}
	return cosines
}

func sRGBToLinear(value int) float64 {
	return sRGBToLinearTable[value]
}

func linearToSRGB(value float64) int {
	v := math.Max(0, math.Min(1, value))
	index := int(v * 4095)
	if index < 0 {
		index = 0
	} else if index > 4095 {
		index = 4095
	}
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

func decodeDC(value int) (float64, float64, float64) {
	r := sRGBToLinear(value >> 16)
	g := sRGBToLinear((value >> 8) & 255)
	b := sRGBToLinear(value & 255)
	return r, g, b
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

func encodeAC(r, g, b, maximumValue float64) int {
	quantR := int(math.Max(0, math.Min(18, math.Floor(signPow(r/maximumValue, 0.5)*9+9.5))))
	quantG := int(math.Max(0, math.Min(18, math.Floor(signPow(g/maximumValue, 0.5)*9+9.5))))
	quantB := int(math.Max(0, math.Min(18, math.Floor(signPow(b/maximumValue, 0.5)*9+9.5))))
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

// PixelsToImage converts raw RGB pixel data into an *image.RGBA
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

	// img := image.NewRGBA(image.Rect(0, 0, width, height))
	// copy(img.Pix, pixels)

	return img
}
