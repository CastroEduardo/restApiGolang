package connectDb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ClientMongo *mongo.Client

func Connect() bool {

	URL_MONGO := "mongodb://admin:password1@127.0.0.1:27017/admin?clickshield?replicaSet=rs0&connect=direct"
	//mongodb://castro:555555@172.16.18:27017
	urlDb := URL_MONGO //os.Getenv("URL_MONGODB")
	fmt.Println(urlDb)

	// Rest of the code will go here
	// Set client options
	clientOptions := options.Client().ApplyURI(urlDb)
	//Context = context.TODO()
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	ClientMongo = client
	if err != nil {

		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
		return false
	}
	//fmt.Println("Connected to MongoDB!")
	return true

}

func Status() *mongo.Client {

	// Check the connection
	err := ClientMongo.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
		result := Connect()
		if result {
			return ClientMongo
		}
		return ClientMongo
	}

	return ClientMongo

}
