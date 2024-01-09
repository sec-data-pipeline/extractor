package external

type filingsResponse struct {
	Name    string  `json:"name"`
	CIK     string  `json:"cik"`
	Filings filings `json:"filings"`
}

type filings struct {
	Recent recent `json:"recent"`
}

type recent struct {
	AccessNumber []string `json:"accessionNumber"`
	FilingDate   []string `json:"filingDate"`
	AcceptDate   []string `json:"acceptanceDateTime"`
	ReportDate   []string `json:"reportDate"`
	Form         []string `json:"form"`
}

type filesResponse struct {
	Dir directory `json:"directory"`
}

type directory struct {
	Items []item `json:"item"`
}

type item struct {
	Name         string `json:"name"`
	LastModified string `json:"last-modified"`
}
