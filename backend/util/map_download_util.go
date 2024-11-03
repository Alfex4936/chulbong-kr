package util

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/go-pdf/fpdf"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
	"go.uber.org/fx"

	"golang.org/x/image/webp"
	// "github.com/chai2010/webp"
)

const (
	watermarkScale = 0.8
)

// TODO: DI?
// Global variables
var (
	watermarkImg    image.Image
	watermarkWidth  int
	watermarkHeight int

	markerIcon   image.Image
	markerWidth  int
	markerHeight int

	monthsInKorean = map[time.Month]string{
		time.January:   "1",
		time.February:  "2",
		time.March:     "3",
		time.April:     "4",
		time.May:       "5",
		time.June:      "6",
		time.July:      "7",
		time.August:    "8",
		time.September: "9",
		time.October:   "10",
		time.November:  "11",
		time.December:  "12",
	}

	HTTPClientUtil = &http.Client{
		Timeout: 10 * time.Second, // Set a timeout to avoid hanging requests indefinitely
	}
)

func RegisterPdfInitLifecycle(lifecycle fx.Lifecycle) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {

			// Load watermark image once
			watermarkImagePath := "resource/logo-text.png"
			watermarkImg, _, _ = loadImage(watermarkImagePath)
			watermarkBounds := watermarkImg.Bounds()
			watermarkWidth = watermarkBounds.Dx()
			watermarkHeight = watermarkBounds.Dy()

			// Load marker icon once
			markerIconPath := "resource/marker_40x40.webp"
			markerIcon, _ = LoadWebP(markerIconPath)
			markerBounds := markerIcon.Bounds()
			markerWidth = markerBounds.Dx()
			markerHeight = markerBounds.Dy()
			return nil
		},
		OnStop: func(context.Context) error {
			// Cleanup if necessary
			return nil
		},
	})

}

func OverlayImages(baseImageFile, markerImagePath string) (string, error) {
	originalBaseImg, _, err := loadImage(baseImageFile) // Load the original base image
	if err != nil {
		return "", err
	}
	originalBaseBounds := originalBaseImg.Bounds()

	resultImg := image.NewRGBA(originalBaseBounds)
	draw.Draw(resultImg, originalBaseBounds, originalBaseImg, image.Point{}, draw.Src)

	files, err := os.ReadDir(markerImagePath)
	if err != nil {
		return "", fmt.Errorf("failed to read marker image path: %w", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".png" {
			markerImg, _, err := loadImage(filepath.Join(markerImagePath, file.Name()))
			if err != nil {
				fmt.Println("Warning: skipping file due to error:", err)
				continue // Skip files that can't be loaded
			}
			overlayDifferences(resultImg, markerImg, originalBaseImg)
		}
	}

	resultPath := filepath.Join(markerImagePath, "result.png")
	if err := saveImage(resultImg, resultPath); err != nil {
		return "", fmt.Errorf("failed to save result image: %w", err)
	}

	return resultPath, nil
}

// PlaceMarkersOnImage places markers on the given base image according to their WCONGNAMUL coordinates.
func PlaceMarkersOnImage(baseImageFile string, markers []WCONGNAMULCoord, centerCX, centerCY float64) (string, error) {
	baseImg, _, err := loadImage(baseImageFile)
	if err != nil {
		return "", err
	}
	bounds := baseImg.Bounds()
	resultImg := image.NewRGBA(bounds)
	draw.Draw(resultImg, bounds, baseImg, image.Point{}, draw.Src)

	// SCALE by 2.5 in 1280x1080 image only, center (centerCX, centerCY).

	for _, marker := range markers {
		x, y := PlaceMarkerOnImage(marker.X, marker.Y, centerCX, centerCY, bounds.Dx(), bounds.Dy())

		// Calculate the top-left position to start drawing the marker icon
		// Ensure the entire marker icon is within bounds before drawing
		startX := x - markerWidth/2 - 5 // 5px out
		startY := y - markerHeight

		// Draw the resized marker icon
		draw.Draw(resultImg, image.Rect(startX, startY, startX+markerWidth, startY+markerHeight), markerIcon, image.Point{0, 0}, draw.Over)
	}

	// Add watermark image (logo-text.png)
	// Resize watermark image based on the given percentage
	newWatermarkWidth := uint(float64(watermarkWidth) * watermarkScale)
	newWatermarkHeight := uint(float64(watermarkHeight) * watermarkScale)
	resizedWatermarkImg := resize.Resize(newWatermarkWidth, newWatermarkHeight, watermarkImg, resize.Lanczos3)

	// Calculate the center position for the resized watermark image
	centerX := (bounds.Dx() - int(newWatermarkWidth)) / 2
	centerY := (bounds.Dy() - int(newWatermarkHeight)) / 2

	// Create an alpha mask for the watermark image with the desired transparency
	temp := 255 * 0.1
	alpha := uint8(temp) // 10% opacity
	alphaMask := image.NewUniform(color.Alpha{A: alpha})

	// Draw the watermark image with the alpha mask at the center position
	draw.DrawMask(resultImg, image.Rect(centerX, centerY, centerX+int(newWatermarkWidth), centerY+int(newWatermarkHeight)), resizedWatermarkImg, image.Point{}, alphaMask, image.Point{}, draw.Over)

	resultPath := filepath.Join(filepath.Dir(baseImageFile), "result_with_markers.png")
	if err := saveImage(resultImg, resultPath); err != nil {
		return "", fmt.Errorf("failed to save image with markers: %w", err)
	}

	return resultPath, nil
}

func loadImage(path string) (image.Image, string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()
	return image.Decode(file) // decoded image
}

func overlayDifferences(base *image.RGBA, overlay image.Image, originalBase image.Image) {
	bounds := base.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			baseColor := originalBase.At(x, y)
			overlayColor := overlay.At(x, y)
			if !colorsAreSimilar(baseColor, overlayColor) {
				base.Set(x, y, overlayColor)
			}
		}
	}
}

