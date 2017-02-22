package main

import (
	. "net/http"
	"strconv"
	"encoding/json"
	"strings"
	"net/url"
	"github.com/gorilla/mux"
)

type Book struct {
	Id    int
	Title string
}

var library []Book

func createBook(w ResponseWriter, request *Request) {
	decoder := json.NewDecoder(request.Body)
	var book Book
	err := request.ParseForm()
	if err != nil {
		Error(w, err.Error(), StatusInternalServerError)
		return
	}
	decodeError := decoder.Decode(&book)
	if decodeError != nil {
		Error(w, decodeError.Error(), StatusInternalServerError)
		return
	}
	if (strings.TrimSpace(book.Title) == "") {
		Error(w, "The title param is missing from the body", StatusBadRequest)
		return
	}
	book.Id = len(library)
	library = append(library, book)
	w.WriteHeader(StatusCreated)
}

func getBook(w ResponseWriter, request *Request) {
	values, err := url.ParseQuery(request.URL.RawQuery)
	limit := values.Get("limit")
	offset := values.Get("offset")
	id := values.Get("id")
	title := values.Get("title")
	if(len(offset) == 0){
		offset = 0
	}
	if (len(limit) == 0){
		limit = 10
	}
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	if (len(values) == 0) {
		json.NewEncoder(w).Encode(library)
	} else {
		json.NewEncoder(w).Encode(library[offset:limit])
	}
}

func updateBook(w ResponseWriter, request *Request) {
	urlParams := mux.Vars(request)
	var bookFromRequest Book
	decoder := json.NewDecoder(request.Body)
	err := request.ParseForm()
	if err != nil {
		Error(w, err.Error(), StatusInternalServerError)
		return
	}
	decodeError := decoder.Decode(&bookFromRequest)

	if decodeError != nil {
		Error(w, decodeError.Error(), StatusInternalServerError)
		return
	}
	id, err := strconv.Atoi(urlParams["id"])
	if (err != nil) {
		Error(w, err.Error(), StatusBadRequest)
		return
	}
	for index, _ := range library {
		if library[index].Id == id {
			bookFromRequest.Id = id
			library[index] = bookFromRequest
			break
		}
	}
	w.WriteHeader(StatusOK)
}

func deleteBook(w ResponseWriter, request *Request) {
	urlParams := mux.Vars(request)
	id, err := strconv.Atoi(urlParams["id"])
	if (err != nil) {
		Error(w, err.Error(), StatusBadRequest)
		return
	}
	for index, _ := range library {
		if library[index].Id == id {
			library = append(library[:index], library[index + 1:]...)
			break
		}
	}
	w.WriteHeader(StatusOK)
}

func main() {
	gorillaRoute := mux.NewRouter()
	gorillaRoute.HandleFunc("/book", func(w ResponseWriter, request *Request) {
		switch request.Method {

		case MethodPost:
			createBook(w, request)
		case MethodGet:
			getBook(w, request)

		}
	})
	gorillaRoute.HandleFunc("/book/{id:[0-9]+}", func(w ResponseWriter, request *Request) {
		switch request.Method {

		case MethodPut:
			updateBook(w, request)
		case MethodDelete:
			deleteBook(w, request)

		}
	})
	Handle("/", gorillaRoute)
	ListenAndServe(":8080", nil)
}