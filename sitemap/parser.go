package main

import (
	"io"

	"golang.org/x/net/html"
)

type link struct {
	Href string
}

// Check all nodes and return list of selected nodes
func splitNodes(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}

	nodes := []*html.Node{}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, splitNodes(c)...)
	}

	return nodes
}

// Parse HTML page and return all <a> links on page
func parseHTML(r io.Reader) ([]link, error) {
	nodesObjs, err := html.Parse(r)

	if err != nil {
		panic(err)
	}

	nodes := splitNodes(nodesObjs)
	result := []link{}

	for _, node := range nodes {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				result = append(result, link{Href: attr.Val})
				break
			}
		}
	}

	return result, nil
}
