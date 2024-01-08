package external

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
)

type API struct {
	client    client
	fileURL   string
	filingURL string
}

func NewAPI() *API {
	return &API{
		client:    newWebClient(),
		fileURL:   "https://www.sec.gov/Archives/edgar/data/",
		filingURL: "https://data.sec.gov/submissions/CIK",
	}
}

func (api *API) GetFilings(cik string) ([]*Filing, error) {
	req, err := api.client.buildRequest(api.filingURL + cik + ".json")
	if err != nil {
		return nil, err
	}
	res, err := api.client.sendRequest(req)
	if err != nil {
		return nil, err
	}
	data, err := api.client.getData(res)
	if err != nil {
		return nil, err
	}
	filRes := &filingsResponse{}
	if err := json.Unmarshal(data, filRes); err != nil {
		return nil, errors.New("Could not process JSON into struct filingsResponse, " + err.Error())
	}
	return transformFilings(filRes), nil
}

func (api *API) GetMainFile(cik string, fil *Filing) (*file, error) {
	req, err := api.client.buildRequest(api.fileURL + cik + "/" + fil.GetID() + "/index.json")
	if err != nil {
		return nil, err
	}
	res, err := api.client.sendRequest(req)
	if err != nil {
		return nil, err
	}
	data, err := api.client.getData(res)
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
	mfName, err := getMainFileName(content)
	if err != nil {
		return nil, err
	}
	files := transformFiles(filRes)
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

func (api *API) getFileContent(cik string, secID string, name string) ([]byte, error) {
	req, err := api.client.buildRequest(api.fileURL + cik + "/" + secID + "/" + name)
	if err != nil {
		return nil, err
	}
	res, err := api.client.sendRequest(req)
	if err != nil {
		return nil, err
	}
	data, err := api.client.getData(res)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type Filing struct {
	secID      string
	Form       string
	FilingDate sql.NullTime
	ReportDate sql.NullTime
	AcceptDate sql.NullTime
}

func (f *Filing) GetID() string {
	return strings.Replace(f.secID, "-", "", -1)
}

type file struct {
	Name         string
	Content      []byte
	LastModified sql.NullTime
}

func (f *file) GetExtension() (string, error) {
	if !strings.Contains(f.Name, ".") {
		return "", errors.New("File extension could not be found")
	}
	result := ""
	for i := len(f.Name) - 1; i >= 0; i-- {
		result = string(f.Name[i]) + result
		if string(f.Name[i]) == "." {
			break
		}
	}
	return result, nil
}

func getFile(files []*file, name string) (*file, error) {
	for _, file := range files {
		if file.Name == name {
			return file, nil
		}
	}
	return nil, errors.New("File not in provided list")
}
