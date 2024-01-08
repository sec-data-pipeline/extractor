package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/sec-data-pipeline/filing-extractor/external"
	"github.com/sec-data-pipeline/filing-extractor/service"
	"github.com/sec-data-pipeline/filing-extractor/storage"
)

var extractor *service.Extractor

func main() {
	err := extractor.Run()
	if err != nil {
		panic(err)
	}
}

func init() {
	region := os.Getenv("REGION")
	var secrets storage.Secrets
	var archive storage.FileStorage
	var logger storage.Logger
	var err error
	if len(region) < 1 {
		secrets, err = storage.NewEnvLoader()
		if err != nil {
			panic(err)
		}
		archive = storage.NewFolder(envOrPanic("ARCHIVE_PATH"))
		logger = storage.NewConsole()
	} else {
		awsSession, err := session.NewSession(&aws.Config{
			Region: aws.String(region),
		})
		if err != nil {
			panic(err)
		}
		secrets = storage.NewSecretsManager(awsSession, envOrPanic("SECRETS_ARN"))
		archive = storage.NewS3Bucket(awsSession, envOrPanic("ARCHIVE_BUCKET"))
		logger = storage.NewCloudWatch()
	}
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
