package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Client *mongo.Client
	DB     *mongo.Database
}

var mg MongoInstance

type Employee struct {
	ID     string `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string `json:"name"`
	Age    uint8  `json:"age"`
	Salary uint64 `json:"salary"`
}

func Connect() error {
	// Load .env File
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Extract Mongo DB URI from .env file
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	// Create connection
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	mg = MongoInstance{
		Client: client,
		DB:     client.Database("fire-hrms"),
	}

	return nil
}

func main() {
	// Connect to Monggo DB
	err := Connect()
	if err != nil {
		log.Fatal("Couldn't  connect to database:", err)
		return
	}

	app := fiber.New()

	e := app.Group("/employee")
	e.Get("/", ListEmployees)
	e.Get("/:id", GetEmployee)
	e.Post("/", CreateEmployee)
	e.Put("/:id", UpdateEmployee)
	e.Delete("/:id", DeleteEmployee)

	app.Listen(":3000")
}

func ListEmployees(c *fiber.Ctx) error {
	collection := mg.DB.Collection("employees")

	query := bson.D{{}}
	cur, err := collection.Find(c.Context(), query)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	employess := []Employee{}
	if err := cur.All(c.Context(), &employess); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(employess)
}

func GetEmployee(c *fiber.Ctx) error {
	id := c.Params("id")
	collection := mg.DB.Collection("employees")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}

	employee := new(Employee)

	filter := bson.D{{Key: "_id", Value: objectId}}
	err = collection.FindOne(c.Context(), filter).Decode(employee)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "employee not found",
		})
	}

	return c.Status(201).JSON(employee)
}

func CreateEmployee(c *fiber.Ctx) error {
	collection := mg.DB.Collection("employees")

	employee := new(Employee)
	if err := c.BodyParser(employee); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	rst, err := collection.InsertOne(c.Context(), employee)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	filter := bson.D{{Key: "_id", Value: rst.InsertedID}}
	createdRecord := collection.FindOne(c.Context(), filter)

	createdEmployee := new(Employee)
	createdRecord.Decode(createdEmployee)

	return c.Status(201).JSON(createdEmployee)
}

func UpdateEmployee(c *fiber.Ctx) error {
	id := c.Params("id")
	collection := mg.DB.Collection("employees")

	// create ObjectID from ID string
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}

	// extract employee from request
	employee := new(Employee)
	if err := c.BodyParser(employee); err != nil {
		c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// create filter query and update
	filter := bson.D{{Key: "_id", Value: objectId}}
	update := bson.D{{
		Key: "$set",
		Value: bson.D{
			{Key: "name", Value: employee.Name},
			{Key: "age", Value: employee.Age},
			{Key: "salary", Value: employee.Salary},
		},
	}}

	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Couldn't update the record"})
	}

	if result.MatchedCount < 1 {
		return c.Status(404).JSON(fiber.Map{"message": "record not found"})
	}

	return c.JSON(fiber.Map{"message": "record updated"})
}

func DeleteEmployee(c *fiber.Ctx) error {
	id := c.Params("id")
	collection := mg.DB.Collection("employees")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid ID"})
	}

	filter := bson.D{{Key: "_id", Value: objectId}}
	result, err := collection.DeleteOne(c.Context(), filter)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Couldn't Delete the record"})
	}

	if result.DeletedCount < 1 {
		return c.Status(404).JSON(fiber.Map{"message": "record not found"})
	}

	return c.JSON(fiber.Map{"message": "record deleted"})
}
