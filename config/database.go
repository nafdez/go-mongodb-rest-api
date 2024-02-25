package config

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// ConnectToDB connects to the database using the provided uri and database name in the .env file and returns the database
func ConnectToDB() *mongo.Database {
	ctx := context.TODO()
	uri := os.Getenv("MONGO_URI_DEV")
	dbName := os.Getenv("MONGO_DATABASE")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Error connecting to database. Error: ", err)
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal("Error connecting to database. Error: ", err)
	}
	return client.Database(dbName)
}
