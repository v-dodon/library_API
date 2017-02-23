package main

import (
	. "net/http"
	"github.com/gorilla/mux"
)

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