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
type libraryHandler struct {
	library map[string]book
	mux     sync.Mutex
}

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

func (l *libraryHandler) getBooks(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	fmt.Println("Sending all books")

	bytes, marshallingError := json.Marshal(l.library)
	if marshallingError != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(marshallingError.Error()))
	}
	fmt.Fprintf(w, string(bytes))
}

func (l *libraryHandler) getBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	fmt.Println("Sending one book")

	wantedBook := r.URL.Query().Get("ISBN")
	fmt.Println(wantedBook)

	value, found := l.library[wantedBook]
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

}

func (l *libraryHandler) addBook(w http.ResponseWriter, r *http.Request) {
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

	l.mux.Lock()
	l.library[newBook.ISBN] = newBook
	l.mux.Unlock()

}

func (l *libraryHandler) updateBook(w http.ResponseWriter, r *http.Request) {
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

	l.mux.Lock()
	_, found := l.library[newBook.ISBN]
	if found == false {
		fmt.Fprintf(w, "Book with ISBN %s doesn't exist, ignoring", newBook.ISBN)
	} else {

		marshallingError := json.Unmarshal(body, &newBook)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(marshallingError.Error()))
			return
		}

		l.library[newBook.ISBN] = newBook

	}
	l.mux.Unlock()
}

func (l *libraryHandler) deleteBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	wantedBook := r.URL.Query().Get("ISBN")
	l.mux.Lock()
	_, found := l.library[wantedBook]
	if found == false {
		fmt.Fprintf(w, "Doesn't exist")
	} else {
		delete(l.library, wantedBook)
		fmt.Fprintf(w, "Successfully deleted book with ISBN : %s from library \n", wantedBook)
	}
	l.mux.Unlock()
}

func createLibrary() *libraryHandler {
	return &libraryHandler{
		library: map[string]book{},
	}
}

func main() {
	// Mock data initialization
	newLibrary := createLibrary()
	testbook1 := book{ISBN: strconv.Itoa(111111), Title: "Cooking 1", Pages: strconv.Itoa(240), Year: strconv.Itoa(2003), Author: &author{Name: "Tarik", Surname: "Redzepagic"}}
	testbook2 := book{ISBN: strconv.Itoa(222222), Title: "Farming 1", Pages: strconv.Itoa(300), Year: strconv.Itoa(2005), Author: &author{Name: "Kirat", Surname: "Pagicredz"}}
	newLibrary.library[testbook1.ISBN] = testbook1
	newLibrary.library[testbook2.ISBN] = testbook2

	http.HandleFunc("/get-books", newLibrary.getBooks)
	http.HandleFunc("/get-book", newLibrary.getBook)
	http.HandleFunc("/add-book", newLibrary.addBook)
	http.HandleFunc("/update-book", newLibrary.updateBook)
	http.HandleFunc("/delete-book", newLibrary.deleteBook)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
