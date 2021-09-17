package logsuserservice

import (
	"context"
	"fmt"
	"log"
	"os" // get an object type
	"rest-api-golang/src/connectDb"
	"rest-api-golang/src/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ClientMongo *mongo.Client
var nameCollection = "logsUsers"

//var client *mongo.Client
var collection *mongo.Collection

func settingsCollections() {
	ClientMongo = connectDb.ClientMongo
	if ClientMongo != nil {
		//fmt.Println(os.Getenv("TOKEN_HASH")))
		collection = ClientMongo.Database(os.Getenv("DB_NAME")).Collection(nameCollection)
	}

}

func Add(level int, idUser string, idCompany string, logUser string) string {
	settingsCollections()

	newLog := models.LogUser{
		Log:       logUser,
		Level:     level,
		Status:    1,
		Date:      time.Now(),
		IdUser:    idUser,
		IdCompany: idCompany,
	}

	if collection != nil {
		insertResult, err := collection.InsertOne(context.TODO(), newLog)
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

func GetList() []models.LogUser {
	settingsCollections()
	var list []models.LogUser
	if collection != nil {
		//transform string _id to Object
		//docID, _ := primitive.ObjectIDFromHex("5e78131bcf026003ec8cb639")
		findOptions := options.Find()
		// Sort by `price` field descending
		findOptions.SetSort(bson.D{{"date", -1}})

		doc, _ := collection.Find(context.TODO(), bson.M{})

		//doc.Decode(&hero)
		var hero models.LogUser
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

func init() {
	//fmt.Println("init Service1")
}

func Demo() string {

	return "Envio desde Servicio"
}
