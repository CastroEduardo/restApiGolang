package sessionuserservice

import (
	"context"
	"fmt"
	"log"
	"os" // get an object type
	"rest-api-golang/src/connectDb"
	"rest-api-golang/src/dbContext/companyservice"
	"rest-api-golang/src/dbContext/privilegeroluserservice"
	"rest-api-golang/src/dbContext/usersservice"
	"rest-api-golang/src/models/authinterfaces"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ClientMongo *mongo.Client
var nameCollection = "sessionsUsers"

//var client *mongo.Client
var collection *mongo.Collection

func settingsCollections() {
	ClientMongo = connectDb.ClientMongo
	if ClientMongo != nil {
		//fmt.Println(os.Getenv("TOKEN_HASH")))
		collection = ClientMongo.Database(os.Getenv("DB_NAME")).Collection(nameCollection)
	}
}

func Add(Model authinterfaces.SessionUser) string {
	settingsCollections()

	LogoutSessionToIdUser(Model.IdUser) //Disable Active Session

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

func LogoutSessionToIdUser(idUser string) bool {
	var sessionUser authinterfaces.SessionUser
	sessionUser = FindToActiveIdUser(idUser)
	if sessionUser.IdCompany != "" {
		if sessionUser.Active {
			sessionUser.Active = false
			sessionUser.DateLogout = time.Now()
			UpdateOne(sessionUser)
		}
	}

	return true
}

func UpdateOne(ModelUpdate authinterfaces.SessionUser) bool {
	settingsCollections()

	//var modelSend authinterfaces.User
	if collection != nil {

		var id = ModelUpdate.ID
		ModelUpdate.ID = ""

		update := bson.M{
			"$set": ModelUpdate,
		}
		// update := bson.M{"$set": bson.M{}}
		docID, _ := primitive.ObjectIDFromHex(id)
		_, err := collection.UpdateOne(context.TODO(), bson.M{"_id": docID}, update)

		if err != nil {
			log.Fatalln("Error on inserting new Hero", err)
			return false
		}
	}

	return false
}

func FindToId(id string) authinterfaces.SessionUser {
	settingsCollections()
	var modelSend authinterfaces.SessionUser
	if collection != nil {
		//transform string _id to Object
		docID, _ := primitive.ObjectIDFromHex(id)
		doc := collection.FindOne(context.TODO(), bson.M{"_id": docID})
		doc.Decode(&modelSend)

	}
	return modelSend
}

func FindToActiveIdUser(idUser string) authinterfaces.SessionUser {
	settingsCollections()
	var modelSend authinterfaces.SessionUser
	if collection != nil {
		//transform string _id to Object
		//docID, _ := primitive.ObjectIDFromHex(idUser)
		doc := collection.FindOne(context.TODO(), bson.M{"iduser": idUser, "active": true})
		doc.Decode(&modelSend)
		// if err != nil {
		// 	log.Fatal(err)
		// 	return modelSend
		// }
	}
	return modelSend
}

func GetList() []authinterfaces.SessionUser {
	settingsCollections()
	var list []authinterfaces.SessionUser
	if collection != nil {
		//transform string _id to Object
		//docID, _ := primitive.ObjectIDFromHex("5e78131bcf026003ec8cb639")
		doc, _ := collection.Find(context.TODO(), bson.M{})
		//doc.Decode(&hero)
		var hero authinterfaces.SessionUser
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

func GetIdSessionToToken(tokenHeader string) string {
	splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
	tokenPart := splitted[1]                    //Grab the token part, what we are truly interested in
	tk := &authinterfaces.Token{}

	//fmt.Println(tokenPart)
	_, errt := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_HASH")), nil
	})
	if errt != nil {
		return ""
	}
	return tk.IdSession

}

func GetClaimForToken(tokenHeader string) authinterfaces.ClaimSession {

	SendModel := authinterfaces.ClaimSession{}

	if tokenHeader == "" {
		return SendModel
	}

	settingsCollections()

	splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
	tokenPart := splitted[1]                    //Grab the token part, what we are truly interested in
	tk := &authinterfaces.Token{}

	//fmt.Println(tokenPart)
	_, errt := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_HASH")), nil
	})

	if errt != nil {
		return SendModel
	}

	if tk.IdSession == "" {
		return SendModel
	}

	var dataSession authinterfaces.SessionUser
	if collection != nil {
		//transform string _id to Object
		docID, _ := primitive.ObjectIDFromHex(tk.IdSession)
		doc := collection.FindOne(context.TODO(), bson.M{"_id": docID})
		doc.Decode(&dataSession)
	}

	var dataCompany authinterfaces.Company
	dataCompany = companyservice.FindToId(dataSession.IdCompany)
	SendModel.Company = dataCompany

	var dataUser authinterfaces.User
	dataUser = usersservice.FindToId(dataSession.IdUser)
	dataUser.Password = ""
	SendModel.User = dataUser

	var dataPrivilegesRol authinterfaces.PrivilegeRolUser
	dataPrivilegesRol = privilegeroluserservice.FindToIdRol(dataUser.IdRol)
	SendModel.PrivilegeRolUser = dataPrivilegesRol

	//fmt.Println(dataPrivilegesRol)

	return SendModel
}
