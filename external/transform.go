package external

import (
	"database/sql"
	"time"
)

func transformFilings(data *filingsResponse) []*Filing {
	var filings []*Filing
	for i, v := range data.Filings.Recent.Form {
		if v != "10-K" && v != "10-Q" {
			continue
		}
		test := &file{Name: data.Filings.Recent.PrimDoc[i]}
		check, err := test.GetExtension()
		if err != nil {
			continue
		}
		if check != ".htm" {
			continue
		}
		fil := &Filing{
			secID:      data.Filings.Recent.AccessNumber[i],
			mainFile:   data.Filings.Recent.PrimDoc[i],
			Form:       v,
			FilingDate: parseNullTime("2006-01-02", data.Filings.Recent.FilingDate[i]),
			AcceptDate: parseNullTime(time.RFC3339, data.Filings.Recent.AcceptDate[i]),
			ReportDate: parseNullTime("2006-01-02", data.Filings.Recent.ReportDate[i]),
		}
		filings = append(filings, fil)
	}
	return filings
}

func transformFiles(data *filesResponse) []*file {
	var files []*file
	for _, v := range data.Dir.Items {
		fil := &file{
			Name:         v.Name,
			LastModified: parseNullTime("2006-01-02 15:04:05", v.LastModified),
		}
		files = append(files, fil)
	}
	return files
}

func parseNullTime(layout string, value string) sql.NullTime {
	t, err := time.Parse(layout, value)
	if err != nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: t, Valid: true}
}
