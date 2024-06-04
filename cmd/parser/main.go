package main

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {
	err := parseStations()
	if err != nil {
		log.Fatal(err)
	}
}

func parseStations() error {
	resp, err := http.Get("https://ru.wikipedia.org/wiki/%D0%A1%D0%BF%D0%B8%D1%81%D0%BE%D0%BA_%D1%81%D1%82%D0%B0%D0%BD%D1%86%D0%B8%D0%B9_%D0%9C%D0%BE%D1%81%D0%BA%D0%BE%D0%B2%D1%81%D0%BA%D0%BE%D0%B3%D0%BE_%D0%BC%D0%B5%D1%82%D1%80%D0%BE%D0%BF%D0%BE%D0%BB%D0%B8%D1%82%D0%B5%D0%BD%D0%B0")
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("bad status code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	doc, err := html.Parse(strings.NewReader(string(body)))
	table, err := findTable(doc)
	if err != nil {
		return err
	}
	if table == nil {
		return fmt.Errorf("table not found")
	}
	fmt.Println(table.Data)

	return nil
}

func getStrNodeContent(n *html.Node) (string, error) {
	buf := &bytes.Buffer{}
	err := html.Render(buf, n)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(buf.String()), nil
}

func findTable(n *html.Node) (*html.Node, error) {
	if n.Type == html.ElementNode && n.Data == "caption" {
		content, err := getStrNodeContent(n)
		if err != nil {
			return nil, err
		}
		if strings.Contains(
			content,
			"Список может быть отсортирован по названиям станций в алфавитном порядке",
		) {
			return n.Parent, nil
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		table, err := findTable(c)
		if err != nil {
			return nil, err
		}
		if table != nil {
			return table, nil
		}
	}

	return nil, nil
}
