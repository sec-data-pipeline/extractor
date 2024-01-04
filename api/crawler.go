package api

import (
	"errors"
	"strings"

	"golang.org/x/net/html"
)

func getMainFileName(b []byte) (string, error) {
	document, err := html.Parse(strings.NewReader(string(b)))
	if err != nil {
		return "", err
	}
	tables, err := getTables(document)
	for _, v := range tables {
		elem, err := findMatchingRow(v)
		if err != nil {
			continue
		}
		return getFile(elem)
	}
	return "", errors.New("Checked all rows in tables and found no match")
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

func findMatchingRow(table *html.Node) (*html.Node, error) {
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
		return nil, errors.New("No matching row was found")
	}
	return row, nil
}

func getFile(node *html.Node) (string, error) {
	result := ""
	var crawler func(node *html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.TextNode {
			if len(result) < 1 && len(node.Data) > 4 && node.Data[len(node.Data)-4:] == ".htm" {
				result = node.Data
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(node)
	if len(result) < 1 {
		return "", errors.New("String could not be found in selected row")
	}
	return result, nil
}
