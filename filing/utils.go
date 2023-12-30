package filing

import (
	"strings"
	"time"
)

func TransfromFilings(data *FilingsResponse) []Filing {
	filings := []Filing{}
	for i, v := range data.Filings.Recent.Form {
		if len(v) < 4 {
			continue
		}
		if v[:3] != "10-" {
			continue
		}
		filing := Filing{}
		filing.SECID = transformID(data.Filings.Recent.AccessionNumber[i])
		filing.RawID = data.Filings.Recent.AccessionNumber[i]
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

func TransformFiles(data *FilesResponse) []File {
	files := []File{}
	for _, v := range data.Dir.Items {
		file := File{Name: v.Name}
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

func MissingFilings(localFilings []Filing, newFilings []Filing) []Filing {
	missing := []Filing{}
outer:
	for _, newFiling := range newFilings {
		for _, localFiling := range localFilings {
			if newFiling.SECID == localFiling.SECID {
				continue outer
			}
		}
		missing = append(missing, newFiling)
	}
	return missing
}

func GetFileExtension(fileName string) string {
	result := ""
	for i := len(fileName) - 1; i >= 0; i-- {
		result = string(fileName[i]) + result
		if string(fileName[i]) == "." {
			break
		}
	}
	return result
}
