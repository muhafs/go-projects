package models

import (
	"github.com/muhafs/go-bookstore/pkg/config"
	"gorm.io/gorm"
)

var db *gorm.DB

type Book struct {
	gorm.Model
	Title       string `json:"title"`
	Author      string `json:"author"`
	Publication string `json:"publication"`
}

func init() {
	// Connect to database
	config.Conntect()
	db = config.GetDB()

	// Migrate the schema
	db.AutoMigrate(&Book{})
}

func ListBooks() []Book {
	var books []Book
	db.Find(&books)

	return books
}

func GetBook(id int) (*Book, *gorm.DB) {
	var book Book
	d := db.First(&book, id)

	return &book, d
}

func (b *Book) CreateBook() *Book {
	db.Create(&b)

	return b
}

// func (b *Book) UpdateBook() {}

func DeleteBook(id int) Book {
	var book Book
	db.Where("ID = ?", id).Unscoped().Delete(&book)

	return book
}
