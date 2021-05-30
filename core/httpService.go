package core

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"golang.org/x/net/html"
)

type HttpService struct {
	Client http.Client
}

func NewHttpService() *HttpService {
	httpService := new(HttpService)
	httpService.Client = http.Client{}
	return httpService
}

func RenderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}

func Body(doc *html.Node) (*html.Node, error) {
	var body *html.Node
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "body" {
			body = node
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(doc)
	if body != nil {
		return body, nil
	}
	return nil, errors.New("Missing <body> in the node tree")
}

func (httpService HttpService) GetBody(URL string) (string, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	htmlTree, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}
	bodyNode, err := Body(htmlTree)
	if err != nil {
		return "", err
	}
	bodyString := RenderNode(bodyNode)
	return bodyString, nil
}
