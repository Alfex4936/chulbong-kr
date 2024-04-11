package utils

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/go-pdf/fpdf"
	"github.com/google/uuid"
)

var HTTPClientUtil = &http.Client{
	Timeout: 10 * time.Second, // Set a timeout to avoid hanging requests indefinitely
}

func OverlayImages(baseImageFile, markerImagePath string) (string, error) {
	originalBaseImg, _, err := loadImage(baseImageFile) // Load the original base image
	if err != nil {
		panic(err) // Handle the error properly
	}
	originalBaseBounds := originalBaseImg.Bounds()

	resultImg := image.NewRGBA(originalBaseBounds)
	draw.Draw(resultImg, originalBaseBounds, originalBaseImg, image.Point{}, draw.Src)

	files, err := ioutil.ReadDir(markerImagePath)
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

func GenerateMapPDF(imagePath, tempDir, title string) (string, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.AddUTF8Font("NanumGothic", "", "fonts/nanum.ttf")
	// Korean font

	pdf.SetFont("NanumGothic", "", 10)
	pdf.CellFormat(190, 10, "More at k-pullup.com", "0", 1, "C", false, pdf.AddLink(), "https://k-pullup.com")
	pdf.SetFont("NanumGothic", "", 16)
	pdf.CellFormat(190, 10, title, "0", 1, "C", false, 0, "")

	pdf.ImageOptions(imagePath, 10, 30, 190, 0, false, fpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")

	pdfName := fmt.Sprintf("kpullup-%s.pdf", uuid.New().String())
	pdfPath := path.Join(tempDir, pdfName)
	err := pdf.OutputFileAndClose(pdfPath)
	if err != nil {
		return "", err
	}

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
