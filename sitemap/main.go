package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	// Get domain address
	domain := flag.String("domain", "http://calhoun.io", "Domain address which will be parsed.")
	flag.Parse()

	// Check path url. http and :// should be included
	pathChecker(*domain)

	// Base client
	var client http.Client
	t, err := pageParseHTML(client, *domain)
	if err != nil {
		log.Fatal(err)
	}

	nr := strings.NewReader(t)
	links, err := parseHTML(nr)
	if err != nil {
		log.Fatal(err)
	}

	links = removeDuplicates(links)
	fmt.Println(links)

}

func pageParseHTML(c http.Client, p string) (string, error) {
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

func pathChecker(path string) {
	if p, pr := strings.Index(path, "http"), strings.Index(path, "://"); p == -1 || pr == -1 {
		panic("Full path is required")
	}
}

func removeDuplicates(ls []link) []link {
	cl := make(map[link]bool)
	keyList := []link{}

	for _, l := range ls {
		if _, ok := cl[l]; ok == false {
			cl[l] = true
			keyList = append(keyList, l)
		}

	}

	return keyList
}
