package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Our library for storing books, key being its ISBN (presumed unique), value being the book itself.
var library = map[int]book{}

type author struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

type book struct {
	ISBN   int     `json:"ISBN"`
	Title  string  `json:"randnumber"`
	Pages  int     `json:"pages"`
	Year   int     `json:"year"`
	Author *author `json:"author"`
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	bytes, _ := json.Marshal(library)
	fmt.Fprintf(w, string(bytes))
}

func getSpecBook(w http.ResponseWriter, r *http.Request) {

}

func createBook(w http.ResponseWriter, r *http.Request) {

}

func updateBook(w http.ResponseWriter, r *http.Request) {

}

func deleteBook(w http.ResponseWriter, r *http.Request) {

}

func main() {
	// Mock data initialization
	testbook1 := book{ISBN: 111111, Title: "Cooking 1", Pages: 240, Year: 2003, Author: &author{Name: "Tarik", Surname: "Redzepagic"}}
	testbook2 := book{ISBN: 222222, Title: "Farming 1", Pages: 300, Year: 2005, Author: &author{Name: "Kirat", Surname: "Pagicredz"}}
	library[testbook1.ISBN] = testbook1
	library[testbook2.ISBN] = testbook2

	http.HandleFunc("/", getBooks)
	http.HandleFunc("/getBooks", getBooks)
	http.HandleFunc("/getSpecificBook", getSpecBook)
	http.HandleFunc("/createBook", createBook)
	http.HandleFunc("/updateBook", updateBook)
	http.HandleFunc("/deleteBook", deleteBook)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
