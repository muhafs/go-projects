package routes

import (
	"github.com/gorilla/mux"
	"github.com/muhafs/go-bookstore/pkg/controllers"
)

var RegisterBookRoutes = func(router *mux.Router) {
	router.HandleFunc("/books", controllers.GetBooks).Methods("GET")
	router.HandleFunc("/books", controllers.CreateBook).Methods("POST")
	router.HandleFunc("/books/{bookID}", controllers.GetBook).Methods("GET")
	router.HandleFunc("/books/{bookID}", controllers.UpdateBook).Methods("PUT")
	router.HandleFunc("/books/{bookID}", controllers.DeleteBook).Methods("DELETE")
}
