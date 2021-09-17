package logssystemservice

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
)

var ClientMongo *mongo.Client
var nameCollection = "logsSystem"

//var client *mongo.Client
var collection *mongo.Collection

func settingsCollections() {
	ClientMongo = connectDb.ClientMongo
	if ClientMongo != nil {
		//fmt.Println(os.Getenv("TOKEN_HASH")))
		collection = ClientMongo.Database(os.Getenv("DB_NAME")).Collection(nameCollection)
	}

}

func Add(level int, logSytem string) string {
	settingsCollections()

	newLog := models.LogSystem{
		Log:    logSytem,
		Level:  level,
		Status: 1,
		Date:   time.Now(),
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

func GetList() []models.LogSystem {
	settingsCollections()
	var list []models.LogSystem
	if collection != nil {
		//transform string _id to Object
		//docID, _ := primitive.ObjectIDFromHex("5e78131bcf026003ec8cb639")
		doc, _ := collection.Find(context.TODO(), bson.M{})
		//doc.Decode(&hero)
		var hero models.LogSystem
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
	//fmt.Println("init.....")
}

func Demo() string {

	return "Envio desde Servicio"
}
