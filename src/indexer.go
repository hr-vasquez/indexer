package main

import (
	"log"
	"os"
)

var logger = log.New(os.Stdout, "", log.LstdFlags)

func main() {
	// Gets the first argument, it should be the path to the database files
	sourceFolder := os.Args[1]

	// Gets the file info, if the argument is a valid path
	_, err := os.Stat(sourceFolder)
	if err != nil {
		logger.Println("The argument is not a valid folder or file.")
		os.Exit(1)
	}

	// Index database source into a full text search engine
	response := IndexSource(sourceFolder)

	logger.Println(response)
	logger.Println("Success in indexing the data of [" + sourceFolder + "]")
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
