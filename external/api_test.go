package external

import (
	"database/sql"
	"errors"
	"net/http"
	"testing"
	"time"
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
	var tests = []struct {
		name    string
		mockRes []byte
		err     error
		want    []*Filing
	}{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			api := newTestAPI([][]byte{test.mockRes})
			got, err := api.GetFilings("")
			if err != nil && test.err == nil {
				t.Errorf(err.Error())
				return
			}
			if test.err != nil && err == nil {
				t.Errorf("expected error to be trown, but got list of length: %d", len(got))
				return
			}
			if test.err != nil && err != nil {
				return
			}
			for i, v := range got {
				if v.secID != test.want[i].secID {
					t.Errorf("got %s, want %s", v.secID, test.want[i].secID)
				}
				if v.Form != test.want[i].Form {
					t.Errorf("got %s, want %s", v.Form, test.want[i].Form)
				}
				msg, ok := checkNullTime(v.AcceptDate, test.want[i].AcceptDate)
				if !ok {
					t.Errorf(msg)
				}
				msg, ok = checkNullTime(v.ReportDate, test.want[i].ReportDate)
				if !ok {
					t.Errorf(msg)
				}
				msg, ok = checkNullTime(v.FilingDate, test.want[i].FilingDate)
				if !ok {
					t.Errorf(msg)
				}
			}
		})
	}
}

func TestGetMainFile(t *testing.T) {
	type mock struct {
		firstRes  []byte
		secondRes []byte
		thirdeRes []byte
	}
	var mocks []mock = []mock{
		{
			firstRes: []byte(`
				{
					"directory":{
						"item":[
							{
								"last-modified": "2004-09-10 16:47:30",
								"name":"ex32.txt"
							},
							{
								"last-modified": "2004-09-10 16:47:30",
								"name":"k2004.htm"
							}
						]
					}
				}
			`),
			secondRes: []byte(`
				<table>
					<tr>
						<td>1</td>
						<td>k2004.htm</td>
					</tr>
				</table>
			`),
			thirdeRes: []byte(`test`),
		},
		{
			firstRes: []byte(`
				{
					"directory":{
						"item":[
							{
								"last-modified": "2004-09-10 16:47:30",
								"name":"ex32.txt"
							},
							{
								"last-modified": "2004-09-10 16:47:30",
								"name":"k2004.htm"
							}
						]
					}
				}
			`),
			secondRes: []byte(`
				<table>
					<tr>
						<td>2</td>
						<td>k2004.htm</td>
					</tr>
				</table>
			`),
			thirdeRes: []byte(`test`),
		},
		{
			firstRes: []byte(`
				{
					"directory":{
						"item":[
							{
								"last-modified": "2004-09-10 16:47:30",
								"name":"ex32.txt"
							}
						]
					}
				}
			`),
			secondRes: []byte(`
				<table>
					<tr>
						<td>1</td>
						<td>k2004.htm</td>
					</tr>
				</table>
			`),
			thirdeRes: []byte(`test`),
		},
		{
			firstRes: []byte(`
				{
					"directory":{
						{
							"last-modified": "2004-09-10 16:47:30",
							"name":"ex32.txt"
						}
					}
				}
			`),
			secondRes: []byte(``),
			thirdeRes: []byte(`test`),
		},
	}
	var tests = []struct {
		name    string
		mockRes mock
		err     error
		want    *file
	}{
		{
			"Found valid file",
			mocks[0],
			nil,
			&file{
				Name: "k2004.htm",
				LastModified: sql.NullTime{
					Time:  time.Date(2004, time.September, 10, 16, 47, 30, 0, time.UTC),
					Valid: true,
				},
				Content: mocks[0].thirdeRes,
			},
		},
		{"Main file not in first row", mocks[1], errors.New(""), &file{}},
		{"Main file not in file list", mocks[2], errors.New(""), &file{}},
		{"Invalid JSON", mocks[3], errors.New(""), &file{}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			api := newTestAPI(
				[][]byte{test.mockRes.firstRes, test.mockRes.secondRes, test.mockRes.thirdeRes},
			)
			got, err := api.GetMainFile("", &Filing{})
			if err != nil && test.err == nil {
				t.Errorf(err.Error())
				return
			}
			if test.err != nil && err == nil {
				t.Errorf("expected error to be trown, but got: %s", got.Name)
				return
			}
			if test.err != nil && err != nil {
				return
			}
			if got.Name != test.want.Name {
				t.Errorf("got %s, want %s", got.Name, test.want.Name)
			}
			msg, ok := checkNullTime(got.LastModified, test.want.LastModified)
			if !ok {
				t.Errorf(msg)
			}
			if string(got.Content) != string(test.want.Content) {
				t.Errorf("got: %s, want: %s", string(got.Content), string(test.want.Content))
			}
		})
	}
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
