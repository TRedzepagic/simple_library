package main

import (
	"net/http"
)

type author struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

type book struct {
	Title  string `json:"randnumber"`
	Pages  int    `json:"pages"`
	Year   int    `json:"year"`
	Author author `json:"author"`
}

func getBooks(w http.ResponseWriter, r *http.Request) {

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

	http.HandleFunc("/getBooks", getBooks)
	http.HandleFunc("getSpecificBook", getSpecBook)
	http.HandleFunc("/createBook", createBook)
	http.HandleFunc("/updateBook", updateBook)
	http.HandleFunc("/deleteBook", deleteBook)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
