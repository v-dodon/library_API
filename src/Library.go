package main

import (
	. "net/http"
	. "strconv"
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

var limitConst = "10"
var offsetConst = "0"

func CreateBook(w ResponseWriter, request *Request) {
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

func GetBook(w ResponseWriter, request *Request) {
	if (len(library) == 0) {
		Error(w, "There are no books", StatusNotFound)
		return
	}
	var found bool
	values, err := url.ParseQuery(request.URL.RawQuery)
	var limit int
	var offset int
	var id int
	limitQuery := values.Get("limit")
	offsetQuery := values.Get("offset")
	idQuery := values.Get("id")
	if (len(idQuery) > 0) {
		id, err = Atoi(idQuery)
		if (err != nil) {
			Error(w, err.Error(), StatusBadRequest)
			return
		}
	}

	title := values.Get("title")
	if (len(offsetQuery) == 0) {
		offset, err = Atoi(offsetConst)
	} else {
		offset, _ = Atoi(offsetQuery)
		if (err != nil) {
			Error(w, err.Error(), StatusBadRequest)
			return
		}
	}
	if (len(limitQuery) == 0) {
		limit, _ = Atoi(limitConst)
	} else {
		limit, _ = Atoi(limitQuery)
		if (err != nil) {
			Error(w, err.Error(), StatusBadRequest)
			return
		}
	}
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	if (len(values) == 0) {
		json.NewEncoder(w).Encode(library)
	} else if (len(idQuery) != 0 || len(title) != 0) {
		for _, element := range library {
			if (element.Id == id || element.Title == title) {
				json.NewEncoder(w).Encode(element)
				w.WriteHeader(StatusOK)
				found = true
				return
			}
		}
		if (!found) {
			Error(w, "Not found book with id: " + Itoa(id), StatusNotFound)
			return
		}
	} else {
		if (limit + offset > len(library)) {
			json.NewEncoder(w).Encode(library[offset:len(library)])
		} else {
			json.NewEncoder(w).Encode(library[offset:limit + offset])
		}

	}
}

func UpdateBook(w ResponseWriter, request *Request) {
	urlParams := mux.Vars(request)
	var bookFromRequest Book
	var found bool
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
	id, err := Atoi(urlParams["id"])
	if (err != nil) {
		Error(w, err.Error(), StatusBadRequest)
		return
	}
	for index, _ := range library {
		if library[index].Id == id {
			bookFromRequest.Id = id
			library[index] = bookFromRequest
			found = true
			break
		}
	}
	if (!found) {
		Error(w, "Not found book with id: " + Itoa(id), StatusNotFound)
		return
	}
	w.WriteHeader(StatusOK)
}

func DeleteBook(w ResponseWriter, request *Request) {
	var found bool
	urlParams := mux.Vars(request)
	id, err := Atoi(urlParams["id"])
	if (err != nil) {
		Error(w, err.Error(), StatusBadRequest)
		return
	}
	for index, _ := range library {
		if library[index].Id == id {
			library = append(library[:index], library[index + 1:]...)
			break
			found = true
		}
	}
	if (!found) {
		Error(w, "Not found book with id: " + Itoa(id), StatusNotFound)
		return
	}
	w.WriteHeader(StatusOK)
}

func main() {
	gorillaRoute := mux.NewRouter()
	gorillaRoute.HandleFunc("/book", func(w ResponseWriter, request *Request) {
		switch request.Method {

		case MethodPost:
			CreateBook(w, request)
		case MethodGet:
			GetBook(w, request)

		}
	})
	gorillaRoute.HandleFunc("/book/{id:[0-9]+}", func(w ResponseWriter, request *Request) {
		switch request.Method {

		case MethodPut:
			UpdateBook(w, request)
		case MethodDelete:
			DeleteBook(w, request)

		}
	})
	Handle("/", gorillaRoute)
	ListenAndServe(":8080", nil)
}