package api

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func transformFilings(data *filingsResponse) []*Filing {
	var filings []*Filing
	for i, v := range data.Filings.Recent.Form {
		if v != "10-K" && v != "10-Q" {
			continue
		}
		var tmp Filing
		tmp.SECID = transformID(data.Filings.Recent.AccessionNumber[i])
		tmp.rawID = data.Filings.Recent.AccessionNumber[i]
		tmp.Form = v
		fd, err := time.Parse("2006-01-02", data.Filings.Recent.FilingDate[i])
		if err != nil {
			fd = time.Time{}
		}
		tmp.FilingDate = fd
		ad, err := time.Parse(time.RFC3339, data.Filings.Recent.AcceptanceDateTime[i])
		if err != nil {
			ad = time.Time{}
		}
		tmp.AcceptanceDate = ad
		rd, err := time.Parse("2006-01-02", data.Filings.Recent.ReportDate[i])
		if err != nil {
			rd = time.Time{}
		}
		tmp.ReportDate = rd
		filings = append(filings, &tmp)
	}
	return filings
}

func transformFiles(data *filesResponse) []*File {
	var files []*File
	for _, v := range data.Dir.Items {
		var tmp File
		tmp.Name = v.Name
		size, err := strconv.Atoi(v.Size)
		if err != nil {
			tmp.Size = -1
		} else {
			tmp.Size = size
		}
		lm, err := time.Parse("2006-01-02 15:04:05", v.LastModified)
		if err != nil {
			lm = time.Time{}
		}
		tmp.LastModified = lm
		files = append(files, &tmp)
	}
	return files
}

func transformID(apiID string) string {
	return strings.Replace(apiID, "-", "", -1)
}

func missingFilings(filings []*Filing, ids []string) []*Filing {
	var missing []*Filing
outer:
	for _, flng := range filings {
		for _, id := range ids {
			if flng.SECID == id {
				continue outer
			}
		}
		missing = append(missing, flng)
	}
	return missing
}

func getFileExtension(fileName string) (string, error) {
	if !strings.Contains(fileName, ".") {
		return "", errors.New("File: '" + fileName + "' does not have a file extension")
	}
	result := ""
	for i := len(fileName) - 1; i >= 0; i-- {
		result = string(fileName[i]) + result
		if string(fileName[i]) == "." {
			break
		}
	}
	return result, nil
}
