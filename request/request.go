package request

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/sec-data-pipeline/extractor/models"
)

const (
	baseURL = "https://www.sec.gov/Archives/edgar/data/"
	apiURL  = "https://data.sec.gov/submissions/CIK"
)

func GetFilingsData(company *models.Company) (*models.FilingsResponse, error) {
	req, err := buildRequest(apiURL + company.CIK + ".json")
	if err != nil {
		return nil, err
	}
	res, err := sendRequest(req)
	if err != nil {
		return nil, err
	}
	data, err := handleFilingsResponse(res)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetFilesData(filing *models.Filing) (*models.FilesResponse, error) {
	req, err := buildRequest(
		baseURL + filing.Company.CIK + "/" + filing.SECID + "/index.json",
	)
	if err != nil {
		return nil, err
	}
	res, err := sendRequest(req)
	if err != nil {
		return nil, err
	}
	data, err := handleFilesResponse(res)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetFileContent(file *models.File) ([]byte, error) {
	req, err := buildRequest(
		baseURL + file.Filing.Company.CIK + "/" + file.Filing.SECID + "/" + file.Name,
	)
	if err != nil {
		return nil, err
	}
	res, err := sendRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func handleFilingsResponse(res *http.Response) (*models.FilingsResponse, error) {
	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	body := &models.FilingsResponse{}
	if err = json.Unmarshal(bodyBytes, &body); err != nil {
		return nil, err
	}
	return body, nil
}

func handleFilesResponse(res *http.Response) (*models.FilesResponse, error) {
	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	body := models.FilesResponse{}
	if err = json.Unmarshal(bodyBytes, &body); err != nil {
		return nil, err
	}
	return &body, nil
}

func buildRequest(urlStr string) (*http.Request, error) {
	req, err := http.NewRequest("GET", urlStr, nil)
	req.Header.Add("User-Agent", "example.com info@example.com")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Connection", "keep-alive")
	if err != nil {
		return nil, err
	}
	return req, nil
}

func sendRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
