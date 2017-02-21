package main

import (
	"net/http"
	//"io/ioutil"
	"encoding/json"
	"strings"
	"net/url"
)

type Book struct {
	Id    int
	Title string
}

var library []Book

func createBook(w http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var book Book
	err := request.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	decodeError := decoder.Decode(&book)
	if decodeError != nil {
		http.Error(w, decodeError.Error(), http.StatusInternalServerError)
		return
	}
	if (strings.TrimSpace(book.Title) == "") {
		http.Error(w, "The title param is missing from the body", http.StatusBadRequest)
		return
	}
	book.Id = len(library)
	library = append(library, book)
	w.WriteHeader(http.StatusCreated)
}

func getBook(w http.ResponseWriter, request *http.Request) {
	_, err := url.ParseQuery(request.URL.RawQuery)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(library)
}

func main() {
	http.HandleFunc("/book", func(w http.ResponseWriter, request *http.Request) {
		switch request.Method {

		case "POST":
			createBook(w, request)
		case "GET":
			getBook(w, request)

		}
	})
	http.ListenAndServe(":8080", nil)
}