func colorsAreSimilar(c1, c2 color.Color) bool {
	const threshold = 10

	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()

	diffR := uint8(r1>>8) - uint8(r2>>8)
	diffG := uint8(g1>>8) - uint8(g2>>8)
	diffB := uint8(b1>>8) - uint8(b2>>8)

	dist := int(diffR)*int(diffR) + int(diffG)*int(diffG) + int(diffB)*int(diffB)

	return dist < threshold*threshold
}

func saveImage(img image.Image, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		return err
	}

	return nil
}

func GenerateMapPDF(imagePath, tempDir, title string, markerId int) (string, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.AddUTF8Font("NanumGothic", "", "resource/nanum.ttf")
	// Set the Korean font
	pdf.SetFont("NanumGothic", "", 10)
	totalWidth := 190.0 // the full width of the page area used

	// Print the date string, centered
	now := time.Now()
	dateString := fmt.Sprintf("%d년 %s월 %d일", now.Year(), monthsInKorean[now.Month()], now.Day())
	pdf.CellFormat(totalWidth, 10, dateString, "0", 1, "C", false, 0, "")

	// Print the title, centered
	pdf.SetFont("NanumGothic", "", 16) // Set font size to 16 for title
	pdf.SetTextColor(0, 0, 0)          // Reset text color to black for the title
	pdf.CellFormat(totalWidth, 10, title, "0", 1, "C", false, 0, "")

	url := fmt.Sprintf("https://k-pullup.com/pullup/%d", markerId)

	// Add image
	pdf.ImageOptions(imagePath, 10, 40, totalWidth, 0, false, fpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, url)

	// Print the "More at k-pullup.com" as one unit, centered
	pdf.SetFont("NanumGothic", "", 8) // Reset the font size to 10 if needed
	pdf.SetTextColor(0, 0, 0)         // Black color for "More at "
	moreAtText := "More at "
	fullLinkText := moreAtText + url
	// fullLinkWidth := pdf.GetStringWidth(fullLinkText)

	// Calculate starting X position to center the full link text
	// linkStartX := (totalWidth - fullLinkWidth) / 2
	// log.Printf("Starting X position to center the full link text: %f", linkStartX)
	pdf.SetX(87)
	pdf.SetY(200)
	pdf.CellFormat(totalWidth, 10, fullLinkText, "0", 0, "R", false, 0, "")
	// Create an invisible link over "k-pullup.com"
	pdf.LinkString(pdf.GetX()+pdf.GetStringWidth(moreAtText), pdf.GetY(), pdf.GetStringWidth(url), 10, url)

	// Save the PDF
	pdfName := fmt.Sprintf("kpullup-%s.pdf", uuid.New().String())
	pdfPath := path.Join(tempDir, pdfName)
	err := pdf.OutputFileAndClose(pdfPath)
	if err != nil {
		return "", err
	}

	// Remove the temporary image file
	os.Remove(imagePath)

	return pdfPath, nil
}

