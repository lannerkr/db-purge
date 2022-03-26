package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// const (
// 	uri = "mongodb://mongoadmin:secret@192.168.0.157:27017/?retryWrites=true&w=majority"
// )

func connectdb() (client *mongo.Client) {

	uri := "mongodb://" + configuration.MongoDB + "/?retryWrites=true&w=majority"

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	// defer func() {
	// 	if err = client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()

	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	log.Println("Successfully connected and pinged.")

	// insert or update user history
	//dbConn := client.Database("ldapDB").Collection("user_history")

	return client

}

func getColl(client *mongo.Client) *mongo.Collection {
	return client.Database("ldapDB").Collection("user_history")
}
