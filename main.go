package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/sec-data-pipeline/extractor/external"
	"github.com/sec-data-pipeline/extractor/service"
	"github.com/sec-data-pipeline/extractor/storage"
)

var extractor *service.Extractor

func main() {
	err := extractor.Run()
	if err != nil {
		panic(err)
	}
}

func init() {
	var err error
	region := envOrPanic("REGION")
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		panic(err)
	}
	var secrets storage.Secrets = storage.NewSecretsManager(awsSession, envOrPanic("SECRETS_ARN"))
	var archive storage.FileStorage = storage.NewS3Bucket(awsSession, envOrPanic("ARCHIVE_BUCKET"))
	var logger storage.Logger = storage.NewCloudWatch()
	connParams, err := secrets.GetConnParams()
	if err != nil {
		panic(err)
	}
	db, err := storage.NewPostgresConn(connParams)
	if err != nil {
		panic(err)
	}
	api := external.NewAPI()
	extractor = service.NewExtractorService(api, db, archive, logger)
}

func envOrPanic(key string) string {
	value := os.Getenv(key)
	if len(value) < 1 {
		panic(errors.New(fmt.Sprintf("Environment variable '%s' must be specified", key)))
	}
	return value
}
