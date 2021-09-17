package dbcontext

import (
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
var nameCollection = "heroes"

//var client *mongo.Client
var collection *mongo.Collection

func init() {
	//user.Area()
	//fmt.Println("Iniciando Context")
}

func settingsCollections() {
	ClientMongo = connectDb.ClientMongo
	if ClientMongo != nil {
		//fmt.Println(os.Getenv("TOKEN_HASH")))
		collection = ClientMongo.Database(os.Getenv("DB_NAME")).Collection(nameCollection)
	}

}

func GetOne(name string) models.Profile {
	settingsCollections()

	var hero models.Profile
	if collection != nil {
		//transform string _id to Object
		docID, _ := primitive.ObjectIDFromHex("5e78131bcf026003ec8cb639")
		doc := collection.FindOne(context.TODO(), bson.M{"_id": docID})
		doc.Decode(&hero)

		//fmt.Println(hero.ID)
	}

	return hero
}

func GetList() []models.Profile {
	settingsCollections()
	var heros []models.Profile
	if collection != nil {

		//transform string _id to Object
		//docID, _ := primitive.ObjectIDFromHex("5e78131bcf026003ec8cb639")
		doc, _ := collection.Find(context.TODO(), bson.M{})
		//doc.Decode(&hero)
		var hero models.Profile
		for doc.Next(context.TODO()) {
			// Declare a result BSON object
			//var result bson.M
			err := doc.Decode(&hero)
			if err != nil {
				fmt.Println(hero)
			}
			heros = append(heros, hero)
		}
	}

	return heros
}

func InsertNew(profile models.Profile) bool {
	settingsCollections()

	if collection != nil {
		//transform string _id to Object
		//docID, _ := primitive.ObjectIDFromHex("5e78131bcf026003ec8cb639")
		_, err := collection.InsertOne(context.TODO(), profile)

		if err != nil {
			log.Fatalln("Error on inserting new Hero", err)
			return false
		}

		//fmt.Println(doc)
	}

	return true
}

func Demo() string {

	// user := usersservice.Demo()
	// fmt.Println(user)

	return "admin"
}
