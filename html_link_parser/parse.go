package main

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

type link struct {
	Href string
	Text string
}

// Split taked nodes and returl list of selected
func splitNodes(n *html.Node, t string) []*html.Node {
	if n.Type == html.ElementNode && n.Data == t {
		return []*html.Node{n}
	}

	nodes := []*html.Node{}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, splitNodes(c, t)...)
	}

	return nodes
}

func parse(r io.Reader, t string) ([]link, error) {
	nodesObjs, err := html.Parse(r)

	if err != nil {
		panic(err)
	}

	nodes := splitNodes(nodesObjs, t)
	result := []link{}

	for _, node := range nodes {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				result = append(result, link{Href: attr.Val, Text: getLinkText(node)})
				break
			}
		}
	}

	return result, nil
}

func getLinkText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}

	text := ""
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += getLinkText(c)
	}
	return strings.Join(strings.Fields(text), " ")
}
