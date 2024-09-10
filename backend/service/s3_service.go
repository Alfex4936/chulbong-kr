package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	myconfig "github.com/Alfex4936/chulbong-kr/config"
	"go.uber.org/zap"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type S3Service struct {
	Config *myconfig.S3Config
	Redis  *RedisService

	logger *zap.Logger
}

func NewS3Service(config *myconfig.S3Config, redis *RedisService, logger *zap.Logger) *S3Service {
	return &S3Service{
		Config: config,
		Redis:  redis,
		logger: logger,
	}
}

func (s *S3Service) UploadFileToS3(folder string, file *multipart.FileHeader) (string, error) {
	// Open the uploaded file
	fileData, err := file.Open()
	if err != nil {
		return "", err
	}
	defer fileData.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Load the AWS credentials
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(s.Config.AwsRegion))
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
	fileExtension := strings.ToLower(filepath.Ext(file.Filename))
	randomName := uuid.String()
	key := fmt.Sprintf("%s/%s%s", folder, randomName, fileExtension)

	// Upload the file to S3
	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &s.Config.S3BucketName,
		Key:    &key,
		Body:   fileData,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	// Construct the file URL
	fileURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.Config.S3BucketName, key)

	return fileURL, nil
}

// DeleteDataFromS3 deletes a photo from S3 given its URL.
func (s *S3Service) DeleteDataFromS3(dataURL string) error {
	var bucketName, key string

	// Attempt to parse the input as a URL
	parsedURL, err := url.Parse(dataURL)
	if err == nil && parsedURL.Scheme != "" && parsedURL.Host != "" {
		// It's a valid URL
		parts := strings.SplitN(parsedURL.Host, ".", 2)
		if len(parts) < 2 {
			return fmt.Errorf("invalid S3 URL format")
		}
		bucketName = parts[0]
		key = strings.TrimPrefix(parsedURL.Path, "/")
	} else {
		// It's not a valid URL, treat it as a key
		bucketName = s.Config.S3BucketName
		key = dataURL
	}

	if key == "" {
		return fmt.Errorf("invalid key")
	}

	// if isImage(filepath.Ext(key)) {
	// 	s.Redis.ResetCache("image:" + key)
	// }

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Load the AWS credentials
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(s.Config.AwsRegion),
	)
	if err != nil {
		return fmt.Errorf("could not load AWS credentials: %w", err)
	}

	// Create an S3 client
	s3Client := s3.NewFromConfig(cfg)

	// Delete the object
	_, err = s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucketName,
		Key:    &key,
	})
	if err != nil {
		return fmt.Errorf("failed to delete object from S3: %w", err)
	}

	// Wait until the object is deleted
	return nil
}

func (s *S3Service) ListAllObjectsInS3() ([]map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Load the AWS credentials
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(s.Config.AwsRegion),
	)
	if err != nil {
		return nil, fmt.Errorf("error loading AWS config: %w", err)
	}

	s3Client := s3.NewFromConfig(cfg)
	input := &s3.ListObjectsV2Input{
		Bucket: &s.Config.S3BucketName,
	}

	var s3Objects []map[string]interface{}
	var sumKB int64
	paginator := s3.NewListObjectsV2Paginator(s3Client, input)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, fmt.Errorf("error listing S3 objects: %w", err)
		}
		for _, item := range output.Contents {
			sizeKB := *item.Size / 1024 // Size in KB

			sumKB += sizeKB

			s3Objects = append(s3Objects, map[string]interface{}{
				"Key":  *item.Key,
				"Size": sizeKB,
			})
		}
	}

	s.logger.Info("💖 Total image size",
		zap.Int("number_of_images", len(s3Objects)),
		zap.Int64("total_size_kb", sumKB),
		zap.Int64("total_size_mb", sumKB/1024),
	)
	return s3Objects, nil
}

// FindUnreferencedS3Objects finds S3 objects that are not referenced in the database.
func (s *S3Service) FindUnreferencedS3Objects(dbURLs []string, s3Keys []string) []string {
	dbURLMap := make(map[string]struct{})
	for _, dbURL := range dbURLs {
		parsedURL, err := url.Parse(dbURL)
		if err == nil && parsedURL.Scheme != "" && parsedURL.Host != "" {
			key := strings.TrimPrefix(parsedURL.Path, "/")
			dbURLMap[key] = struct{}{}
		} else {
			// maybe just a key
			dbURLMap[dbURL] = struct{}{}
		}
	}

	var unreferenced []string
	for _, key := range s3Keys {
		if _, found := dbURLMap[key]; !found {
			unreferenced = append(unreferenced, key)
		}
	}

	return unreferenced
}

func (s *S3Service) MoveFileInS3(sourceKey string, destinationKey string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Load the AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(s.Config.AwsRegion),
	)
	if err != nil {
		return fmt.Errorf("could not load AWS configuration: %w", err)
	}

	// Create an S3 client
	s3Client := s3.NewFromConfig(cfg)

	copySource := url.PathEscape(s.Config.S3BucketName + "/" + sourceKey)

	// Copy the object to the new location
	_, err = s3Client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(s.Config.S3BucketName),
		CopySource: aws.String(copySource),
		Key:        aws.String(destinationKey),
	})
	if err != nil {
		return fmt.Errorf("failed to copy file in S3: %w", err)
	}

	// Delete the original object
	_, err = s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.Config.S3BucketName),
		Key:    aws.String(sourceKey),
	})
	if err != nil {
		return fmt.Errorf("failed to delete original file in S3: %w", err)
	}

	return nil
}

// Helper function to determine if a file extension corresponds to an image
func isImage(ext string) bool {
	// Normalize the extension to lower case
	ext = strings.ToLower(ext)
	// List of supported image extensions
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif", ".webp":
		return true
	default:
		return false
	}
}
