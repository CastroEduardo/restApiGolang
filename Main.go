package main

import (
	"fmt"
	"os"
	"rest-api-golang/src/connectDb"
	"rest-api-golang/src/loggerservice"

	"rest-api-golang/src/routes"
)

func main() {

	var modeApp = os.Getenv("MODE_APP")
	if modeApp == "DEV" {

		fmt.Println("dev")
		//send enviroment
		os.Setenv("APP_PORT", "1111")
		os.Setenv("TOKEN_HASH", "PassworsStrongTo254444Token")
		os.Setenv("DB_NAME", "dbdemo")
		fmt.Println("APP_MODE: " + modeApp)
	}

	loggerservice.Add("  ")
	loggerservice.Add("#####*** SETTINGS APP ***#####")
	loggerservice.Add("APP_PORT:" + os.Getenv("APP_PORT"))
	loggerservice.Add("TOKEN_HASH:" + os.Getenv("TOKEN_HASH"))
	loggerservice.Add("DB_NAME:" + os.Getenv("DB_NAME"))
	loggerservice.Add("URL_MONGODB:" + os.Getenv("URL_MONGODB"))

	//os.Setenv("URL_MONGODB", "PATH DB")
	// fmt.Println(os.Getenv("DB_NAME"))
	// var pass string
	// pass = "$2a$04$C6RHde4Vmnj6n/49J2Nfd.eEcsqKj0yxk78/xTt.f2A.8Qbv8a47y"
	// rsu := utils.Descrypt(pass, []byte("231154"))
	// fmt.Println(rsu)

	connect := connectDb.Connect() //try connect mongoDB
	if connect {
		fmt.Println("conectado Db ")
		loggerservice.Add("dbConnect:" + "True")
	}
	loggerservice.Add("#####*** END SETTINGS APP ***###")
	loggerservice.Add("  ")

	var APP_PORT = os.Getenv("APP_PORT")

	if APP_PORT == "" {
		APP_PORT = "1111"
	}

	routes.SetupServer(APP_PORT)
}
