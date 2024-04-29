package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"QA-System/config/config"
)

var MDB *mongo.Collection

func MongodbInit() {
	// Get MongoDB connection information from the configuration file
	user := config.Config.GetString("mongodb.user")
	pass := config.Config.GetString("mongodb.pass")
	host := config.Config.GetString("mongodb.host")
	name := config.Config.GetString("mongodb.db")
	collection := config.Config.GetString("mongodb.collection")

	// Build the MongoDB connection string
	dsn := fmt.Sprintf("mongodb://%v:%v@%v/%v", user, pass, host, name)

	// Create MongoDB client options
	clientOptions := options.Client().ApplyURI(dsn)

	// Create connection context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create MongoDB client
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to create MongoDB client:", err)
	}

	if err :=client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	// Set the MongoDB database
	MDB = client.Database(name).Collection(collection)

	// Print a log message to indicate successful connection to MongoDB
	log.Println("Connected to MongoDB")
}
