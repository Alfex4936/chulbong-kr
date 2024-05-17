package csw.chulbongkr.service.local;

import csw.chulbongkr.util.CoordinatesConverter;
import jakarta.annotation.PostConstruct;
import lombok.Getter;
import lombok.RequiredArgsConstructor;
import org.apache.pdfbox.pdmodel.PDDocument;
import org.apache.pdfbox.pdmodel.PDPage;
import org.apache.pdfbox.pdmodel.PDPageContentStream;
import org.apache.pdfbox.pdmodel.font.PDType0Font;
import org.apache.pdfbox.pdmodel.graphics.image.PDImageXObject;
import org.springframework.stereotype.Service;

import javax.imageio.ImageIO;
import javax.imageio.ImageReader;
import javax.imageio.stream.ImageInputStream;
import java.awt.*;
import java.awt.geom.AffineTransform;
import java.awt.image.BufferedImage;
import java.io.File;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.nio.file.Path;
import java.time.LocalDate;
import java.time.format.DateTimeFormatter;
import java.util.Iterator;
import java.util.List;
import java.util.Random;
import java.util.UUID;

@RequiredArgsConstructor
@Service
public class ImageProcessorService {

    private final FileCleanupService fileCleanupService;

    @Getter
    private BufferedImage markerIcon;
    @Getter
    private File nanumFont;

    @PostConstruct
    public void init() throws IOException {
        markerIcon = loadWebP("map_marker.webp");
        nanumFont = loadNanumFont("fonts/nanum.ttf");
    }

    public String placeMarkersOnImage(String baseImageFile, List<CoordinatesConverter.XYCoordinate> markers, double centerCX, double centerCY) throws IOException {
        BufferedImage baseImage = ImageIO.read(new File(baseImageFile));
        Graphics2D graphics = baseImage.createGraphics();

        for (var marker : markers) {
            int x = calculateX(marker.latitude(), centerCX, baseImage.getWidth());
            int y = calculateY(marker.longitude(), centerCY, baseImage.getHeight());
            drawMarker(graphics, markerIcon, x, y);
        }

        // Add watermark text
        addWatermarkText(graphics, baseImage);

        String outputFilePath = "markers_" + UUID.randomUUID() + ".png";
        Path destPath = fileCleanupService.getTempDir().resolve(outputFilePath);

        ImageIO.write(baseImage, "png", new File(destPath.toString()));
        return destPath.toString();
    }

    protected void addWatermarkText(Graphics2D graphics, BufferedImage baseImage) {
        String watermarkText = "k-pullup.com";
        Font font = new Font("Arial", Font.BOLD, 40);
        AlphaComposite alphaChannel = AlphaComposite.getInstance(AlphaComposite.SRC_OVER, 0.1f); // transparency
        graphics.setComposite(alphaChannel);
        graphics.setColor(Color.GRAY);
        graphics.setFont(font);

        FontMetrics fontMetrics = graphics.getFontMetrics();
        int centerX = (baseImage.getWidth() - fontMetrics.stringWidth(watermarkText)) / 2;
        int centerY = baseImage.getHeight() / 2;

        // Create a random rotation angle between -10 and 10 degrees
        Random random = new Random();
        double rotationAngle = Math.toRadians(random.nextInt(21) - 10); // random angle between -10 and 10 degrees

        // Save the original transform
        AffineTransform originalTransform = graphics.getTransform();

        // Apply rotation around the center of the text
        graphics.rotate(rotationAngle, centerX + fontMetrics.stringWidth(watermarkText) / 2.0, centerY);

        // Draw the watermark text
        graphics.drawString(watermarkText, centerX, centerY);

        // Restore the original transform
        graphics.setTransform(originalTransform);
    }

    private int calculateX(double cx, double centerCX, int imageWidth) {
        double deltaX = cx - centerCX;
        double unitsPerPixel = 3190.0 / imageWidth;
        return (int) ((imageWidth / 2) + (deltaX / unitsPerPixel));
    }

