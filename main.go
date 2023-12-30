package main

import (
	"errors"
	"log"

	"github.com/sec-data-pipeline/extractor/filing"
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
		for _, company := range companies {
			filings, err := getMissingFilings(&company)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			for _, fil := range filings {
				err := extendFiling(&company, &fil)
				if err != nil {
					log.Println(err.Error())
					continue
				}
				err = db.InsertFiling(&fil)
				if err != nil {
					log.Println(err.Error())
					continue
				}
				err = archive.PutObject(&fil)
				if err != nil {
					log.Println(err.Error())
				}
			}
		}
	}
}

func getMissingFilings(company *filing.Company) ([]filing.Filing, error) {
	filingsData, err := filing.GetFilingsData(company.CIK)
	if err != nil {
		return nil, err
	}
	newFilings := filing.TransfromFilings(filingsData)
	localFilings, err := db.GetFilings(company)
	return filing.MissingFilings(localFilings, newFilings), nil
}

func extendFiling(company *filing.Company, fil *filing.Filing) error {
	fil.Company = company
	indexFile := fil.RawID + "-index.html"
	fileBytes, err := filing.GetFileContent(fil.Company.CIK, fil.SECID, indexFile)
	if err != nil {
		return err
	}
	mainFile, err := filing.GetMainFileName(fileBytes)
	if err != nil {
		return err
	}
	filesData, err := filing.GetFilesData(fil.Company.CIK, fil.SECID)
	if err != nil {
		return err
	}
	files := filing.TransformFiles(filesData)
	for _, v := range files {
		if mainFile == v.Name {
			fil.File = &v
			fil.File.Content, err = filing.GetFileContent(fil.Company.CIK, fil.SECID, fil.File.Name)
			if err != nil {
				log.Println(err.Error())
				fil.File = nil
			}
			break
		}
	}
	if fil.File == nil {
		return errors.New("Main file: '" + mainFile + "' or it's content wasn't found")
	}
	return nil
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
