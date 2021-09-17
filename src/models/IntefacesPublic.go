package models

import (
	"time"
)

// Profile - is the memory representation of one user profile
type Profile struct {
	ID          string    `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string    `json:"username"`
	Password    string    `json:"password"`
	Age         int       `json:"age"`
	LastUpdated time.Time `json:"lastTime"`
}

type LogUser struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	IdUser    string    `json:"idUser,omitempty" bson:"iduser,omitempty"`
	IdCompany string    `json:"idCompany" bson:"idcompany"`
	Log       string    `json:"log" bson:"log"`
	Status    int       `json:"status" bson:"status"`
	Level     int       `json:"level" bson:"level"`
	Date      time.Time `json:"date" bson:"date"`
}

type LogSystem struct {
	ID     string    `json:"id,omitempty" bson:"_id,omitempty"`
	Log    string    `json:"log" bson:"log"`
	Status int       `json:"status" bson:"status"`
	Level  int       `json:"level" bson:"level"`
	Date   time.Time `json:"date" bson:"date"`
}

type GenericList struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	IdKey     string    `json:"idKey" bson:"idKey"`
	Identity  string    `json:"identity" bson:"identity"`
	Name      string    `json:"name" bson:"name"`
	Status    int       `json:"status" bson:"status"`
	Value1    string    `json:"value1" bson:"value1"`
	Value2    string    `json:"value2" bson:"value2"`
	Value3    string    `json:"value3" bson:"value3"`
	Value4    string    `json:"value4" bson:"value4"`
	Value5    string    `json:"value5" bson:"value5"`
	Value6    string    `json:"value6" bson:"value6"`
	IdCompany string    `json:"idCompany" bson:"idCompany"`
	Note      string    `json:"note" bson:"note"`
	Date      time.Time `json:"date" bson:"date"`
}

type RequestListGeneral struct {
	IdCompany string `json:"idCompany" bson:"idCompany"`
	Status    int    `json:"status" bson:"status"`
	Identity  string `json:"identity" bson:"identity"`
}

type RequestGeneric struct {
	Params1 string `json:"params1" bson:"idCompany"`
	Params2 string `json:"params2" bson:"params2"`
	Params3 string `json:"params3" bson:"params3"`
	Params4 string `json:"params4" bson:"params4"`
}
