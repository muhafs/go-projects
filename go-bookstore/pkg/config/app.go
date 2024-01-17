package config

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var database *gorm.DB

func Conntect() {
	dsn := "root:root@tcp(127.0.0.1:3306)/bookstore?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("failed to connect database")
	}

	database = db
}

func GetDB() *gorm.DB {
	return database
}
