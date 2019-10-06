package main

import (
	"flag"
	"fmt"
	"net/http"
)

func main() {

	// Parse yaml file if it exists
	var yamlPath = flag.String("yaml", "", "parse yaml file")
	var jsonPath = flag.String("json", "", "parse json file")
	flag.Parse()

	if string(*yamlPath) != "" && string(*jsonPath) != "" {
		panic("Sorry, only yaml or json can be parsed in one time")
	}

	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := mapHandler(pathsToUrls, mux)

	// Build the Handler using the mapHandler as the fallback
	handler, err := combinedHandler(*yamlPath, *jsonPath, mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
