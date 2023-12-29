package models

import "time"

type Company struct {
	ID  int
	CIK string
}

type Filing struct {
	ID             int
	Company        *Company
	SECID          string
	Form           string
	Size           int
	FilingDate     time.Time
	ReportDate     time.Time
	AcceptanceDate time.Time
}

type File struct {
	Name         string
	Filing       *Filing
	Size         int
	LastModified time.Time
	Content      []byte
}

type FilingsResponse struct {
	Name      string     `json:"name"`
	CIK       string     `json:"cik"`
	Tickers   []string   `json:"tickers"`
	Exchanges []string   `json:"exchanges"`
	Filings   APIFilings `json:"filings"`
}

type APIFilings struct {
	Recent APIRecent `json:"recent"`
}

type APIRecent struct {
	AccessionNumber    []string `json:"accessionNumber"`
	FilingDate         []string `json:"filingDate"`
	AcceptanceDateTime []string `json:"acceptanceDateTime"`
	ReportDate         []string `json:"reportDate"`
	Form               []string `json:"form"`
	Size               []int    `json:"size"`
}

type FilesResponse struct {
	Dir Directory `json:"directory"`
}

type Directory struct {
	Items []Item `json:"item"`
}

type Item struct {
	Name         string `json:"name"`
	Size         string `json:"size"`
	LastModified string `json:"last-modified"`
}