    private int calculateY(double cy, double centerCY, int imageHeight) {
        double deltaY = cy - centerCY;
        double unitsPerPixel = 3190.0 / imageHeight;
        return (int) ((imageHeight / 2) - (deltaY / unitsPerPixel));
    }

    private void drawMarker(Graphics2D graphics, BufferedImage markerIcon, int x, int y) {
        int markerWidth = markerIcon.getWidth();
        int markerHeight = markerIcon.getHeight();
        int startX = x - markerWidth / 2 - 5; // 5px out
        int startY = y - markerHeight;
        graphics.drawImage(markerIcon, startX, startY, null);
    }

    public String generateMapPDF(String imagePath, String title) throws IOException {
        try (PDDocument document = new PDDocument()) {
            PDPage page = new PDPage();
            document.addPage(page);

            PDType0Font font = PDType0Font.load(document, nanumFont);
            PDImageXObject pdImage = PDImageXObject.createFromFile(imagePath, document);

            float pageWidth = page.getMediaBox().getWidth();
            float titleFontSize = 16;
            float dateFontSize = 12;

            try (PDPageContentStream contentStream = new PDPageContentStream(document, page)) {
                // Center the date at the top
                String date = getFormattedDate();
                float dateWidth = font.getStringWidth(date) / 1000 * dateFontSize;
                float dateXOffset = (pageWidth - dateWidth) / 2;
                contentStream.beginText();
                contentStream.setFont(font, dateFontSize);
                contentStream.newLineAtOffset(dateXOffset, 750);
                contentStream.showText(date);
                contentStream.endText();

                // Center the title below the date
                float titleWidth = font.getStringWidth(title) / 1000 * titleFontSize;
                float titleXOffset = (pageWidth - titleWidth) / 2;
                contentStream.beginText();
                contentStream.setFont(font, titleFontSize);
                contentStream.newLineAtOffset(titleXOffset, 720); // decreasing = less padding
                contentStream.showText(title);
                contentStream.endText();

                // Draw the image below the title
                contentStream.drawImage(pdImage, 50, 300, 500, 400);

                // Add additional text below the image
                contentStream.beginText();
                contentStream.setFont(font, 10);
                contentStream.newLineAtOffset(50, 280);
                contentStream.showText("More at https://k-pullup.com");
                contentStream.endText();
            }

            String pdfPath = "kpullup-" + UUID.randomUUID() + ".pdf";
            Path destPath = fileCleanupService.getTempDir().resolve(pdfPath);
            document.save(destPath.toFile());
            return destPath.toString();
        }
    }

    private String getFormattedDate() {
        LocalDate now = LocalDate.now();
        DateTimeFormatter formatter = DateTimeFormatter.ofPattern("yyyy년 MM월 dd일");
        return now.format(formatter);
    }

    protected BufferedImage loadWebP(String resourcePath) throws IOException {
        try (InputStream input = getClass().getClassLoader().getResourceAsStream(resourcePath)) {
            if (input == null) {
                throw new IOException("Resource not found: " + resourcePath);
            }

            try (ImageInputStream imageInput = ImageIO.createImageInputStream(input)) {
                Iterator<ImageReader> readers = ImageIO.getImageReadersByFormatName("WEBP");
                if (!readers.hasNext()) {
                    throw new IOException("No WEBP readers found");
                }
                ImageReader reader = readers.next();
                reader.setInput(imageInput);
                return reader.read(0);
            }
        }
    }

    protected File loadNanumFont(String resourcePath) throws IOException {
        try (InputStream input = getClass().getClassLoader().getResourceAsStream(resourcePath)) {
            if (input == null) {
                throw new IOException("Resource not found: " + resourcePath);
            }

            // Create a temporary file
            File tempFile = File.createTempFile("nanum", ".ttf");
            tempFile.deleteOnExit();

            try (FileOutputStream out = new FileOutputStream(tempFile)) {
                byte[] buffer = new byte[1024];
                int bytesRead;
                while ((bytesRead = input.read(buffer)) != -1) {
                    out.write(buffer, 0, bytesRead);
                }
            }

            return tempFile;
        }
    }
}