// Helper function to download a file from URL to a specific destination
func DownloadFile(URL, destPath string) error {
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	resp, err := HTTPClientUtil.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Create the file
	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// PlaceMarkerOnImage calculates (x,y) position in imagem SCALE by 2.5 in 1280x1080 image only.
func PlaceMarkerOnImage(CX, CY, centerCX, centerCY float64, imageWidth, imageHeight int) (int, int) {
	deltaX := CX - centerCX
	deltaY := CY - centerCY

	cxUnitsPerPixel := 3190.0 / float64(imageWidth)
	cyUnitsPerPixel := 3190.0 / float64(imageWidth)

	pixelOffsetX := deltaX / cxUnitsPerPixel
	pixelOffsetY := deltaY / cyUnitsPerPixel

	markerPosX := (imageWidth / 2) + int(pixelOffsetX)
	markerPosY := (imageHeight / 2) - int(pixelOffsetY)

	return markerPosX, markerPosY
}

// LoadWebP loads a WEBP image from the specified file path.
func LoadWebP(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return webp.Decode(file)
}

// PlaceMarkersOnImage places markers on the given base image according to their WCONGNAMUL coordinates.
func PlaceMarkersOnImageDynamic(baseImageFile string, markers []WCONGNAMULCoord, centerCX, centerCY, zoomScale float64) (string, error) {
	baseImg, _, err := loadImage(baseImageFile)
	if err != nil {
		return "", err
	}
	bounds := baseImg.Bounds()
	resultImg := image.NewRGBA(bounds)
	draw.Draw(resultImg, bounds, baseImg, image.Point{}, draw.Src)

	// SCALE by 2.5 in 1280x1080 image only, center (centerCX, centerCY).

	for _, marker := range markers {
		x, y := PlaceMarkerOnImageDynamic(marker.X, marker.Y, centerCX, centerCY, bounds.Dx(), bounds.Dy(), zoomScale)

		// Calculate the top-left position to start drawing the marker icon
		// Ensure the entire marker icon is within bounds before drawing
		startX := x - markerWidth/2 - 5 // 5px out
		startY := y - markerHeight

		// Draw the resized marker icon
		draw.Draw(resultImg, image.Rect(startX, startY, startX+markerWidth, startY+markerHeight), markerIcon, image.Point{0, 0}, draw.Over)
	}

	// Add watermark image (logo-text.png)
	// Resize watermark image based on the given percentage
	newWatermarkWidth := uint(float64(watermarkWidth) * watermarkScale)
	newWatermarkHeight := uint(float64(watermarkHeight) * watermarkScale)
	resizedWatermarkImg := resize.Resize(newWatermarkWidth, newWatermarkHeight, watermarkImg, resize.Lanczos3)

	// Calculate the center position for the resized watermark image
	centerX := (bounds.Dx() - int(newWatermarkWidth)) / 2
	centerY := (bounds.Dy() - int(newWatermarkHeight)) / 2

	// Create an alpha mask for the watermark image with the desired transparency
	temp := 255 * 0.1
	alpha := uint8(temp) // 10% opacity
	alphaMask := image.NewUniform(color.Alpha{A: alpha})

	// Draw the watermark image with the alpha mask at the center position
	draw.DrawMask(resultImg, image.Rect(centerX, centerY, centerX+int(newWatermarkWidth), centerY+int(newWatermarkHeight)), resizedWatermarkImg, image.Point{}, alphaMask, image.Point{}, draw.Over)

	resultPath := filepath.Join(filepath.Dir(baseImageFile), "result_with_markers.png")
	if err := saveImage(resultImg, resultPath); err != nil {
		return "", fmt.Errorf("failed to save image with markers: %w", err)
	}

	return resultPath, nil
}

// PlaceMarkerOnImageDynamic calculates (x, y) position in the image with a given zoom scale and dimensions.
func PlaceMarkerOnImageDynamic(CX, CY, centerCX, centerCY float64, imageWidth, imageHeight int, zoomScale float64) (int, int) {
	deltaX := CX - centerCX
	deltaY := CY - centerCY

	// Adjust units per pixel based on the zoom scale
	cxUnitsPerPixel := (zoomScale * float64(imageWidth))
	cyUnitsPerPixel := (zoomScale * float64(imageHeight))

	pixelOffsetX := deltaX / cxUnitsPerPixel
	pixelOffsetY := deltaY / cyUnitsPerPixel

	markerPosX := (imageWidth / 2) + int(pixelOffsetX)
	markerPosY := (imageHeight / 2) - int(pixelOffsetY)

	return markerPosX, markerPosY
}
