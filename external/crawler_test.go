package external

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestGetMainFileName(t *testing.T) {
	var tests = []struct {
		name  string
		input []byte
		want  string
	}{
		{
			"Multiple first rows in one table",
			[]byte(
				`
					<body>
      			<table>
        			<caption>Table</caption>
        			<tr>
          			<td>1</td>
          			<td>foo.htm</td>
        			</tr>
        			<tr>
          			<td>1</td>
          			<td>bar.htm</td>
        			</tr>
        			<tr>
          			<td>3</td>
        			</tr>
      			</table>
  				</body>
				`,
			),
			"foo.htm",
		},
		{
			"Multiple valid files in first row",
			[]byte(
				`
					<body>
      			<table>
        			<caption>Table</caption>
        			<tr>
          			<td>1</td>
          			<td>foo.htm</td>
          			<td>bar.htm</td>
        			</tr>
        			<tr>
          			<td>1</td>
          			<td>baz.htm</td>
        			</tr>
        			<tr>
          			<td>3</td>
        			</tr>
      			</table>
  				</body>
				`,
			),
			"foo.htm",
		},
		{
			"First matchin row without valid file",
			[]byte(
				`
					<body>
      			<table>
        			<caption>Table</caption>
        			<tr>
          			<td>1</td>
        			</tr>
        			<tr>
          			<td>1</td>
          			<td>baz.htm</td>
        			</tr>
        			<tr>
          			<td>3</td>
        			</tr>
      			</table>
  				</body>
				`,
			),
			"",
		},
		{
			"No matching first row",
			[]byte(
				`
					<body>
      			<table>
        			<caption>Table</caption>
        			<tr>
          			<td>2</td>
          			<td>foo.htm</td>
        			</tr>
        			<tr>
          			<td>2</td>
          			<td>baz.htm</td>
        			</tr>
        			<tr>
          			<td>3</td>
          			<td>bar.htm</td>
        			</tr>
      			</table>
  				</body>
				`,
			),
			"",
		},
		{
			"First row in a second table",
			[]byte(
				`
					<body>
      			<table>
        			<caption>Table</caption>
        			<tr>
          			<td>2</td>
          			<td>foo.htm</td>
        			</tr>
        			<tr>
          			<td>2</td>
          			<td>baz.htm</td>
        			</tr>
        			<tr>
          			<td>3</td>
          			<td>bar.htm</td>
        			</tr>
      			</table>
						<table>
        			<caption>Table</caption>
        			<tr>
          			<td>2</td>
          			<td>foo.htm</td>
        			</tr>
        			<tr>
          			<td>1</td>
          			<td>baz.htm</td>
        			</tr>
        			<tr>
          			<td>3</td>
          			<td>bar.htm</td>
        			</tr>
      			</table>

  				</body>
				`,
			),
			"baz.htm",
		},
		{
			"Matching row not in a table",
			[]byte(
				`
					<body>
						<div>
							<tr>
          			<td>1</td>
          			<td>bar.htm</td>
        			</tr>
						</div>
  				</body>
				`,
			),
			"",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := getMainFileName(test.input)
			if err != nil {
				if len(test.want) > 0 {
					t.Errorf(err.Error())
				}
				return
			}
			if len(test.want) < 1 {
				t.Errorf("expected error to be trown, but got: %s", got)
				return
			}
			if got != test.want {
				t.Errorf("got: %s, want: %s", got, test.want)
			}
		})
	}
}

func TestGetFileName(t *testing.T) {
	var tests = []struct {
		name    string
		content string
		want    string
	}{
		{"Deep match", `<div><div>test.htm</div></div>`, "test.htm"},
		{"Not deep match", `<div>test.htm</div>`, "test.htm"},
		{
			"Multiple matches on same height",
			`<div><div>foo.htm</div><div>bar.htm</div></div>`,
			"foo.htm",
		},
		{
			"Multiple matches on differen height",
			`<div><div><div>foo.htm</div></div><div>bar.htm</div></div>`,
			"foo.htm",
		},
		{"Deep no match", `<div><div>test.go</div></div>`, ""},
		{"Not deep no match", `<div>test.go</div>`, ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			node, err := stringToHTML(test.content)
			if err != nil {
				t.Errorf(err.Error())
				return
			}
			got, err := getFileName(node)
			if err != nil {
				if len(test.want) > 0 {
					t.Errorf(err.Error())
				}
				return
			}
			if len(test.want) < 1 {
				t.Errorf("expected error to be trown, but got: %s", got)
				return
			}
			if got != test.want {
				t.Errorf("got: %s, want: %s", got, test.want)
			}
		})
	}
}

