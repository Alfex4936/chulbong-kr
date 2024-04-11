package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	"github.com/go-pdf/fpdf"
)

func main() {
	originalBaseImg, _, err := loadImage("tests/base.png") // Load the original base image
	if err != nil {
		panic(err) // Handle the error properly
	}
	originalBaseBounds := originalBaseImg.Bounds()

	resultImg := image.NewRGBA(originalBaseBounds)
	draw.Draw(resultImg, originalBaseBounds, originalBaseImg, image.Point{}, draw.Src)

	markerImages := []string{"tests/2.png", "tests/1.png"}

	for _, markerPath := range markerImages {
		markerImg, _, err := loadImage(markerPath)
		if err != nil {
			panic(err) // Handle the error properly
		}
		overlayDifferences(resultImg, markerImg, originalBaseImg)
	}

	// Save the resulting image as PNG
	saveImage(resultImg, "tests/result.png")

	// Generate PDF with the image
	generatePDF("tests/result.png", "tests/output.pdf", "제주특별자치도 제주시 영평동 2200") // Korean for "Korean Title"
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

func saveImage(img image.Image, path string) {
	file, err := os.Create(path)
	if err != nil {
		panic(err) // Handle the error properly
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		panic(err) // Handle the error properly
	}
}

func generatePDF(imagePath, pdfPath, title string) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.AddUTF8Font("NanumGothic", "", "fonts/nanum.ttf")
	// Korean font

	pdf.SetFont("NanumGothic", "", 10)
	pdf.CellFormat(190, 10, "More at k-pullup.com", "0", 1, "C", false, pdf.AddLink(), "https://k-pullup.com")
	pdf.SetFont("NanumGothic", "", 16)
	pdf.CellFormat(190, 10, title, "0", 1, "C", false, 0, "")

	pdf.ImageOptions(imagePath, 10, 30, 190, 0, false, fpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")

	err := pdf.OutputFileAndClose(pdfPath)
	if err != nil {
		fmt.Println(err)
	}
}
