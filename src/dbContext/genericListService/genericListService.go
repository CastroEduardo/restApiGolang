package genericlistservice

import (

	// get an object type

	"context"
	"fmt"
	"log"
	"os"
	"rest-api-golang/src/connectDb"
	"rest-api-golang/src/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ClientMongo *mongo.Client
var nameCollection = "genericList"

//var client *mongo.Client
var collection *mongo.Collection

func settingsCollections() {
	ClientMongo = connectDb.ClientMongo
	if ClientMongo != nil {
		//fmt.Println(os.Getenv("TOKEN_HASH")))
		collection = ClientMongo.Database(os.Getenv("DB_NAME")).Collection(nameCollection)
	}
}

func Add(Model models.GenericList) string {
	settingsCollections()

	if collection != nil {
		insertResult, err := collection.InsertOne(context.TODO(), Model)
		if err != nil {
			log.Fatalln("Error on inserting new Hero", err)
			return ""
		}
		//get id Add
		if oid, ok := insertResult.InsertedID.(primitive.ObjectID); ok {
			return oid.Hex()
		}
	}
	return ""
}

func GetListForIdCompany(idCompany string) []models.GenericList {
	settingsCollections()
	var list []models.GenericList
	if collection != nil {
		//transform string _id to Object
		//docID, _ := primitive.ObjectIDFromHex("5e78131bcf026003ec8cb639")
		doc, _ := collection.Find(context.TODO(), bson.M{"idCompany": idCompany})
		//doc.Decode(&hero)
		var hero models.GenericList
		for doc.Next(context.TODO()) {
			// Declare a result BSON object
			//var result bson.M
			err := doc.Decode(&hero)
			if err != nil {
				fmt.Println(hero)
			}
			list = append(list, hero)
		}
	}

	return list
}

func GetListForIdCompanyAndIdentity(idCompany string, Identity string) []models.GenericList {
	settingsCollections()
	var list []models.GenericList
	if collection != nil {
		//transform string _id to Object
		//docID, _ := primitive.ObjectIDFromHex("5e78131bcf026003ec8cb639")
		doc, _ := collection.Find(context.TODO(), bson.M{"idCompany": idCompany, "identity": Identity})
		//doc.Decode(&hero)
		var hero models.GenericList
		for doc.Next(context.TODO()) {
			// Declare a result BSON object
			//var result bson.M
			err := doc.Decode(&hero)
			if err != nil {
				fmt.Println(hero)
			}
			list = append(list, hero)
		}
	}

	return list
}
