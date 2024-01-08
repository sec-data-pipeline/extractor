package external

import (
	"io"
	"net/http"
	"time"
)

type client interface {
	buildRequest(urlStr string) (*http.Request, error)
	sendRequest(req *http.Request) (*http.Response, error)
	getData(res *http.Response) ([]byte, error)
}

type webClient struct{}

func newWebClient() *webClient {
	return &webClient{}
}

func (c *webClient) buildRequest(urlStr string) (*http.Request, error) {
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "example.com info@example.com")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Connection", "keep-alive")
	return req, nil
}

func (c *webClient) sendRequest(req *http.Request) (*http.Response, error) {
	time.Sleep(200 * time.Millisecond)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *webClient) getData(res *http.Response) ([]byte, error) {
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
