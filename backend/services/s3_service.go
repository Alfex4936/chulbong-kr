package services

import (
	"chulbong-kr/database"
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

func UploadFileToS3(folder string, file *multipart.FileHeader) (string, error) {
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
	key := fmt.Sprintf("%s/%s%s", folder, uuid.String(), fileExtension)

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
		bucketName = S3_BUCKET_NAME
		key = dataURL
	}

	if key == "" {
		return fmt.Errorf("invalid key")
	}

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

	// Wait until the object is deleted
	return nil
}

func FetchAllPhotoURLsFromDB() ([]string, error) {
	query := `
        SELECT PhotoURL FROM Photos WHERE PhotoURL IS NOT NULL
        UNION
        SELECT URL FROM MarkerAddressFailures WHERE URL IS NOT NULL
        UNION
        SELECT ReportImageURL FROM Reports WHERE ReportImageURL IS NOT NULL
    `
	// Execute the query
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying database: %w", err)
	}
	defer rows.Close()

	var urls []string
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		urls = append(urls, url)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return urls, nil
}

func ListAllObjectsInS3() ([]string, error) {
	// Load the AWS credentials
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(AWS_REGION),
	)
	if err != nil {
		return nil, fmt.Errorf("error loading AWS config: %w", err)
	}

	s3Client := s3.NewFromConfig(cfg)
	input := &s3.ListObjectsV2Input{
		Bucket: &S3_BUCKET_NAME,
	}

	var s3Keys []string
	paginator := s3.NewListObjectsV2Paginator(s3Client, input)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, fmt.Errorf("error listing S3 objects: %w", err)
		}
		for _, item := range output.Contents {
			s3Keys = append(s3Keys, *item.Key)
		}
	}

	return s3Keys, nil
}

// FindUnreferencedS3Objects finds S3 objects that are not referenced in the database.
func FindUnreferencedS3Objects(dbURLs []string, s3Keys []string) []string {
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
