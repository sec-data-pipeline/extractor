package api

import "time"

type Filing struct {
	SECID          string
	rawID          string
	Form           string
	FilingDate     time.Time
	ReportDate     time.Time
	AcceptanceDate time.Time
}

type File struct {
	Name         string
	Content      []byte
	Extension    string
	Size         int
	LastModified time.Time
}

type filingsResponse struct {
	Name    string     `json:"name"`
	CIK     string     `json:"cik"`
	Filings apiFilings `json:"filings"`
}

type apiFilings struct {
	Recent apiRecent `json:"recent"`
}

type apiRecent struct {
	AccessionNumber    []string `json:"accessionNumber"`
	FilingDate         []string `json:"filingDate"`
	AcceptanceDateTime []string `json:"acceptanceDateTime"`
	ReportDate         []string `json:"reportDate"`
	Form               []string `json:"form"`
}

type filesResponse struct {
	Dir directory `json:"directory"`
}

type directory struct {
	Items []item `json:"item"`
}

type item struct {
	Name         string `json:"name"`
	Size         string `json:"size"`
	LastModified string `json:"last-modified"`
}
