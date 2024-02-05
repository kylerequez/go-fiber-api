package common

import (
	"context"
	"errors"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDB() error {
	hostname, port, uriString, database := os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_URI"), os.Getenv("DB_NAME")
	if hostname == "" || port == "" || uriString == "" || database == "" {
		return errors.New("you do not have the necessary information to connect to the database. Please check your .env file.")
	}

	uri := fmt.Sprintf(uriString, hostname, port)

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	fmt.Println(":::-:::\tSuccessfully Connected\t:::-:::")
	DB = client.Database(database)
	return nil
}

func CloseDB() error {
	return DB.Client().Disconnect(context.Background())
}

func GetCollection(name string) (coll *mongo.Collection) {
	return DB.Collection(name)
}
