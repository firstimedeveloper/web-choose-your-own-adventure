package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func NewHandler(s map[string]StoryArc) http.Handler {
	return handler{s}
}

type handler struct {
	s map[string]StoryArc
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	// "/intro" => "intro"
	path = path[1:]
	//story := h.s[path]

	t := template.Must(template.ParseFiles("index.html"))

	// ["intro"]
	if chapter, ok := h.s[path]; ok {
		err := t.Execute(w, chapter)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "Chapter not found", http.StatusNotFound)

}

func main() {
	//Read the json file and assign it to content
	content, err := ioutil.ReadFile("gopher.json")
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("File contents: %s", content)

	parsedJSON := make(map[string]StoryArc)
	err = json.Unmarshal(content, &parsedJSON)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Print("Parsed Json file: ", parsedJSON)

	port := os.Getenv("PORT")

	if port == ":" {
		log.Fatal("$PORT must be set")
	}

	log.Fatal(http.ListenAndServe(":"+port, NewHandler(parsedJSON)))

}

//StoryArc is a struct that contains the title, storylines, and options
type StoryArc struct {
	Title   string
	Story   []string
	Options []Option
}

//Option is a struct that containst 2 strings, Text and Arc.
type Option struct {
	Text string
	Arc  string
}
