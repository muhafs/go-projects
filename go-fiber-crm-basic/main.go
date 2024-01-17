package main

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2"
	"github.com/muhafs/go-fiber-crm-basic/database"
	"github.com/muhafs/go-fiber-crm-basic/lead"
	"gorm.io/gorm"
)

func setupRoutes(l fiber.Router) {
	l.Get("/", lead.GetLeads)
	l.Get("/:id", lead.GetLead)
	l.Post("/", lead.NewLead)
	l.Delete("/:id", lead.DeleteLead)
}

func initDatabase() {
	var err error
	database.ConnectDB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("Connection opened to database")

	database.ConnectDB.AutoMigrate(&lead.Lead{})
	fmt.Println("Database Migrated")
}

func main() {
	app := fiber.New()
	initDatabase()

	l := app.Group("/api/v1/lead")
	setupRoutes(l)

	app.Listen(":3000")
}
