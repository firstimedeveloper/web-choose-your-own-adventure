package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

func NewHandler(s StoryArc) http.Handler {
	return handler{s}
}

type handler struct {
	s StoryArc
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	story := h.s

	t := template.Must(template.ParseFiles("index.html"))
	err := t.Execute(w, story)

	if err != nil {
		log.Fatalf("template execution: %s", err)
	}
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

	log.Fatal(http.ListenAndServe(":8080", NewHandler(parsedJSON["intro"])))

	fmt.Println("Welcome to Choose your Own Adventure!\nAn interactive story where you dictate what happens.\nPress enter to continue.")
	var temp string
	fmt.Scanln(&temp)
	storyOver := false
	for !storyOver {
		startArc("intro", parsedJSON)
		storyOver = true
	}
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

func startArc(key string, parsedJSON map[string]StoryArc) {
	//print the title of the StoryArc
	fmt.Println(parsedJSON[key].Title)
	//print story lines secuentially, requiring the user to press enter inbetween each line
	var temp int
	for _, lines := range parsedJSON[key].Story {
		fmt.Println(lines)
		fmt.Scanln(&temp)
	}

	if len(parsedJSON[key].Options) == 0 {
		fmt.Println("The End. What a great Adventure!")
		return
	}

	var count int
	for index, option := range parsedJSON[key].Options {
		fmt.Printf("Press %d: %s\n", index+1, option.Text)
		count++
	}
	fmt.Print("Enter a number: ")
	fmt.Scanln(&temp)
	for temp < 1 || temp > count {
		fmt.Print("Please enter a valid number: ")
		fmt.Scanln(&temp)
	}

	startArc(parsedJSON[key].Options[temp-1].Arc, parsedJSON)

}
