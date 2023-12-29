package utils

import (
	"strconv"
	"strings"
	"time"

	"github.com/sec-data-pipeline/extractor/models"
)

func TransfromFilings(data *models.FilingsResponse) []models.Filing {
	filings := []models.Filing{}
	for i, v := range data.Filings.Recent.Form {
		if len(v) < 4 {
			continue
		}
		if v[:3] != "10-" {
			continue
		}
		filing := models.Filing{}
		filing.SECID = transformID(data.Filings.Recent.AccessionNumber[i])
		filing.Size = data.Filings.Recent.Size[i]
		filing.Form = v
		fd, err := time.Parse("2006-01-02", data.Filings.Recent.FilingDate[i])
		if err != nil {
			fd = time.Time{}
		}
		filing.FilingDate = fd
		ad, err := time.Parse(time.RFC3339, data.Filings.Recent.AcceptanceDateTime[i])
		if err != nil {
			ad = time.Time{}
		}
		filing.AcceptanceDate = ad
		rd, err := time.Parse("2006-01-02", data.Filings.Recent.ReportDate[i])
		if err != nil {
			rd = time.Time{}
		}
		filing.ReportDate = rd
		filings = append(filings, filing)
	}
	return filings
}

func TransformFiles(data *models.FilesResponse) []models.File {
	files := []models.File{}
	for _, v := range data.Dir.Items {
		file := models.File{Name: v.Name}
		size, err := strconv.Atoi(v.Size)
		if err != nil {
			size = 0
		}
		file.Size = size
		lm, err := time.Parse("2006-01-02 15:04:05", v.LastModified)
		if err != nil {
			lm = time.Time{}
		}
		file.LastModified = lm
		files = append(files, file)
	}
	return files
}

func transformID(apiID string) string {
	return strings.Replace(apiID, "-", "", -1)
}
