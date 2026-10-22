package main

import (
	"animals/internal/handlers"
	"animals/internal/models"
	"animals/internal/services"
	"log"
	"net/http"
)

// TODO: JSON serialization / deserialization

func main() {
	db, err := services.SetupDbConn()
	if err != nil {
		log.Fatalf("Could not create a database: %s", err)
	}

	app := models.CreateAppState("Animals Shop", db)
	handlerState := handlers.HandlerState{
		AppState: app,
	}

	http.HandleFunc("/", handlerState.HandleRoot)
	// TODO: add handler for a FORM, this form should create a new instance of an animal

	server := &http.Server{
		Addr: ":4000",
	}

	log.Print("Server will be started at port :4000")

	server.ListenAndServe()
}
