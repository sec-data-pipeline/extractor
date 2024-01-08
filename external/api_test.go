package external

import (
	"errors"
	"net/http"
	"testing"
)

func TestGetID(t *testing.T) {
	var tests = []struct {
		name string
		id   string
		want string
	}{
		{"Two hyphens in the middle", "234234-23423-4234", "234234234234234"},
		{"Only hyphens", "------", ""},
		{"Hyphen at the start", "-432523509", "432523509"},
		{"Hyphen at the end", "432523509-", "432523509"},
		{"Hyphen at start and end", "-432523509-", "432523509"},
		{"Empty ID", "", ""},
		{"Only one hyphen", "-", ""},
		{"Multiple hyphens at the start", "----2349234239", "2349234239"},
		{"Multiple hyphens at the end", "2349234239-----", "2349234239"},
		{"Hyphen at start, end and middle", "-2348723-324324-", "2348723324324"},
		{"No hyphen", "2343289472395", "2343289472395"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fil := &Filing{secID: test.id}
			got := fil.GetID()
			if got != test.want {
				t.Errorf("got %s, want %s", got, test.want)
			}
		})
	}
}

func TestGetExtension(t *testing.T) {
	var tests = []struct {
		name string
		want string
	}{
		{"test.htm", ".htm"},
		{"324324fasf", ""},
		{"hehe.pdf", ".pdf"},
		{"test.txt", ".txt"},
		{"sfsdafasdf.fasdfdsa.test", ".test"},
		{".htm", ".htm"},
		{".", "."},
		{"", ""},
		{"ewr897we.", "."},
		{".ewr897we.go", ".go"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fil := &file{Name: test.name}
			got, err := fil.GetExtension()
			if test.want == "" {
				if err == nil {
					t.Errorf("expected an error thrown, but got: %s", got)
				}
				return
			}
			if err != nil {
				t.Errorf(err.Error())
				return
			}
			if got != test.want {
				t.Errorf("got %s, want %s", got, test.want)
			}
		})
	}
}

func TestGetFilings(t *testing.T) {

}

func TestGetMainFile(t *testing.T) {

}

type testClient struct {
	data  [][]byte
	index int
}

func (c *testClient) buildRequest(urlStr string) (*http.Request, error) {
	return nil, nil
}

func (c *testClient) sendRequest(req *http.Request) (*http.Response, error) {
	return nil, nil
}

func (c *testClient) getData(res *http.Response) ([]byte, error) {
	if c.index >= len(c.data) {
		return nil, errors.New("Test data out of range")
	}
	data := c.data[c.index]
	c.index++
	return data, nil
}

func newTestAPI(data [][]byte) *API {
	return &API{client: &testClient{data: data, index: 0}}
}
