package main

import (
	"../../adventure"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := flag.Int("port", 3000, "The port on which to start the web application")
	file := flag.String("file", "gopher.json", "The JSON file with the CYOA story")
	flag.Parse()

	fmt.Printf("Using the story in %s.\n", *file)

	f, err := os.Open(*file)

	if err != nil {
		panic(err)
	}

	story, err := adventure.JsonStory(f)

	if err != nil {
		panic(err)
	}

	h := adventure.NewHandler(story)

	fmt.Printf("Starting the server on port %d...\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
}
