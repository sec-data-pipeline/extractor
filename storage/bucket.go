package storage

import (
	"bytes"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sec-data-pipeline/extractor/filing"
)

type Bucket struct {
	name   string
	client *s3.S3
}

func NewBucket() (*Bucket, error) {
	region := os.Getenv("REGION")
	if len(region) < 1 {
		return nil, errors.New("Environment variable 'REGION' must be specified")
	}
	bucket := os.Getenv("ARCHIVE_BUCKET")
	if len(bucket) < 1 {
		return nil, errors.New("Environment variable 'ARCHIVE_BUCKET' must be specified")
	}
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, err
	}
	client := s3.New(awsSession)
	return &Bucket{name: bucket, client: client}, nil
}

func (bucket *Bucket) PutObject(fil *filing.Filing) error {
	input := &s3.PutObjectInput{
		Bucket: aws.String(bucket.name),
		Key:    aws.String(fil.SECID + filing.GetFileExtension(fil.File.Name)),
		Body:   bytes.NewReader(fil.File.Content),
	}
	_, err := bucket.client.PutObject(input)
	if err != nil {
		return err
	}
	return nil
}
