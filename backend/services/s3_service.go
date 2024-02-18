package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

var AWS_REGION string
var S3_BUCKET_NAME string

func UploadFileToS3(file *multipart.FileHeader) (string, error) {
	// Open the uploaded file
	fileData, err := file.Open()
	if err != nil {
		return "", err
	}
	defer fileData.Close()

	// Load the AWS credentials
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(AWS_REGION),
	)
	if err != nil {
		return "", fmt.Errorf("could not load AWS credentials: %w", err)
	}

	// Create an S3 client
	s3Client := s3.NewFromConfig(cfg)

	// Generate a UUID for a unique filename
	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("failed to generate UUID: %w", err)
	}

	// Use the original file's extension but with a new UUID as the filename
	fileExtension := filepath.Ext(file.Filename)
	key := fmt.Sprintf("%s%s", uuid.String(), fileExtension)

	// Upload the file to S3
	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &S3_BUCKET_NAME,
		Key:    &key,
		Body:   fileData,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	// Construct the file URL
	fileURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", S3_BUCKET_NAME, key)

	return fileURL, nil
}

// DeleteDataFromS3 deletes a photo from S3 given its URL.
func DeleteDataFromS3(dataURL string) error {
	parsedURL, err := url.Parse(dataURL)
	if err != nil {
		return err
	}

	// Extract bucket name and key from the URL
	bucketName := strings.Split(parsedURL.Host, ".")[0]
	key := strings.TrimPrefix(parsedURL.Path, "/")

	// Load the AWS credentials
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(AWS_REGION),
	)
	if err != nil {
		return fmt.Errorf("could not load AWS credentials: %w", err)
	}

	// Create an S3 client
	s3Client := s3.NewFromConfig(cfg)

	// Delete the object
	_, err = s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: &bucketName,
		Key:    &key,
	})
	if err != nil {
		return fmt.Errorf("failed to delete object from S3: %w", err)
	}

	return nil
}