func TestGetFirstRow(t *testing.T) {
	var tests = []struct {
		name    string
		content string
		want    string
	}{
		{
			"Ordered table",
			testStringHTML[3],
			`
				<tr>
					<td>1</td>
				</tr>
			`,
		},
		{
			"Not ordered table",
			testStringHTML[4],
			`
				<tr>
					<td>1</td>
				</tr>
			`,
		},
		{
			"More than one matching row",
			testStringHTML[5],
			`
				<tr> 
					<td>1</td> 
					<td>foo</td> 
				</tr>
			`,
		},
		{
			"No matching row found",
			testStringHTML[0],
			"",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			node, err := stringToHTML(test.content)
			if err != nil {
				t.Errorf("Could not parse string to html")
				return
			}
			row, err := getFirstRow(node)
			if err != nil {
				if len(test.want) > 0 {
					t.Errorf(err.Error())
				}
				return
			}
			got, err := htmlToString(row)
			if err != nil {
				t.Errorf("Could not parse html to string")
				return
			}
			if len(test.want) < 1 {
				t.Errorf("expected error to be trown, but got: %s", got)
				return
			}
			if stripString(got) != stripString(test.want) {
				t.Errorf("got: %s, want: %s", stripString(got), stripString(test.want))
			}
		})
	}
}

func TestGetTables(t *testing.T) {
	var tests = []struct {
		name    string
		content string
		want    int
	}{
		{"No table", testStringHTML[0], 0},
		{"Multiple tables in different depths", testStringHTML[1], 4},
		{"One table", testStringHTML[2], 1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			node, err := stringToHTML(test.content)
			if err != nil {
				t.Errorf("Could not parse string to html")
				return
			}
			got, err := getTables(node)
			if err != nil {
				t.Errorf(err.Error())
				return
			}
			if len(got) != test.want {
				t.Errorf("list length got: %d, want: %d", len(got), test.want)
			}
		})
	}
}

func stringToHTML(str string) (*html.Node, error) {
	document, err := html.Parse(strings.NewReader(str))
	if err != nil {
		return nil, err
	}
	return document, nil
}

func htmlToString(node *html.Node) (string, error) {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	err := html.Render(w, node)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func stripString(str string) string {
	result := strings.Replace(str, "\n", "", -1)
	result = strings.Replace(result, "\t", "", -1)
	return strings.Replace(result, " ", "", -1)
}

var testStringHTML []string = []string{
	`
		<div>
			<h1>Hello world!</h1>
			<p>This is sample 1.</p>
			<p>
				<a href="https://www.dwsamplefiles.com/">Learn More</a>
			</p>
		</div>
	`,
	`
  	<table>
    	<caption>Table 1</caption>
    	<tr>
      	<th>Name</th>
      	<th>Age</th>
    	</tr>
    	<tr>
      	<td>John</td>
      	<td>25</td>
    	</tr>
    	<tr>
      	<td>Jane</td>
      	<td>30</td>
    	</tr>
  	</table>
  	<div>
    	<table>
      	<caption>Table 2</caption>
      	<tr>
        	<th>City</th>
        	<th>Population</th>
      	</tr>
      	<tr>
        	<td>New York</td>
        	<td>8.4 million</td>
      	</tr>
      	<tr>
        	<td>Paris</td>
        	<td>2.1 million</td>
      	</tr>
    	</table>
  	</div>
  	<section>
    	<table>
      	<caption>Table 3</caption>
      	<tr>
        	<th>Product</th>
        	<th>Price</th>
      	</tr>
      	<tr>
        	<td>Laptop</td>
        	<td>$1000</td>
      	</tr>
      	<tr>
        	<td>Phone</td>
        	<td>$500</td>
      	</tr>
    	</table>
  	</section>
  	<table>
    	<caption>Table 4</caption>
    	<tr>
      	<th>Language</th>
      	<th>Framework</th>
    	</tr>
    	<tr>
      	<td>JavaScript</td>
      	<td>React</td>
    	</tr>
    	<tr>
      	<td>Python</td>
      	<td>Flask</td>
    	</tr>
  	</table>
	`,
	`
		<main>
    	<section>
      	<h2>Section 1</h2>
      	<p>Some content before the table...</p>
      	<table>
        	<caption>Depth Table</caption>
        	<tr>
          	<th>Name</th>
          	<th>Age</th>
          	<th>City</th>
        	</tr>
        	<tr>
          	<td>John</td>
          	<td>25</td>
          	<td>New York</td>
        	</tr>
        	<tr>
          	<td>Jane</td>
          	<td>30</td>
          	<td>Paris</td>
        	</tr>
      	</table>
      	<p>Some content after the table...</p>
    	</section>
    	<section>
      	<h2>Section 2</h2>
      	<p>Another section with different content...</p>
    	</section>
  	</main>
	`,
	`
		<body>
      <table>
        <caption>Table</caption>
        <tr>
          <td>1</td>
        </tr>
        <tr>
          <td>2</td>
        </tr>
        <tr>
          <td>3</td>
        </tr>
      </table>
  	</body>
	`,
	`
		<body>
      <table>
        <caption>Table</caption>
        <tr>
          <td>2</td>
        </tr>
        <tr>
          <td>1</td>
        </tr>
        <tr>
          <td>3</td>
        </tr>
      </table>
  	</body>
	`,
	`
		<body>
      <table>
        <caption>Table</caption>
        <tr>
          <td>1</td>
          <td>foo</td>
        </tr>
        <tr>
          <td>1</td>
          <td>bar</td>
        </tr>
        <tr>
          <td>3</td>
        </tr>
      </table>
  	</body>
	`,
}
