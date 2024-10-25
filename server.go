package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Npwskp/GymsbroBackend/api/v1/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"

	_ "github.com/Npwskp/GymsbroBackend/docs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

var mg *MongoInstance

var mongoURI string
var dbname string

func connectDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return err
	}

	mg = &MongoInstance{
		Client: client,
		Db:     client.Database(dbname),
	}

	return nil
}

func disconnectDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := mg.Client.Disconnect(ctx)
	if err != nil {
		return err
	}

	return nil
}

// @title		GymsBro API
// @description	This is a sample server for GymsBro API.
// @BasePath	/
// @schemes		http https
// @host		localhost:8080
// @version		1.0
func main() {
	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New())
	app.Static("/swagger", "./docs/swagger.json")

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Get environment variables
	dbname = os.Getenv("DB_NAME")
	mongoURI = os.Getenv("MONGO_URI")

	connectDB()
	defer disconnectDB()

	utils.InjectApp(app, mg.Db)

	app.Get("/swagger/*", swagger.HandlerDefault) // default

	app.Get("/", func(c *fiber.Ctx) error {
		fmt.Println(time.Now().Format("2006-01-02"))
		return c.SendString("Hello, World 👋!")
	})

	app.Listen(":8080")
}
