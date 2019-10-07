package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// Parse json
	jsn := flag.String("json", "gopher.json", "")
	flag.Parse()

	f, err := os.Open(*jsn)
	if err != nil {
		panic(err)
	}

	// Read file and generate Story struct
	newStory := story{}
	newStory.readJSON(f)

	// Run server
	nh := newHandler(newStory)
	fmt.Println("Start server")
	log.Fatal(http.ListenAndServe("localhost:8080", nh))

}
