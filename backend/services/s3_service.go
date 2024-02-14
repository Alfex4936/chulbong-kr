package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

func UploadFileToS3(file *multipart.FileHeader) (string, error) {
	// Open the uploaded file
	fileData, err := file.Open()
	if err != nil {
		return "", err
	}
	defer fileData.Close()

	// Load the AWS credentials
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("AWS_REGION")),
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
	bucketName := os.Getenv("AWS_BUCKET_NAME")
	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &key,
		Body:   fileData,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	// Construct the file URL
	fileURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, key)

	return fileURL, nil
}
