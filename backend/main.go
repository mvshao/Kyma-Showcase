package main

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

// initAPIHandler initializes a handler for the API.
func initAPIHandler() (Handler, error) {
	address := os.Getenv("REDIS_URL")
	if address == "" {
		log.Fatal("Failed to read REDIS_URL from .env file")
	}

	database := NewDatabaseConnection(address, os.Getenv("REDIS_PASSWORD"))
	err := database.Connect()
	if err != nil {
		return Handler{}, err
	}
	apiHandler := NewHandler(database, NewIdGenerator())
	return apiHandler, nil
}

// main contains all the function handlers and initializes the database connection.
func main() {
	mux := mux.NewRouter()

	APIhandler, err := initAPIHandler()
	if err != nil {
		log.Fatalf("Error connecting to database: %s", err)
	}

	mux.HandleFunc("/v1/images", APIhandler.DBGetAllHandler).Methods("GET")
	mux.HandleFunc("/v1/images/{id}", APIhandler.DBGetHandler).Methods("GET")
	mux.HandleFunc("/v1/images", APIhandler.DBPostHandler).Methods("POST")

	handler := cors.Default().Handler(mux)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Failed to read PORT from .env file")
	}

	err = http.ListenAndServe(":"+port, handler)
	if err != nil {
		log.Fatalf("Starting server at port %s failed!", port)
	}
}
