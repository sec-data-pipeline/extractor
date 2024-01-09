package external

import (
	"errors"
	"strings"

	"golang.org/x/net/html"
)

func getMainFileName(data []byte) (string, error) {
	document, err := html.Parse(strings.NewReader(string(data)))
	if err != nil {
		return "", err
	}
	tables, err := getTables(document)
	for _, table := range tables {
		row, err := getFirstRow(table)
		if err != nil {
			continue
		}
		return getFileName(row)
	}
	return "", errors.New("No match in any table found")
}

func getFileName(node *html.Node) (string, error) {
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

func getFirstRow(table *html.Node) (*html.Node, error) {
	var row *html.Node
	var crawler func(node *html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.TextNode && node.Data == "1" && row == nil {
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

func getTables(document *html.Node) ([]*html.Node, error) {
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
