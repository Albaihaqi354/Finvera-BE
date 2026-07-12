package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
)

type StorageService interface {
	UploadFile(ctx context.Context, file multipart.File, filename string, contentType string) (string, error)
}

type r2Storage struct {
	client     *s3.Client
	bucketName string
	publicURL  string
}

func NewR2Storage(accountID, accessKey, secretKey, bucketName, publicURL string) (StorageService, error) {
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID),
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	return &r2Storage{
		client:     client,
		bucketName: bucketName,
		publicURL:  publicURL,
	}, nil
}

func (s *r2Storage) UploadFile(ctx context.Context, file multipart.File, filename string, contentType string) (string, error) {
	uploader := manager.NewUploader(s.client)

	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(filename),
		Body:        file,
		ContentType: aws.String(contentType),
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	// Return the public URL for the uploaded file
	return fmt.Sprintf("%s/%s", s.publicURL, filename), nil
}
