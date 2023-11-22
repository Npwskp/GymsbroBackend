package main

import (
	"context"
	"time"

	"github.com/Npwskp/GymsbroBackend/src/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

var mg *MongoInstance

const dbname = "GymsBro"
const mongoURI = "mongodb+srv://npwskp:YV57BjDS6DwFzmxT@npwskp.l9cg7pi.mongodb.net/"

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

func main() {
	app := fiber.New()

	app.Use(logger.New())

	connectDB()
	defer disconnectDB()

	utils.InjectApp(app, mg.Db)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Listen(":8080")
}
