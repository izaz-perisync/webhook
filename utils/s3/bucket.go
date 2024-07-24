package bucket

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Service interface {
	Upload(ctx context.Context, fileData []byte, fileName string) (*string, error)
	Delete(ctx context.Context, objectKey string) error
}

type Service struct {
	client *s3.Client
	folder string
	bucket string
	region string
}

func New(baseEndpoint, accessKey, secretKey, region, bucket, folder string) (S3Service, error) {

	opts := s3.Options{
		Region:      *aws.String(region),
		Credentials: credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
	}

	if baseEndpoint != "" {
		opts.BaseEndpoint = aws.String(baseEndpoint)
	}

	client := s3.New(opts)
	if client == nil {
		return nil, errors.New("s3 client is nil")
	}

	return &Service{
		client: client,
		folder: folder,
		bucket: bucket,
		region: region,
	}, nil
}

func (s *Service) Upload(ctx context.Context, fileData []byte, fileName string) (*string, error) {
	switch {
	case fileName == "":
		return nil, errors.New("file name is empty")

	case len(fileData) < 1:
		return nil, errors.New("file is empty")
	}

	filePath := fileName // stores asset on root
	if s.folder != "" {  // stores asset in a specific path
		filePath = s.folder + fileName
	}

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filePath),
		Body:   bytes.NewReader(fileData),
	})
	if err != nil {
		return nil, err
	}

	imageURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, s.region, filePath)

	return &imageURL, nil
}

func (s *Service) Delete(ctx context.Context, objectKey string) error {

	if objectKey == "" {
		return errors.New("url is empty")
	}

	prefix := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/", s.bucket, s.region)

	key := strings.TrimPrefix(objectKey, prefix)

	deleteObjectInput := &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}

	_, err := s.client.DeleteObject(ctx, deleteObjectInput)
	return err
}
