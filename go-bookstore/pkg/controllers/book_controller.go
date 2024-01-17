package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/muhafs/go-bookstore/pkg/models"
	"github.com/muhafs/go-bookstore/pkg/utils"
)

func GetBooks(w http.ResponseWriter, r *http.Request) {
	book := models.ListBooks()

	utils.RespondWithJSON(w, http.StatusOK, book)
}

func GetBook(w http.ResponseWriter, r *http.Request) {
	bookID := mux.Vars(r)["bookID"]

	ID, err := strconv.Atoi(bookID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Couldn't parse ID: %v", err))
		return
	}

	book, _ := models.GetBook(ID)
	if book.Title == "" {
		utils.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Book not found"))
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, book)
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	reqBook := &models.Book{}
	utils.ParseBody(r, reqBook)

	book := reqBook.CreateBook()

	utils.RespondWithJSON(w, http.StatusCreated, book)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	reqBook := &models.Book{}
	utils.ParseBody(r, reqBook)

	strID := mux.Vars(r)["bookID"]
	ID, err := strconv.Atoi(strID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Couldn't parse ID: %v", err))
	}

	book, db := models.GetBook(ID)
	if reqBook.Title != "" {
		book.Title = reqBook.Title
	}
	if reqBook.Author != "" {
		book.Author = reqBook.Author
	}
	if reqBook.Publication != "" {
		book.Publication = reqBook.Publication
	}

	db.Save(&book)

	utils.RespondWithJSON(w, http.StatusOK, book)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	bookID := mux.Vars(r)["bookID"]

	ID, err := strconv.Atoi(bookID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Couldn't parse ID: %v", err))
	}

	book := models.DeleteBook(ID)

	utils.RespondWithJSON(w, http.StatusOK, book)
}
