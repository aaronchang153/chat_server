package main

import (
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	err := RunServer(logger)
	if err != nil {
		log.Fatal(err)
	}
}
