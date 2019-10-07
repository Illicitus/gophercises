package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"text/template"
)

type story map[string]storyBlock

type storyBlock struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}

// Awesome
type handlerOptions func(h *handler)

func withTemplate(t *template.Template) handlerOptions {
	return func(h *handler) {
		h.t = t
	}
}

func (s *story) readJSON(r io.Reader) {
	reader := json.NewDecoder(r)

	if err := reader.Decode(s); err != nil {
		panic(err)
	}
}

func newHandler(s story, opts ...handlerOptions) http.Handler {
	t := template.Must(template.New("").Parse(defaultTmpHandler))
	h := handler{s, t, defaultPath}
	for _, opt := range opts {
		opt(&h)
	}
	return h
}

type handler struct {
	s story
	t *template.Template
	p func(r *http.Request) string
}

func defaultPath(r *http.Request) string {
	path := r.URL.Path
	if path == "" || path == "/" {
		path = "/intro"
	}

	return path[1:]
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.p(r)

	if newStory, ok := h.s[path]; ok {

		// Render the template data and return respons
		err := h.t.Execute(w, newStory)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Who knows...", 400)
		}
		return
	}
	http.Error(w, "Not found", 404)
}

var defaultTmpHandler = func() string {
	file, err := ioutil.ReadFile("templates/story.html")

	if err != nil {
		panic(err)
	}

	return string(file)
}()
