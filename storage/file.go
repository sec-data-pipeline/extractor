package storage

import (
	"bytes"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type FileStorage interface {
	PutObject(key string, data []byte) error
}

type s3Bucket struct {
	name   string
	client *s3.S3
}

func NewS3Bucket(awsSession *session.Session, name string) *s3Bucket {
	return &s3Bucket{name: name, client: s3.New(awsSession)}
}

func (b *s3Bucket) PutObject(key string, data []byte) error {
	input := &s3.PutObjectInput{
		Bucket: aws.String(b.name),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	}
	_, err := b.client.PutObject(input)
	if err != nil {
		return err
	}
	return nil
}

type folder struct {
	path string
}

func NewFolder(path string) *folder {
	return &folder{path: path}
}

func (f *folder) PutObject(key string, data []byte) error {
	err := buildPath(f.path + "/" + key)
	if err != nil {
		return err
	}
	err = os.WriteFile(f.path+"/"+key, data, 0666)
	if err != nil {
		return err
	}
	return nil
}

func buildPath(path string) error {
	var i int = len(path) - 1
	for ; i <= 0; i-- {
		if string(path[i]) == "/" {
			break
		}
	}
	err := os.MkdirAll(path[i:], os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
