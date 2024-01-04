package storage

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Bucket struct {
	name   string
	client *s3.S3
}

func NewBucket(awsSession *session.Session, name string) *Bucket {
	client := s3.New(awsSession)
	return &Bucket{name: name, client: client}
}

func (bucket *Bucket) PutObject(name string, content []byte) error {
	input := &s3.PutObjectInput{
		Bucket: aws.String(bucket.name),
		Key:    aws.String(name),
		Body:   bytes.NewReader(content),
	}
	_, err := bucket.client.PutObject(input)
	if err != nil {
		return err
	}
	return nil
}
