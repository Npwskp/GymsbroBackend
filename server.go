package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	dbmongo "github.com/Npwskp/GymsbroBackend/api/v1/db"
	"github.com/Npwskp/GymsbroBackend/api/v1/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"

	_ "github.com/Npwskp/GymsbroBackend/docs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Npwskp/GymsbroBackend/api/v1/auth"
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
// @BasePath	/api/v1
// @schemes		http https
// @host		35.240.232.32:8080
// @version					1.0

// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description Enter your bearer token in the format: Bearer {token}

// @SecurityDefinition.apiKey cookieAuth
// @in cookie
// @name jwt

// @Security Bearer
// @Security cookieAuth
func main() {
	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000, https://gyms-bro-fe.vercel.app/, https://localhost:8080, http://35.240.232.32:8080",
		AllowCredentials: true,
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin,Authorization",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))
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

	if err := dbmongo.CreateIndexes(mg.Db); err != nil {
		log.Fatalf("Error creating indexes: %v", err)
	}

	utils.InjectApp(app, mg.Db)

	app.Get("/swagger/*", swagger.New(swagger.Config{
		Title: "GymsBro API",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		fmt.Println(time.Now().Format("2006-01-02"))
		return c.SendString("Hello, World 👋!")
	})

	auth.InitGoogleOAuth()

	app.Listen(":8080")
}
