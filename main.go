package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

// Our library for storing books, key being its ISBN (presumed unique), value being the book itself.
var library = map[string]book{}
var mutex sync.Mutex

type author struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

type book struct {
	ISBN   string  `json:"isbn"`
	Title  string  `json:"title"`
	Pages  string  `json:"pages"`
	Year   string  `json:"year"`
	Author *author `json:"author"`
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	fmt.Println("Sending all books")

	mutex.Lock()
	bytes, marshallingError := json.Marshal(library)
	mutex.Unlock()
	if marshallingError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(marshallingError.Error()))
	}
	fmt.Fprintf(w, string(bytes))
}

func getBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	fmt.Println("Sending one book")

	wantedBook := r.URL.Query().Get("ISBN")
	fmt.Println(wantedBook)

	mutex.Lock()
	value, found := library[wantedBook]
	if found == false {
		fmt.Fprintf(w, "Book with ISBN %s doesn't exist", wantedBook)
	} else {
		bytes, marshallingError := json.Marshal(value)
		if marshallingError != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(marshallingError.Error()))
		}
		fmt.Fprintf(w, string(bytes))

	}
	mutex.Unlock()

}

func addBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	contentChecker := r.Header.Get("content-type")
	if contentChecker != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need application/json', but got '%s'", contentChecker)))
		return
	}

	var newBook book
	err = json.Unmarshal(body, &newBook)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	mutex.Lock()
	library[newBook.ISBN] = newBook
	mutex.Unlock()

}

func updateBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	contentChecker := r.Header.Get("content-type")
	if contentChecker != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need application/json', but got '%s'", contentChecker)))
		return
	}

	var newBook book
	err = json.Unmarshal(body, &newBook)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	mutex.Lock()

	_, found := library[newBook.ISBN]
	if found == false {
		fmt.Fprintf(w, "Book with ISBN %s doesn't exist, ignoring", newBook.ISBN)
	} else {

		marshallingError := json.Unmarshal(body, &newBook)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(marshallingError.Error()))
			return
		}

		library[newBook.ISBN] = newBook

	}
	mutex.Unlock()

}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	wantedBook := r.URL.Query().Get("ISBN")
	_, found := library[wantedBook]
	if found == false {
		fmt.Fprintf(w, "Doesn't exist")
	} else {
		delete(library, wantedBook)
		fmt.Fprintf(w, "Successfully deleted book with ISBN : %s from library \n", wantedBook)
	}
}

func main() {
	// Mock data initialization
	testbook1 := book{ISBN: strconv.Itoa(111111), Title: "Cooking 1", Pages: strconv.Itoa(240), Year: strconv.Itoa(2003), Author: &author{Name: "Tarik", Surname: "Redzepagic"}}
	testbook2 := book{ISBN: strconv.Itoa(222222), Title: "Farming 1", Pages: strconv.Itoa(300), Year: strconv.Itoa(2005), Author: &author{Name: "Kirat", Surname: "Pagicredz"}}
	library[testbook1.ISBN] = testbook1
	library[testbook2.ISBN] = testbook2

	http.HandleFunc("/getBooks", getBooks)
	http.HandleFunc("/getBook", getBook)
	http.HandleFunc("/addBook", addBook)
	http.HandleFunc("/updateBook", updateBook)
	http.HandleFunc("/deleteBook", deleteBook)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
