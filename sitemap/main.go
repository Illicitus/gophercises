package main

import (
	"encoding/xml"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type locationList struct {
	Loc string `xml:"loc"`
}

type pageMap struct {
	XMLName xml.Name `xml:"urlset"`
	Path    string   `xml:"xmlns,attr"`
	URL []locationList `xml:"url"`
}

type domain struct {
	Prt    string
	Domain string
}

func newDomain(path string) domain {
	protocolSeparator := "://"

	// Check path url. http and :// should be included
	if p, pr := strings.Index(path, "http"), strings.Index(path, protocolSeparator); p == -1 || pr == -1 {
		panic("Full path is required")
	}

	urlData := strings.Split(path, protocolSeparator)
	return domain{Prt: urlData[0], Domain: strings.Split(urlData[1], "/")[0]}

}

func (d *domain) getFullPath() string {
	return d.Prt + "://" + d.Domain
}

func main() {
	// Get domain address
	d := flag.String("domain", "https://www.calhoun.io", "Domain address which will be parsed.")
	flag.Parse()

	// Create domain struct
	dom := newDomain(*d)
	domPath := dom.getFullPath()

	// Base client
	var client http.Client

	// Generate links from main page
	links := linksGenerator(client, dom, domPath)

	// Generate links from each base link
	result := make(map[string][]link, 0)
	result = getSiteMapUrls(client, dom, links, result)


	// Generate xml
	g := []pageMap{}
	for k, v := range result {
		valuesList := make([]locationList, 0)

		for _, obj := range v {
			valuesList = append(valuesList, locationList{Loc: dom.getFullPath()+obj.Href})
		}
		page := pageMap{Path: dom.getFullPath() + k}
		page.URL = valuesList
		g = append(g, page)

	}

	output, err := xml.MarshalIndent(g, "  ", "    ")
	if err != nil {
		log.Fatal(err)
	}
	os.Stdout.Write(output)

}

func getSiteMapUrls(client http.Client, dom domain, links []link, result map[string][]link) map[string][]link {
	for _, l := range links {
		if _, ok := result[l.Href]; ok == false {
			path := dom.getFullPath() + l.Href
			newLinks := linksGenerator(client, dom, path)
			result[l.Href] = newLinks

			if len(newLinks) != 0 {
				result = getSiteMapUrls(client, dom, newLinks, result)
			}
		}
	}
	return result
}

func getPage(c http.Client, p string) (string, error) {
	resp, err := c.Get(p)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		b := string(body)
		return b, nil
	}
	return "", err
}

func linksGenerator(c http.Client, d domain, p string) []link {

	t, err := getPage(c, p)
	if err != nil {
		log.Fatal(err)
	}

	nr := strings.NewReader(t)
	links, err := parseHTML(nr)
	if err != nil {
		log.Fatal(err)
	}

	// Remove duplicates and not main links
	links = linksCleaner(links, d)

	return links
}

func linksCleaner(ls []link, d domain) []link {
	domainPath := d.getFullPath()

	// Remove related sites
	onlyDomainLinks := []link{}
	for _, l := range ls {
		if string(l.Href[0]) == "/" {
			onlyDomainLinks = append(onlyDomainLinks, l)
		} else {
			if strings.Contains(l.Href, domainPath) {
				onlyDomainLinks = append(onlyDomainLinks, l)

			}
		}
	}

	// Convert to one format
	oneFormatLinks := []link{}
	for _, l := range onlyDomainLinks {
		if string(l.Href[0]) == "/" {

			oneFormatLinks = append(oneFormatLinks, l)

		} else if strings.Contains(l.Href, domainPath) {

			l.Href = l.Href[len(domainPath):]
			oneFormatLinks = append(oneFormatLinks, l)

		} else if strings.Contains(l.Href, d.Domain) {
			l.Href = l.Href[len(d.Domain):]
			oneFormatLinks = append(oneFormatLinks, l)
		}
	}

	// Remove duplicates
	cl := make(map[link]bool)
	uniqueLinks := []link{}

	for _, l := range oneFormatLinks {
		if _, ok := cl[l]; ok == false {
			cl[l] = true
			uniqueLinks = append(uniqueLinks, l)
		}

	}

	return uniqueLinks
}
