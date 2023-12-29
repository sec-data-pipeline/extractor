package main

import (
	"log"
	"time"

	"github.com/sec-data-pipeline/extractor/request"
	"github.com/sec-data-pipeline/extractor/storage"
	"github.com/sec-data-pipeline/extractor/utils"
)

var db *storage.Database
var archive *storage.Bucket

func main() {
	filingCount, err := db.GetFilingCount()
	if err != nil {
		panic(err)
	}
	if filingCount > 0 {
		filing, err := db.GetLatestFiling()
		if err != nil {
			panic(err)
		}
		localFiles, err := db.GetFiles(filing)
		filesData, err := request.GetFilesData(filing)
		if err != nil {
			panic(err)
		}
		newFiles := utils.TransformFiles(filesData)
		missing := utils.MissingFiles(localFiles, newFiles)
		for _, file := range missing {
			file.Filing = filing
			err := db.InsertFile(&file)
			if err != nil {
				panic(err)
			}
			time.Sleep(100 * time.Millisecond)
			file.Content, err = request.GetFileContent(&file)
			err = archive.PutObject(&file)
			if err != nil {
				panic(err)
			}
		}
	}
	for {
		companies, err := db.GetCompanies()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		for _, company := range companies {
			localFilings, err := db.GetFilings(&company)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			time.Sleep(100 * time.Millisecond)
			filingsData, err := request.GetFilingsData(&company)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			newFilings := utils.TransfromFilings(filingsData)
			missing := utils.MissingFilings(localFilings, newFilings)
			for _, filing := range missing {
				filing.Company = &company
				err := db.InsertFiling(&filing)
				if err != nil {
					log.Println(err.Error())
					continue
				}
				time.Sleep(100 * time.Millisecond)
				filesData, err := request.GetFilesData(&filing)
				if err != nil {
					log.Println(err.Error())
					continue
				}
				files := utils.TransformFiles(filesData)
				for _, file := range files {
					file.Filing = &filing
					err := db.InsertFile(&file)
					if err != nil {
						panic(err)
					}
					time.Sleep(100 * time.Millisecond)
					file.Content, err = request.GetFileContent(&file)
					err = archive.PutObject(&file)
					if err != nil {
						panic(err)
					}
				}
			}
		}
	}
}

func init() {
	connParams, err := storage.GetAWSConnParams()
	if err != nil {
		panic(err)
	}
	db, err = storage.NewDB(connParams)
	if err != nil {
		panic(err)
	}
	archive, err = storage.NewBucket()
	if err != nil {
		panic(err)
	}
}
