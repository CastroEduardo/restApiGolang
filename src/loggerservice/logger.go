package loggerservice

import (
	"log"
	"os"
)

var pathFile *os.File

func Init() {

}

func Add(logger string) {
	file, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	pathFile = file
	if err != nil {
		log.Fatal(pathFile)
	}
	log.SetOutput(file)
	log.Print(logger)

	defer pathFile.Close()
}
