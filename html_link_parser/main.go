package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	htmlPath := flag.String("html", "templates/ex1.html", "html file")
	flag.Parse()

	r := strings.NewReader(readHTML(*htmlPath))
	t := "a" //Type of nodes what we want to search
	links, err := parse(r, t)

	if err != nil {
		panic(err)
	}

	for _, link := range links {
		fmt.Println("href: ", link.Href, "text: ", link.Text)
	}

}

func readHTML(f string) string {
	file, err := ioutil.ReadFile(f)

	if err != nil {
		panic(err)
	}

	return string(file)
}
