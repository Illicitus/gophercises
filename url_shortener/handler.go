package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"gopkg.in/yaml.v2"
)

type urlPattern struct {
	Path string
	URL  string
}

func combinedHandler(yml string, jsn string, fallback http.Handler) (http.HandlerFunc, error) {
	ymlFile := getYamlData(yml)
	yh, err := yamlHandler([]byte(ymlFile), fallback)

	if err != nil {
		panic(err)
	}

	jsnFile := getJSONData(jsn)
	jh, err := jsonHandler([]byte(jsnFile), fallback)

	if err != nil {
		panic(err)
	}

	patterns := append(yh, jh...)
	pathsMap := buildMap(patterns)
	return mapHandler(pathsMap, fallback), nil

}

func mapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	view := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		newPath, ok := pathsToUrls[path]

		if ok {
			http.Redirect(w, r, newPath, 301)
		} else {
			fallback.ServeHTTP(w, r)
		}

	}
	return view
}

// Return yaml from file if it exists
func getYamlData(f string) string {
	if f != "" {
		file, err := ioutil.ReadFile(f)

		if err != nil {
			panic(err)
		}

		return string(file)
	}

	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	return yaml
}

// Return json from file if it exists
func getJSONData(f string) string {
	if f != "" {
		file, err := ioutil.ReadFile(f)

		if err != nil {
			panic(err)
		}

		return string(file)
	}

	json := `[
		{
		  "path": "/json",
		  "url": "https://golang.org/pkg/encoding/json/"
		},
		{
		  "path":"/yahoo",
		  "url":"https://yahoo.com"
		}
	  ]`
	return json
}

func yamlHandler(yml []byte, fallback http.Handler) ([]urlPattern, error) {
	parsedYaml, err := parseYaml(yml)

	if err != nil {
		return nil, err
	}

	return parsedYaml, nil
}

func jsonHandler(jsn []byte, fallback http.Handler) ([]urlPattern, error) {
	parsedJSON, err := parseJSON(jsn)

	if err != nil {
		return nil, err
	}

	return parsedJSON, nil
}

// Parse yaml and return map of urlPattern
func parseYaml(yml []byte) ([]urlPattern, error) {
	patterns := []urlPattern{}
	err := yaml.Unmarshal(yml, &patterns)
	if err != nil {
		return nil, err
	}
	return patterns, err
}

// Parse yaml and return map of urlPattern
func parseJSON(jsn []byte) ([]urlPattern, error) {
	patterns := []urlPattern{}
	err := json.Unmarshal(jsn, &patterns)
	if err != nil {
		return nil, err
	}
	return patterns, err
}

// Take parsed yaml or json, create paths list and return it
func buildMap(patterns []urlPattern) map[string]string {
	paths := make(map[string]string)
	for _, pu := range patterns {
		paths[pu.Path] = pu.URL
	}
	return paths
}
