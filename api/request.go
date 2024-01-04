package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

const (
	baseURL = "https://www.sec.gov/Archives/edgar/data/"
	apiURL  = "https://data.sec.gov/submissions/CIK"
)

func getFilingsData(cik string) (*filingsResponse, error) {
	req, err := buildRequest(apiURL + cik + ".json")
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

func getFilesData(cik string, secID string) (*filesResponse, error) {
	req, err := buildRequest(
		baseURL + cik + "/" + secID + "/index.json",
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

func getFileContent(cik string, secID string, name string) ([]byte, error) {
	req, err := buildRequest(
		baseURL + cik + "/" + secID + "/" + name,
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

func handleFilingsResponse(res *http.Response) (*filingsResponse, error) {
	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	body := &filingsResponse{}
	if err = json.Unmarshal(bodyBytes, &body); err != nil {
		return nil, errors.New(
			"Could not process JSON into struct filingsResponse\n" +
				"json.Unmarshal threw: '" + err.Error() + "'\n" + string(bodyBytes),
		)
	}
	return body, nil
}

func handleFilesResponse(res *http.Response) (*filesResponse, error) {
	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	body := filesResponse{}
	if err = json.Unmarshal(bodyBytes, &body); err != nil {
		return nil, errors.New(
			"Could not process JSON into struct filesResponse\n" +
				"json.Unmarshal threw: '" + err.Error() + "'\n" + string(bodyBytes),
		)
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
	time.Sleep(100 * time.Millisecond)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
