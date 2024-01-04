package main

import (
	"errors"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/sec-data-pipeline/extractor/api"
	"github.com/sec-data-pipeline/extractor/storage"
)

var db *storage.Database
var archive *storage.Bucket

func main() {
	for {
		companies, err := db.GetCompanies()
		if err != nil {
			panic(err)
		}
		for _, cmp := range companies {
			filingIDs, err := db.GetFilingIDs(cmp)
			if err != nil {
				panic(err)
			}
			filings, err := api.GetNewFilings(cmp.CIK, filingIDs)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			for _, flng := range filings {
				file, err := api.GetMainFile(cmp.CIK, flng)
				if err != nil {
					log.Println(err.Error())
					continue
				}
				err = db.InsertFiling(
					cmp.ID,
					flng.SECID,
					flng.Form,
					file.Name,
					flng.FilingDate,
					flng.ReportDate,
					flng.AcceptanceDate,
					file.LastModified,
				)
				if err != nil {
					log.Println(err.Error())
					continue
				}
				err = archive.PutObject(flng.SECID+file.Extension, file.Content)
				if err != nil {
					log.Println(err.Error())
				}
			}
		}
	}
}

func init() {
	region := os.Getenv("REGION")
	if len(region) < 1 {
		panic(errors.New("Environment variable 'REGION' must be specified"))
	}
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		panic(err)
	}
	secretsARN := os.Getenv("SECRETS")
	if len(secretsARN) < 1 {
		panic(errors.New("Environment variable 'SECRETS' must be specified"))
	}
	connParams, err := storage.GetAWSConnParams(awsSession, secretsARN)
	if err != nil {
		panic(err)
	}
	db, err = storage.NewDB(connParams)
	if err != nil {
		panic(err)
	}
	archiveName := os.Getenv("ARCHIVE_BUCKET")
	if len(archiveName) < 1 {
		panic(errors.New("Environment variable 'ARCHIVE_BUCKET' must be specified"))
	}
	archive = storage.NewBucket(awsSession, archiveName)
}
