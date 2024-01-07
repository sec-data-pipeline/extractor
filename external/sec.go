package external

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type secAPI struct {
	baseURL string
	dataURL string
	request *struct {
		client  *http.Client
		headers map[string]string
		delay   time.Duration
	}
}

func NewSECAPI() *secAPI {
	return &secAPI{
		baseURL: "https://www.sec.gov/Archives/edgar/data/",
		dataURL: "https://data.sec.gov/submissions/CIK",
		request: &struct {
			client  *http.Client
			headers map[string]string
			delay   time.Duration
		}{
			client: &http.Client{},
			headers: map[string]string{
				"User-Agent": "example.com info@example.com",
				"Accept":     "*/*",
				"Connection": "keep-alive",
			},
			delay: 200,
		},
	}
}

func (api *secAPI) GetFilings(cik string) ([]*Filing, error) {
	req, err := api.buildRequest(api.dataURL + cik + ".json")
	if err != nil {
		return nil, err
	}
	res, err := api.sendRequest(req)
	if err != nil {
		return nil, err
	}
	data, err := api.getData(res)
	if err != nil {
		return nil, err
	}
	filRes := &filingsResponse{}
	if err := json.Unmarshal(data, filRes); err != nil {
		return nil, errors.New("Could not process JSON into struct filingsResponse, " + err.Error())
	}
	return api.transformFilings(filRes), nil
}

func (api *secAPI) GetMainFile(cik string, fil *Filing) (*file, error) {
	req, err := api.buildRequest(api.baseURL + cik + "/" + fil.GetID() + "/index.json")
	if err != nil {
		return nil, err
	}
	res, err := api.sendRequest(req)
	if err != nil {
		return nil, err
	}
	data, err := api.getData(res)
	if err != nil {
		return nil, err
	}
	filRes := &filesResponse{}
	if err := json.Unmarshal(data, filRes); err != nil {
		return nil, errors.New("Could not process JSON into struct filesResponse, " + err.Error())
	}
	indexFile := fil.secID + "-index.html"
	content, err := api.getFileContent(cik, fil.secID, indexFile)
	if err != nil {
		return nil, errors.New("Could not get content of index file, " + err.Error())
	}
	mfName, err := api.getMainFileName(content)
	if err != nil {
		return nil, err
	}
	files := api.transformFiles(filRes)
	mainFile, err := getFile(files, mfName)
	if err != nil {
		return nil, err
	}
	mainFile.Content, err = api.getFileContent(cik, fil.GetID(), mfName)
	if err != nil {
		return nil, errors.New("Could not get content of main file, " + err.Error())
	}
	return mainFile, nil
}

func (api *secAPI) getMainFileName(data []byte) (string, error) {
	document, err := html.Parse(strings.NewReader(string(data)))
	if err != nil {
		return "", err
	}
	tables, err := api.getTables(document)
	for _, table := range tables {
		row, err := api.getFirstRow(table)
		if err != nil {
			continue
		}
		return api.getFileName(row)
	}
	return "", errors.New("Checked all rows in tables and found no match")
}

func (api *secAPI) getFileName(node *html.Node) (string, error) {
	var fileName string = ""
	var crawler func(node *html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.TextNode {
			if len(fileName) < 1 && len(node.Data) > 4 && node.Data[len(node.Data)-4:] == ".htm" {
				fileName = node.Data
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(node)
	if len(fileName) < 1 {
		return "", errors.New("String could not be found in selected row")
	}
	return fileName, nil
}

func (api *secAPI) getFirstRow(table *html.Node) (*html.Node, error) {
	var row *html.Node
	var crawler func(node *html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.TextNode && node.Data == "1" {
			row = node.Parent.Parent
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(table)
	if row == nil {
		return nil, errors.New("Table does not contain the first row")
	}
	return row, nil
}

func (api *secAPI) getTables(document *html.Node) ([]*html.Node, error) {
	tables := []*html.Node{}
	var crawler func(node *html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "table" {
			tables = append(tables, node)
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(document)
	return tables, nil
}

func (api *secAPI) transformFilings(data *filingsResponse) []*Filing {
	var filings []*Filing
	for i, v := range data.Filings.Recent.Form {
		if v != "10-K" && v != "10-Q" {
			continue
		}
		var tmp Filing
		tmp.secID = data.Filings.Recent.AccessionNumber[i]
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

func (api *secAPI) transformFiles(data *filesResponse) []*file {
	var files []*file
	for _, v := range data.Dir.Items {
		var tmp file
		tmp.Name = v.Name
		lm, err := time.Parse("2006-01-02 15:04:05", v.LastModified)
		if err != nil {
			lm = time.Time{}
		}
		tmp.LastModified = lm
		files = append(files, &tmp)
	}
	return files
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

func (api *secAPI) getFileContent(cik string, secID string, name string) ([]byte, error) {
	req, err := api.buildRequest(
		api.baseURL + cik + "/" + secID + "/" + name,
	)
	if err != nil {
		return nil, err
	}
	res, err := api.sendRequest(req)
	if err != nil {
		return nil, err
	}
	data, err := api.getData(res)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (api *secAPI) getData(res *http.Response) ([]byte, error) {
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (api *secAPI) buildRequest(urlStr string) (*http.Request, error) {
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range api.request.headers {
		req.Header.Add(k, v)
	}
	return req, nil
}

func (api *secAPI) sendRequest(req *http.Request) (*http.Response, error) {
	time.Sleep(api.request.delay * time.Millisecond)
	res, err := api.request.client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
